package procurement

import (
	"errors"
	"fmt"
	"time"

	"encore.dev/types/uuid"
)

type Response struct {
	Message string `json:"message"`
}

type CreateSupplierRequest struct {
	Name          string `json:"name"`
	ContactPerson string `json:"contact_person"`
	Email         string `json:"email"`
	Phone         string `json:"phone"`
	Address       string `json:"address"`
	City          string `json:"city"`
	Country       string `json:"country"`
}

func (s *CreateSupplierRequest) Validate() error {
	if s.Name == "" {
		return errors.New("name is required")
	}
	if len(s.Name) > 200 {
		return errors.New("name must be less than 200 characters")
	}
	if len(s.ContactPerson) > 100 {
		return errors.New("contact_person must be less than 100 characters")
	}
	if len(s.Email) > 100 {
		return errors.New("email must be less than 100 characters")
	}
	if len(s.Phone) > 20 {
		return errors.New("phone must be less than 20 characters")
	}
	if len(s.City) > 100 {
		return errors.New("city must be less than 100 characters")
	}
	if len(s.Country) > 100 {
		return errors.New("country must be less than 100 characters")
	}
	return nil
}

type SupplierListItem struct {
	ID            uuid.UUID `json:"id"`
	Name          string    `json:"name"`
	ContactPerson string    `json:"contact_person"`
	Email         string    `json:"email"`
	Phone         string    `json:"phone"`
	Address       string    `json:"address"`
	City          string    `json:"city"`
	Country       string    `json:"country"`
	IsActive      bool      `json:"is_active"`
}

type SupplierResponse struct {
	Message string            `json:"message"`
	Data    *SupplierListItem `json:"data,omitempty"`
}

type ListSuppliersResponse struct {
	Message string             `json:"message"`
	Data    []SupplierListItem `json:"data"`
}

type PurchaseItemRequest struct {
	ProductID uuid.UUID `json:"product_id"`
	Quantity  int       `json:"quantity"`
}

type CreatePurchaseRequest struct {
	SupplierID       uuid.UUID             `json:"supplier_id"`
	OrderDate        time.Time             `json:"order_date"`
	InvoiceNumber    string                `json:"invoice_number"`
	ExpectedDelivery time.Time             `json:"expected_delivery"`
	Notes            string                `json:"notes"`
	Items            []PurchaseItemRequest `json:"items"`
	CreatedBy        uuid.UUID             `json:"created_by"`
}

func (p *CreatePurchaseRequest) Validate() error {
	if p.SupplierID == uuid.Nil {
		return errors.New("supplier_id is required")
	}
	if p.InvoiceNumber == "" {
		return errors.New("invoice_number is required")
	}
	if len(p.InvoiceNumber) > 50 {
		return errors.New("invoice_number must be less than 50 characters")
	}
	if p.OrderDate.IsZero() {
		return errors.New("order_date is required")
	}
	if len(p.Items) == 0 {
		return errors.New("at least one item is required")
	}
	if p.CreatedBy == uuid.Nil {
		return errors.New("created_by is required")
	}

	for i, item := range p.Items {
		itemNum := i + 1
		if item.ProductID == uuid.Nil {
			return fmt.Errorf("product_id is required for item %d", itemNum)
		}
		if item.Quantity <= 0 {
			return fmt.Errorf("quantity must be greater than 0 for item %d", itemNum)
		}
	}

	return nil
}

type PurchaseListItem struct {
	ID               uuid.UUID  `json:"id"`
	Invoice          string     `json:"invoice"`
	Supplier         string     `json:"supplier"`
	OrderDate        time.Time  `json:"order_date"`
	ExpectedDelivery *time.Time `json:"expected_delivery,omitempty"`
	Total            float64    `json:"total"`
	Status           string     `json:"status"`
	TotalItem        int        `json:"total_item"`
}

type ListPurchasesResponse struct {
	Message string             `json:"message"`
	Data    []PurchaseListItem `json:"data"`
}

type UpdatePurchaseStatusRequest struct {
	Status string `json:"status"`
}

func (u *UpdatePurchaseStatusRequest) Validate() error {
	if u.Status == "" {
		return errors.New("status is required")
	}
	validStatuses := map[string]bool{
		"pending":   true,
		"completed": true,
		"cancelled": true,
	}
	if !validStatuses[u.Status] {
		return errors.New("status must be one of: pending, completed, cancelled")
	}
	return nil
}
