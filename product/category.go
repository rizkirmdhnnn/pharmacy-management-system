package product

import (
	"context"
	"errors"
	"time"

	"encore.dev/types/uuid"
)

// Category model
type Category struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// GetAllCategories retrieves all categories
//
//encore:api public method=GET path=/api/categories
func GetAllCategories(ctx context.Context) (ListCategoriesResponse, error) {
	var listCategoriesResponse ListCategoriesResponse
	rows, err := db.Query(ctx, "SELECT id, name, description FROM categories")
	if err != nil {
		return ListCategoriesResponse{
			Message: "Failed to retrieve categories",
			Data:    []CategoryListItem{},
		}, errors.New("failed to retrieve categories")
	}
	defer rows.Close()
	for rows.Next() {
		var category CategoryListItem
		err = rows.Scan(
			&category.ID,
			&category.Name,
			&category.Description,
		)
		if err != nil {
			return ListCategoriesResponse{Message: "Failed to scan category"}, errors.New("failed to scan category")
		}
		listCategoriesResponse.Data = append(listCategoriesResponse.Data, category)
	}
	return ListCategoriesResponse{
		Message: "Categories retrieved successfully",
		Data:    listCategoriesResponse.Data,
	}, nil
}

// IsCategoryExists checks if a category with the given name already exists
func IsCategoryExists(ctx context.Context, name string) (bool, error) {
	var count int
	err := db.QueryRow(ctx, "SELECT COUNT(*) FROM categories WHERE name = $1", name).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
