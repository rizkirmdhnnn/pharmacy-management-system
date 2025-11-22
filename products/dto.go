package products

// CreateCategoryRequest represents the request for creating a category
type CreateCategoryRequest struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}

// GetAllCategoriesResponse represents the response for listing categories
type GetAllCategoriesResponse struct {
	Categories []Category `json:"categories"`
}

// CreateProductRequest represents the request for creating a product
type CreateProductRequest struct {
	Name        string  `json:"name"`
	Description string  `json:"description,omitempty"`
	CategoryID  string  `json:"category_id"`
	BasePrice   float64 `json:"base_price"`
}

// GetAllProductsResponse represents the response for listing products
type GetAllProductsResponse struct {
	Products []Product `json:"products"`
}

// Response represents the response for all endpoints
type Response struct {
	Message string `json:"message"`
}
