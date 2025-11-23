package procurement

import (
	"errors"

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
	TaxID         string `json:"tax_id"`
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
	if len(s.TaxID) > 50 {
		return errors.New("tax_id must be less than 50 characters")
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
	TaxID         string    `json:"tax_id"`
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
