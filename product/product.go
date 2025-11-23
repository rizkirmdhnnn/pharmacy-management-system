package product

import (
	"context"
	"errors"
	"time"

	"encore.dev/types/uuid"
)

// Product model
type Product struct {
	ID            uuid.UUID `json:"id"`
	Name          string    `json:"name"`
	CategoryID    uuid.UUID `json:"category_id"`
	Description   string    `json:"description"`
	BasePrice     float64   `json:"base_price"`
	MinStockLevel int       `json:"min_stock_level"`
	Barcode       string    `json:"barcode"`
	IsActive      bool      `json:"is_active"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// CreateProduct creates a new product
//
//encore:api public method=POST path=/api/products
func CreateProduct(ctx context.Context, product *CreateProductRequest) (Response, error) {
	// Check if product already exists
	exists, err := IsProductExists(ctx, product.Name)
	if err != nil {
		return Response{Message: "Failed to check if product already exists"}, err
	}
	if exists {
		return Response{Message: "Product with this name already exists"}, nil
	}

	// Create product
	var productID uuid.UUID
	err = db.QueryRow(ctx, `
		INSERT INTO products (name, category_id, description, base_price, min_stock_level, barcode, is_active) 
		VALUES ($1, $2, $3, $4, $5, $6, $7) 
		RETURNING id
	`, product.Name, product.CategoryID, product.Description, product.SellingPrice, product.MinimumStockQuantity, product.Barcode, true).Scan(&productID)
	if err != nil {
		return Response{Message: "Failed to create product"}, err
	}

	// Create batch
	CreateBatch(ctx, &Batch{
		ProductID:      productID,
		BatchNumber:    product.BatchNumber,
		Quantity:       product.StockQuantity,
		PurchasePrice:  product.CostPrice,
		SellingPrice:   product.SellingPrice,
		ExpirationDate: product.ExpirationDate,
		SupplierID:     &product.SupplierID,
		PurchaseID:     nil, // Will be set when purchase is completed
	})

	return Response{Message: "Product created successfully"}, nil
}

// GetAllProducts retrieves all products with aggregated batch information
//
//encore:api public method=GET path=/api/products
func GetAllProducts(ctx context.Context) (*ProductResponse, error) {
	query := `
		SELECT 
			p.id,
			p.name,
			c.name as category_name,
			p.description,
			p.min_stock_level,
			COALESCE(SUM(b.quantity), 0) as total_quantity,
			(
				SELECT batch_number 
				FROM batches 
				WHERE product_id = p.id 
				ORDER BY created_at DESC 
				LIMIT 1
			) as latest_batch_number,
			(
				SELECT selling_price 
				FROM batches 
				WHERE product_id = p.id 
				ORDER BY created_at DESC 
				LIMIT 1
			) as latest_selling_price,
			(
				SELECT expiration_date 
				FROM batches 
				WHERE product_id = p.id 
				ORDER BY expiration_date ASC 
				LIMIT 1
			) as earliest_expiration_date
		FROM products p
		LEFT JOIN categories c ON p.category_id = c.id
		LEFT JOIN batches b ON p.id = b.product_id
		GROUP BY p.id, p.name, c.name, p.description, p.min_stock_level
		ORDER BY p.name
	`

	rows, err := db.Query(ctx, query)
	if err != nil {
		return nil, errors.New("failed to retrieve products: " + err.Error())
	}
	defer rows.Close()

	var products []ProductWithBatchListItem
	for rows.Next() {
		var product ProductWithBatchListItem
		var totalQuantity int
		var nullableBatchNumber *string
		var nullableSellingPrice *float64
		var nullableExpirationDate *time.Time

		err := rows.Scan(
			&product.ID,
			&product.Name,
			&product.Category,
			&product.Description,
			&product.MinimumStockQuantity,
			&totalQuantity,
			&nullableBatchNumber,
			&nullableSellingPrice,
			&nullableExpirationDate,
		)
		if err != nil {
			return nil, errors.New("failed to scan product: " + err.Error())
		}

		// Set aggregated and latest batch information
		product.Quantity = totalQuantity
		if nullableBatchNumber != nil {
			product.BatchNumber = *nullableBatchNumber
		}
		if nullableSellingPrice != nil {
			product.SellingPrice = *nullableSellingPrice
		}
		if nullableExpirationDate != nil {
			product.ExpirationDate = *nullableExpirationDate
		}

		products = append(products, product)
	}

	if err = rows.Err(); err != nil {
		return nil, errors.New("error iterating products: " + err.Error())
	}

	return &ProductResponse{Message: "Products retrieved successfully", Data: products}, nil
}

// GetProduct retrieves a product by ID
//
//encore:api public method=GET path=/api/products/:id
func GetProduct(ctx context.Context, id uuid.UUID) (*Product, error) {
	var product Product
	err := db.QueryRow(ctx, `
		SELECT id, name, category_id, description, base_price, min_stock_level, barcode, is_active, created_at, updated_at 
		FROM products 
		WHERE id = $1
	`, id).Scan(
		&product.ID,
		&product.Name,
		&product.CategoryID,
		&product.Description,
		&product.BasePrice,
		&product.MinStockLevel,
		&product.Barcode,
		&product.IsActive,
		&product.CreatedAt,
		&product.UpdatedAt,
	)
	if err != nil {
		return nil, errors.New("product not found")
	}
	return &product, nil
}

// IsProductExists checks if a product with the given name already exists
func IsProductExists(ctx context.Context, name string) (bool, error) {
	var count int
	err := db.QueryRow(ctx, "SELECT COUNT(*) FROM products WHERE name = $1", name).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
