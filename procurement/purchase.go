package procurement

import (
	"context"
	"errors"
	"regexp"
	"strings"
	"time"

	"encore.app/product"
	"encore.dev/types/uuid"
)

// Purchase model
type Purchase struct {
	ID             uuid.UUID `json:"id"`
	PurchaseNumber string    `json:"purchase_number"`
	SupplierID     uuid.UUID `json:"supplier_id"`
	PurchaseDate   time.Time `json:"purchase_date"`
	TotalAmount    float64   `json:"total_amount"`
	Status         string    `json:"status"`
	Notes          string    `json:"notes"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	CreatedBy      uuid.UUID `json:"created_by"`
}

// PurchaseItem model
type PurchaseItem struct {
	ID         uuid.UUID `json:"id"`
	PurchaseID uuid.UUID `json:"purchase_id"`
	ProductID  uuid.UUID `json:"product_id"`
	Quantity   int       `json:"quantity"`
	TotalPrice float64   `json:"total_price"`
	CreatedAt  time.Time `json:"created_at"`
}

// CreatePurchase creates a new purchase order with items
//
//encore:api public method=POST path=/api/purchases
func CreatePurchase(ctx context.Context, req *CreatePurchaseRequest) (Response, error) {
	// Validate request
	if err := req.Validate(); err != nil {
		return Response{Message: "Validation failed"}, err
	}

	// Check if supplier exists
	var supplierExists bool
	err := db.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM suppliers WHERE id = $1)", req.SupplierID).Scan(&supplierExists)
	if err != nil {
		return Response{Message: "Failed to check supplier"}, err
	}
	if !supplierExists {
		return Response{Message: "Supplier not found"}, errors.New("supplier not found")
	}

	// Check if purchase number already exists
	var purchaseNumberExists bool
	err = db.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM purchases WHERE purchase_number = $1)", req.InvoiceNumber).Scan(&purchaseNumberExists)
	if err != nil {
		return Response{Message: "Failed to check purchase number"}, err
	}
	if purchaseNumberExists {
		return Response{Message: "Purchase number already exists"}, errors.New("purchase number already exists")
	}

	// Fetch product prices from product service and validate products exist
	productPrices := make(map[uuid.UUID]float64)
	var totalAmount float64
	for _, item := range req.Items {
		// Get product from product service to get base_price
		productData, err := product.GetProduct(ctx, item.ProductID)
		if err != nil {
			return Response{Message: "Product not found: " + item.ProductID.String()}, errors.New("product not found")
		}

		// Store product price for later use
		productPrices[item.ProductID] = productData.BasePrice

		// Calculate item total using product's base_price
		itemTotal := float64(item.Quantity) * productData.BasePrice
		totalAmount += itemTotal
	}

	// Build notes with expected delivery if provided
	notes := req.Notes
	if !req.ExpectedDelivery.IsZero() {
		if notes != "" {
			notes += "\n"
		}
		notes += "Expected delivery: " + req.ExpectedDelivery.Format("2006-01-02")
	}

	// Create purchase
	var purchaseID uuid.UUID
	err = db.QueryRow(ctx, `
		INSERT INTO purchases (purchase_number, supplier_id, purchase_date, total_amount, status, notes, created_by)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id
	`, req.InvoiceNumber, req.SupplierID, req.OrderDate, totalAmount, "pending", notes, req.CreatedBy).Scan(&purchaseID)
	if err != nil {
		return Response{Message: "Failed to create purchase"}, err
	}

	// Create purchase items using cached product prices
	for _, item := range req.Items {
		// Get cached product price
		unitPrice := productPrices[item.ProductID]

		// Auto-calculate total_price from quantity * product base_price
		totalPrice := float64(item.Quantity) * unitPrice
		_, err = db.Exec(ctx, `
			INSERT INTO purchase_items (purchase_id, product_id, quantity, total_price)
			VALUES ($1, $2, $3, $4)
		`, purchaseID, item.ProductID, item.Quantity, totalPrice)
		if err != nil {
			// If item creation fails, try to clean up the purchase
			db.Exec(ctx, "DELETE FROM purchases WHERE id = $1", purchaseID)
			return Response{Message: "Failed to create purchase item"}, err
		}
	}

	return Response{Message: "Purchase created successfully"}, nil
}

// GetAllPurchases retrieves all purchases with supplier information and item counts
//
//encore:api public method=GET path=/api/purchases
func GetAllPurchases(ctx context.Context) (ListPurchasesResponse, error) {
	rows, err := db.Query(ctx, `
		SELECT 
			p.id,
			p.purchase_number,
			s.name as supplier_name,
			p.purchase_date,
			p.total_amount,
			p.status,
			p.notes,
			COALESCE(COUNT(pi.id), 0) as total_item
		FROM purchases p
		LEFT JOIN suppliers s ON p.supplier_id = s.id
		LEFT JOIN purchase_items pi ON p.id = pi.purchase_id
		GROUP BY p.id, p.purchase_number, s.name, p.purchase_date, p.total_amount, p.status, p.notes
		ORDER BY p.purchase_date DESC
	`)
	if err != nil {
		return ListPurchasesResponse{
			Message: "Failed to retrieve purchases",
			Data:    []PurchaseListItem{},
		}, errors.New("failed to retrieve purchases")
	}
	defer rows.Close()

	var purchases []PurchaseListItem
	for rows.Next() {
		var purchase PurchaseListItem
		var notes string
		err = rows.Scan(
			&purchase.ID,
			&purchase.Invoice,
			&purchase.Supplier,
			&purchase.OrderDate,
			&purchase.Total,
			&purchase.Status,
			&notes,
			&purchase.TotalItem,
		)
		if err != nil {
			return ListPurchasesResponse{Message: "Failed to scan purchase"}, errors.New("failed to scan purchase")
		}

		// Parse expected delivery from notes
		expectedDelivery := parseExpectedDelivery(notes)
		if expectedDelivery != nil {
			purchase.ExpectedDelivery = expectedDelivery
		}

		purchases = append(purchases, purchase)
	}

	if err = rows.Err(); err != nil {
		return ListPurchasesResponse{Message: "Error iterating purchases"}, errors.New("error iterating purchases: " + err.Error())
	}

	return ListPurchasesResponse{
		Message: "Purchases retrieved successfully",
		Data:    purchases,
	}, nil
}

// parseExpectedDelivery extracts the expected delivery date from notes
// Format: "Expected delivery: YYYY-MM-DD"
func parseExpectedDelivery(notes string) *time.Time {
	if notes == "" {
		return nil
	}

	// Look for "Expected delivery: YYYY-MM-DD" pattern
	re := regexp.MustCompile(`Expected delivery:\s*(\d{4}-\d{2}-\d{2})`)
	matches := re.FindStringSubmatch(notes)
	if len(matches) < 2 {
		return nil
	}

	// Try to parse the date
	dateStr := strings.TrimSpace(matches[1])
	parsedDate, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return nil
	}

	return &parsedDate
}

// UpdatePurchaseStatus updates the status of a purchase
//
//encore:api public method=PUT path=/api/purchases/:id/status
func UpdatePurchaseStatus(ctx context.Context, id uuid.UUID, req *UpdatePurchaseStatusRequest) (Response, error) {
	// Validate request
	if err := req.Validate(); err != nil {
		return Response{Message: "Validation failed"}, err
	}

	// Check if purchase exists
	var purchaseExists bool
	err := db.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM purchases WHERE id = $1)", id).Scan(&purchaseExists)
	if err != nil {
		return Response{Message: "Failed to check purchase"}, err
	}
	if !purchaseExists {
		return Response{Message: "Purchase not found"}, errors.New("purchase not found")
	}

	// Update purchase status
	result, err := db.Exec(ctx, `
		UPDATE purchases 
		SET status = $1, updated_at = NOW()
		WHERE id = $2
	`, req.Status, id)
	if err != nil {
		return Response{Message: "Failed to update purchase status"}, err
	}

	// Check if any rows were affected
	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return Response{Message: "Purchase not found"}, errors.New("purchase not found")
	}

	return Response{Message: "Purchase status updated successfully"}, nil
}
