package product

import (
	"errors"
	"time"

	"encore.dev/types/uuid"
)

type Response struct {
	Message string `json:"message"`
}

type CategoryListItem struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
}

type ListCategoriesResponse struct {
	Message string             `json:"message"`
	Data    []CategoryListItem `json:"data"`
}

type CreateProductRequest struct {
	Name                 string    `json:"name"`
	CategoryID           uuid.UUID `json:"category_id"`
	SellingPrice         float64   `json:"selling_price"`
	CostPrice            float64   `json:"cost_price"`
	StockQuantity        int       `json:"stock_quantity"`
	MinimumStockQuantity int       `json:"minimum_stock_quantity"`
	Barcode              string    `json:"barcode"`
	ExpirationDate       time.Time `json:"expiration_date"`
	BatchNumber          string    `json:"batch_number"`
	SupplierID           uuid.UUID `json:"supplier_id"`
	Description          string    `json:"description"`
}

func (p *CreateProductRequest) Validate() error {
	if p.Name == "" {
		return errors.New("name is required")
	}
	if len(p.Name) > 255 {
		return errors.New("name must be less than 255 characters")
	}
	if p.CategoryID == uuid.Nil {
		return errors.New("category_id is required")
	}
	if p.SellingPrice <= 0 {
		return errors.New("selling_price must be greater than 0")
	}
	if p.CostPrice <= 0 {
		return errors.New("cost_price must be greater than 0")
	}
	if p.StockQuantity < 0 {
		return errors.New("stock_quantity must be non-negative")
	}
	if p.MinimumStockQuantity < 0 {
		return errors.New("minimum_stock_quantity must be non-negative")
	}
	if p.ExpirationDate.IsZero() {
		return errors.New("expiration_date is required")
	}
	if p.BatchNumber == "" {
		return errors.New("batch_number is required")
	}
	if p.SupplierID == uuid.Nil {
		return errors.New("supplier_id is required")
	}
	return nil
}

type CreateBatchRequest struct {
	ProductID      uuid.UUID `json:"product_id"`
	BatchNumber    string    `json:"batch_number"`
	Quantity       int       `json:"quantity"`
	PurchasePrice  float64   `json:"purchase_price"`
	SellingPrice   float64   `json:"selling_price"`
	ExpirationDate time.Time `json:"expiration_date"`
}

type ProductWithBatchListItem struct {
	ID                   uuid.UUID `json:"id"`
	Name                 string    `json:"name"`
	Category             string    `json:"category"`
	Description          string    `json:"description"`
	BatchNumber          string    `json:"batch_number"`
	Quantity             int       `json:"quantity"`
	MinimumStockQuantity int       `json:"minimum_stock_quantity"`
	SellingPrice         float64   `json:"selling_price"`
	ExpirationDate       time.Time `json:"expiration_date"`
}

type ProductResponse struct {
	Message string                     `json:"message"`
	Data    []ProductWithBatchListItem `json:"data,omitempty"`
}
