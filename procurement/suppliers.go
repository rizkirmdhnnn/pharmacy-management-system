package procurement

import (
	"context"
	"errors"
	"time"

	"encore.dev/types/uuid"
)

// Supplier model
type Supplier struct {
	ID            uuid.UUID `json:"id"`
	Name          string    `json:"name"`
	ContactPerson string    `json:"contact_person"`
	Email         string    `json:"email"`
	Phone         string    `json:"phone"`
	Address       string    `json:"address"`
	City          string    `json:"city"`
	Country       string    `json:"country"`
	IsActive      bool      `json:"is_active"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// CreateSupplier creates a new supplier
//
//encore:api public method=POST path=/api/suppliers
func CreateSupplier(ctx context.Context, req *CreateSupplierRequest) (Response, error) {
	// Check if supplier already exists
	exists, err := IsSupplierExists(ctx, req.Name)
	if err != nil {
		return Response{Message: "Failed to check if supplier already exists"}, err
	}
	if exists {
		return Response{Message: "Supplier with this name already exists"}, nil
	}

	// Create supplier
	var supplierID uuid.UUID
	err = db.QueryRow(ctx, `
		INSERT INTO suppliers (name, contact_person, email, phone, address, city, country, is_active) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8) 
		RETURNING id
	`, req.Name, req.ContactPerson, req.Email, req.Phone, req.Address, req.City, req.Country, true).Scan(&supplierID)
	if err != nil {
		return Response{Message: "Failed to create supplier"}, err
	}

	return Response{Message: "Supplier created successfully"}, nil
}

// GetAllSuppliers retrieves all suppliers
//
//encore:api public method=GET path=/api/suppliers
func GetAllSuppliers(ctx context.Context) (ListSuppliersResponse, error) {
	var listSuppliersResponse ListSuppliersResponse
	rows, err := db.Query(ctx, `
		SELECT id, name, contact_person, email, phone, address, city, country, is_active
		FROM suppliers 
		ORDER BY name
	`)
	if err != nil {
		return ListSuppliersResponse{
			Message: "Failed to retrieve suppliers",
			Data:    []SupplierListItem{},
		}, errors.New("failed to retrieve suppliers")
	}
	defer rows.Close()

	for rows.Next() {
		var supplier SupplierListItem
		err = rows.Scan(
			&supplier.ID,
			&supplier.Name,
			&supplier.ContactPerson,
			&supplier.Email,
			&supplier.Phone,
			&supplier.Address,
			&supplier.City,
			&supplier.Country,
			&supplier.IsActive,
		)
		if err != nil {
			return ListSuppliersResponse{Message: "Failed to scan supplier"}, errors.New("failed to scan supplier")
		}
		listSuppliersResponse.Data = append(listSuppliersResponse.Data, supplier)
	}

	if err = rows.Err(); err != nil {
		return ListSuppliersResponse{Message: "Error iterating suppliers"}, errors.New("error iterating suppliers: " + err.Error())
	}

	return ListSuppliersResponse{
		Message: "Suppliers retrieved successfully",
		Data:    listSuppliersResponse.Data,
	}, nil
}

// GetSupplier retrieves a supplier by ID
//
//encore:api public method=GET path=/api/suppliers/:id
func GetSupplier(ctx context.Context, id uuid.UUID) (SupplierResponse, error) {
	var supplier SupplierListItem
	err := db.QueryRow(ctx, `
		SELECT id, name, contact_person, email, phone, address, city, country, is_active
		FROM suppliers 
		WHERE id = $1
	`, id).Scan(
		&supplier.ID,
		&supplier.Name,
		&supplier.ContactPerson,
		&supplier.Email,
		&supplier.Phone,
		&supplier.Address,
		&supplier.City,
		&supplier.Country,
		&supplier.IsActive,
	)
	if err != nil {
		return SupplierResponse{Message: "Supplier not found"}, errors.New("supplier not found")
	}
	return SupplierResponse{
		Message: "Supplier retrieved successfully",
		Data:    &supplier,
	}, nil
}

// IsSupplierExists checks if a supplier with the given name already exists
func IsSupplierExists(ctx context.Context, name string) (bool, error) {
	var count int
	err := db.QueryRow(ctx, "SELECT COUNT(*) FROM suppliers WHERE name = $1", name).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
