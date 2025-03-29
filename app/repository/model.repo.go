package repository

import (
	"github.com/google/uuid"
)

type OrderBy struct {
	Field string `json:"field"`
	Order string `json:"order"`
}
type Filter struct {
	Fields     []string
	OmitFields []string
}

type FindProductByFilter struct {
	Filter
	ID          *uuid.UUID
	Name        *string
	CategoryIDs []uuid.UUID
	Page        *int
	Limit       *int
	Status      *string
	Order       OrderBy

	// Relationship
	CategoryFields []string
}

type FindCategoryByFilter struct {
	Filter
	ID    *uuid.UUID
	Name  *string
	Page  *int
	Limit *int
}

type FindUserByFilter struct {
	Filter
	ID       *uuid.UUID
	Page     *int
	Limit    *int
	Role     *string
	Name     *string
	Username *string
}

type FindReviewByFilter struct {
	Filter
	ID          *uuid.UUID
	ProductID   *uuid.UUID
	ProductName *string
	UserID      *uuid.UUID
	Page        *int
	Limit       *int

	// Relationship
	ProductFields []string
	UserFields    []string
}

type FindWishlistByFilter struct {
	Filter
	ID        *uuid.UUID
	UserID    *uuid.UUID
	Page      *int
	Limit     *int
	ProductID *uuid.UUID

	// Relationship
	ProductFields []string
	UserFields    []string
}
