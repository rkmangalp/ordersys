package product

import (
	"time"
)

// Product represents a product in the system
type Product struct {
	ID          string    `json:"id"`
	SKU         string    `json:"sku"`
	Name        string    `json:"name"`
	Description *string   `json:"description,omitempty"`
	PriceCents  int       `json:"price_cents"`
	Category    string    `json:"category"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// CreateProductRequest represents the request to create a product
type CreateProductRequest struct {
	SKU         string  `json:"sku" validate:"required,max=100"`
	Name        string  `json:"name" validate:"required,max=255"`
	Description *string `json:"description,omitempty"`
	PriceCents  int     `json:"price_cents" validate:"required,gt=0"`
	Category    string  `json:"category" validate:"required,max=100"`
}

// UpdateProductRequest represents the request to update a product
type UpdateProductRequest struct {
	SKU         *string `json:"sku,omitempty" validate:"omitempty,max=100"`
	Name        *string `json:"name,omitempty" validate:"omitempty,max=255"`
	Description *string `json:"description,omitempty"`
	PriceCents  *int    `json:"price_cents,omitempty" validate:"omitempty,gt=0"`
	Category    *string `json:"category,omitempty" validate:"omitempty,max=100"`
}

// ProductListResponse represents a paginated list of products
type ProductListResponse struct {
	Data []Product `json:"data"`
	Meta Meta      `json:"meta"`
}

// Meta contains pagination metadata
type Meta struct {
	Page       int `json:"page"`
	PageSize   int `json:"page_size"`
	Total      int `json:"total"`
	TotalPages int `json:"total_pages"`
}

// Validate validates the create request
func (r *CreateProductRequest) Validate() error {
	if r.SKU == "" {
		return ErrInvalidSKU
	}
	if len(r.SKU) > 100 {
		return ErrSKUTooLong
	}
	if r.Name == "" {
		return ErrInvalidName
	}
	if len(r.Name) > 255 {
		return ErrNameTooLong
	}
	if r.PriceCents <= 0 {
		return ErrInvalidPrice
	}
	if r.Category == "" {
		return ErrInvalidCategory
	}
	if len(r.Category) > 100 {
		return ErrCategoryTooLong
	}
	return nil
}

// Validate validates the update request
func (r *UpdateProductRequest) Validate() error {
	if r.SKU != nil && *r.SKU == "" {
		return ErrInvalidSKU
	}
	if r.SKU != nil && len(*r.SKU) > 100 {
		return ErrSKUTooLong
	}
	if r.Name != nil && *r.Name == "" {
		return ErrInvalidName
	}
	if r.Name != nil && len(*r.Name) > 255 {
		return ErrNameTooLong
	}
	if r.PriceCents != nil && *r.PriceCents <= 0 {
		return ErrInvalidPrice
	}
	if r.Category != nil && *r.Category == "" {
		return ErrInvalidCategory
	}
	if r.Category != nil && len(*r.Category) > 100 {
		return ErrCategoryTooLong
	}
	return nil
}

// Domain errors
var (
	ErrInvalidSKU        = &ValidationError{Field: "sku", Message: "sku is required"}
	ErrSKUTooLong        = &ValidationError{Field: "sku", Message: "sku must be <= 100 characters"}
	ErrInvalidName       = &ValidationError{Field: "name", Message: "name is required"}
	ErrNameTooLong       = &ValidationError{Field: "name", Message: "name must be <= 255 characters"}
	ErrInvalidPrice      = &ValidationError{Field: "price_cents", Message: "price_cents must be > 0"}
	ErrInvalidCategory   = &ValidationError{Field: "category", Message: "category is required"}
	ErrCategoryTooLong   = &ValidationError{Field: "category", Message: "category must be <= 100 characters"}
	ErrProductNotFound   = &NotFoundError{Resource: "product"}
	ErrSKUAlreadyExists  = &ConflictError{Field: "sku", Message: "sku already exists"}
)

// ValidationError represents a validation error
type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return e.Field + ": " + e.Message
}

// NotFoundError represents a not found error
type NotFoundError struct {
	Resource string
}

func (e *NotFoundError) Error() string {
	return e.Resource + " not found"
}

// ConflictError represents a conflict error
type ConflictError struct {
	Field   string
	Message string
}

func (e *ConflictError) Error() string {
	return e.Field + ": " + e.Message
}