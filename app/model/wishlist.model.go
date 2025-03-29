package model

import (
	"sondth-test_soa/app/entity"

	"github.com/google/uuid"
)

// AddToWishlistRequest struct
type AddToWishlistRequest struct {
	ProductID uuid.UUID `json:"product_id"`
}
type AddToWishlistResponse struct {
}

// RemoveFromWishlistRequest struct
type DeleteFromWishlistRequest struct {
	ID uuid.UUID `param:"id" validate:"required"`
}
type DeleteFromWishlistResponse struct {
}

// GetWishlistRequest struct
type GetWishlistsRequest struct {
	Page  int `json:"page"`
	Limit int `json:"limit"`
}
type GetWishlistsResponse struct {
	Wishlists []entity.Wishlist `json:"wishlists"`
	Count     int64             `json:"count"`
}

// GetWishlistsSummaryRequest struct
type GetWishlistsSummaryRequest struct {
}
type GetWishlistsSummaryResponse struct {
	Count int64 `json:"count"`
}
