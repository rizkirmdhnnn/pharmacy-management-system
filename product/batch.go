package product

import (
	"context"
	"time"

	"encore.dev/types/uuid"
)

type Batch struct {
	ID             uuid.UUID  `json:"id"`
	ProductID      uuid.UUID  `json:"product_id"`
	BatchNumber    string     `json:"batch_number"`
	Quantity       int        `json:"quantity"`
	PurchasePrice  float64    `json:"purchase_price"`
	SellingPrice   float64    `json:"selling_price"`
	ExpirationDate time.Time  `json:"expiration_date"`
	SupplierID     *uuid.UUID `json:"supplier_id,omitempty"`
	PurchaseID     *uuid.UUID `json:"purchase_id,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

// CreateBatch creates a new batch
func CreateBatch(ctx context.Context, batch *Batch) error {
	_, err := db.Exec(ctx, `
		INSERT INTO batches (product_id, batch_number, quantity, purchase_price, selling_price, expiration_date, supplier_id, purchase_id) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`, batch.ProductID, batch.BatchNumber, batch.Quantity, batch.PurchasePrice, batch.SellingPrice, batch.ExpirationDate, batch.SupplierID, batch.PurchaseID)
	if err != nil {
		return err
	}
	return nil
}
