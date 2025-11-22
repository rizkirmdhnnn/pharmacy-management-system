package products

type Category struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}

type Product struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	CategoryID  string    `json:"category_id"`
	Category    *Category `json:"category,omitempty"`
	IsActive    bool      `json:"is_active"`
	BasePrice   float64   `json:"base_price"`
}
