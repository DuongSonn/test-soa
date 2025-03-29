package model

import (
	"sondth-test_soa/app/entity"
	"sondth-test_soa/app/repository"

	"github.com/google/uuid"
)

// CreateProductRequest struct
type CreateProductRequest struct {
	Name        string    `json:"name" validate:"required"`
	Description *string   `json:"description"`
	Price       float64   `json:"price" validate:"required"`
	Quantity    uint64    `json:"quantity" validate:"required"`
	CategoryID  uuid.UUID `json:"category_id" validate:"required"`
}
type CreateProductResponse struct{}

// GetProductRequest struct
type GetProductRequest struct {
	Name        *string            `json:"name"`
	CategoryIDs []uuid.UUID        `json:"category_ids"`
	Page        *int               `json:"page"`
	Limit       *int               `json:"limit"`
	Status      *string            `json:"status"`
	Order       repository.OrderBy `json:"order"`
}
type GetProductResponse struct {
	Count  int64            `json:"count"`
	Result []entity.Product `json:"result"`
}

// UpdateProductRequest struct
type UpdateProductRequest struct {
	ID uuid.UUID `json:"id" validate:"required"`
	CreateProductRequest
}
type UpdateProductResponse struct {
	Product entity.Product `json:"product"`
}

// DeleteProductRequest struct
type DeleteProductRequest struct {
	ID uuid.UUID `json:"id" validate:"required"`
}
type DeleteProductResponse struct{}
