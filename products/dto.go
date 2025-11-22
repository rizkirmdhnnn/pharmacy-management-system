package products

type CreateCategoryRequest struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}

// CreateCategoryResponse represents the response for creating a category
type CreateCategoryResponse struct {
	Category Category `json:"category"`
}

// GetAllCategoriesResponse represents the response for listing categories
type GetAllCategoriesResponse struct {
	Categories []Category `json:"categories"`
}
