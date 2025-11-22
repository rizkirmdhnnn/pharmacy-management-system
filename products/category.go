package products

import (
	"context"
	"strings"
	"time"

	"encore.dev/beta/errs"
)

var (
	ErrCategoryExists   = errs.B().Code(errs.AlreadyExists).Msg("category already exists").Err()
	ErrCategoryNotFound = errs.B().Code(errs.NotFound).Msg("category not found").Err()
	ErrInvalidInput     = errs.B().Code(errs.InvalidArgument).Msg("invalid input").Err()
)

// create category
// encore:api public path=/api/categories method=POST
func CreateCategory(ctx context.Context, request CreateCategoryRequest) (*Response, error) {
	// Validate input
	request.Name = strings.TrimSpace(request.Name)
	if request.Name == "" {
		return nil, errs.B().Code(errs.InvalidArgument).Msg("category name is required").Err()
	}

	// Check if category exists
	exists, err := checkIfCategoryExists(ctx, request.Name)
	if err != nil {
		return nil, errs.WrapCode(err, errs.Internal, "failed to check category existence")
	}
	if exists {
		return nil, ErrCategoryExists
	}

	// Create category
	now := time.Now()
	var category Category
	err = db.QueryRow(ctx, `
		INSERT INTO categories (name, description, created_at, updated_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id, name, description
	`, request.Name, request.Description, now, now).Scan(&category.ID, &category.Name, &category.Description)
	if err != nil {
		return nil, errs.WrapCode(err, errs.Internal, "failed to create category")
	}

	return &Response{
		Message: "Category created successfully",
	}, nil
}

// get all categories
// encore:api public path=/api/categories method=GET
func GetAllCategories(ctx context.Context) (*GetAllCategoriesResponse, error) {
	rows, err := db.Query(ctx, `
		SELECT id, name, description
		FROM categories
		ORDER BY name ASC
	`)
	if err != nil {
		return nil, errs.WrapCode(err, errs.Internal, "failed to fetch categories")
	}
	defer rows.Close()

	var categories []Category
	for rows.Next() {
		var cat Category
		err := rows.Scan(&cat.ID, &cat.Name, &cat.Description)
		if err != nil {
			return nil, errs.WrapCode(err, errs.Internal, "failed to scan category")
		}
		categories = append(categories, cat)
	}

	if err = rows.Err(); err != nil {
		return nil, errs.WrapCode(err, errs.Internal, "error iterating categories")
	}

	return &GetAllCategoriesResponse{
		Categories: categories,
	}, nil
}

// check if category exists
func checkIfCategoryExists(ctx context.Context, name string) (bool, error) {
	var exists bool
	err := db.QueryRow(ctx, `
		SELECT EXISTS(SELECT 1 FROM categories WHERE name = $1)
	`, name).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

// check if category exists by id
func checkIfCategoryExistsById(ctx context.Context, id string) (bool, error) {
	var exists bool
	err := db.QueryRow(ctx, `
		SELECT EXISTS(SELECT 1 FROM categories WHERE id = $1)
	`, id).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}
