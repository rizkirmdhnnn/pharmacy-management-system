package products

import (
	"context"
	"math"
	"strings"
	"time"

	"encore.dev/beta/errs"
)

var (
	ErrProductExists             = errs.B().Code(errs.AlreadyExists).Msg("product already exists").Err()
	ErrProductNotFound           = errs.B().Code(errs.NotFound).Msg("product not found").Err()
	ErrProductNameRequired       = errs.B().Code(errs.InvalidArgument).Msg("product name is required").Err()
	ErrProductCategoryIDRequired = errs.B().Code(errs.InvalidArgument).Msg("product category id is required").Err()
	ErrProductBasePriceRequired  = errs.B().Code(errs.InvalidArgument).Msg("product base price is required").Err()
)

// create new product
// encore:api public path=/api/products method=POST
func CreateProduct(ctx context.Context, request CreateProductRequest) (*Response, error) {
	// Validate input
	request.Name = strings.TrimSpace(request.Name)
	if request.Name == "" {
		return nil, ErrProductNameRequired
	}
	request.CategoryID = strings.TrimSpace(request.CategoryID)
	if request.CategoryID == "" {
		return nil, ErrProductCategoryIDRequired
	}
	request.BasePrice = math.Abs(request.BasePrice)
	if request.BasePrice <= 0 {
		return nil, ErrProductBasePriceRequired
	}

	// Check if category exists
	exists, err := checkIfCategoryExistsById(ctx, request.CategoryID)
	if err != nil {
		return nil, errs.WrapCode(err, errs.Internal, "failed to check category existence")
	}
	if !exists {
		return nil, ErrCategoryNotFound
	}

	// Check if product exists
	exists, err = checkIfProductExists(ctx, request.Name)
	if err != nil {
		return nil, errs.WrapCode(err, errs.Internal, "failed to check product existence")
	}
	if exists {
		return nil, ErrProductExists
	}

	// Create product
	now := time.Now()
	var product Product
	err = db.QueryRow(ctx, `
		INSERT INTO products (name, description, category_id, is_active, base_price, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, name, description, category_id, is_active, base_price
	`, request.Name, request.Description, request.CategoryID, true, request.BasePrice, now, now).Scan(&product.ID, &product.Name, &product.Description, &product.CategoryID, &product.IsActive, &product.BasePrice)
	if err != nil {
		return nil, errs.WrapCode(err, errs.Internal, "failed to create product")
	}

	return &Response{
		Message: "Product created successfully",
	}, nil
}

// check if product exists
func checkIfProductExists(ctx context.Context, name string) (bool, error) {
	var exists bool
	err := db.QueryRow(ctx, `
		SELECT EXISTS(SELECT 1 FROM products WHERE name = $1)
	`, name).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}
