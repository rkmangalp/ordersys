package product

import (
	"context"
	"errors"
	"fmt"
)

// Service defines the business logic for products
type Service struct {
	repo Repository
}

// NewService creates a new product service
func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

// CreateProduct creates a new product
func (s *Service) CreateProduct(ctx context.Context, req CreateProductRequest) (*Product, error) {
	// Validate request
	if err := req.Validate(); err != nil {
		return nil, err
	}

	// Check if SKU already exists
	existing, err := s.repo.GetBySKU(ctx, req.SKU)
	if err == nil && existing != nil {
		return nil, ErrSKUAlreadyExists
	} else if err != nil && !isNotFoundError(err) {
		return nil, fmt.Errorf("failed to check SKU: %w", err)
	}

	// Create product
	p := &Product{
		SKU:         req.SKU,
		Name:        req.Name,
		Description: req.Description,
		PriceCents:  req.PriceCents,
		Category:    req.Category,
	}

	if err := s.repo.Create(ctx, p); err != nil {
		return nil, err
	}

	return p, nil
}

// GetProduct retrieves a product by ID
func (s *Service) GetProduct(ctx context.Context, id string) (*Product, error) {
	if id == "" {
		return nil, ErrProductNotFound
	}
	return s.repo.GetByID(ctx, id)
}

// ListProducts retrieves products with pagination
func (s *Service) ListProducts(ctx context.Context, page, pageSize int) (*ProductListResponse, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	products, total, err := s.repo.List(ctx, page, pageSize)
	if err != nil {
		return nil, err
	}

	totalPages := (total + pageSize - 1) / pageSize

	return &ProductListResponse{
		Data: products,
		Meta: Meta{
			Page:       page,
			PageSize:   pageSize,
			Total:      total,
			TotalPages: totalPages,
		},
	}, nil
}

// UpdateProduct updates a product
func (s *Service) UpdateProduct(ctx context.Context, id string, req UpdateProductRequest) (*Product, error) {
	if id == "" {
		return nil, ErrProductNotFound
	}

	// Validate request
	if err := req.Validate(); err != nil {
		return nil, err
	}

	// Get existing product
	existing, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Apply updates
	if req.SKU != nil {
		// Check if new SKU conflicts
		if *req.SKU != existing.SKU {
			conflict, err := s.repo.GetBySKU(ctx, *req.SKU)
			if err == nil && conflict != nil {
				return nil, ErrSKUAlreadyExists
			} else if err != nil && !isNotFoundError(err) {
				return nil, fmt.Errorf("failed to check SKU: %w", err)
			}
		}
		existing.SKU = *req.SKU
	}
	if req.Name != nil {
		existing.Name = *req.Name
	}
	if req.Description != nil {
		existing.Description = req.Description
	}
	if req.PriceCents != nil {
		existing.PriceCents = *req.PriceCents
	}
	if req.Category != nil {
		existing.Category = *req.Category
	}

	// Save updates
	if err := s.repo.Update(ctx, existing); err != nil {
		return nil, err
	}

	return existing, nil
}

// DeleteProduct deletes a product
func (s *Service) DeleteProduct(ctx context.Context, id string) error {
	if id == "" {
		return ErrProductNotFound
	}
	return s.repo.Delete(ctx, id)
}

// isNotFoundError checks if error is a NotFoundError
func isNotFoundError(err error) bool {
	var notFound *NotFoundError
	return errors.As(err, &notFound)
}