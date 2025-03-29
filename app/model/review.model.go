package model

import (
	"sondth-test_soa/app/entity"

	"github.com/google/uuid"
)

// CreateReviewRequest struct
type CreateReviewRequest struct {
	ProductID uuid.UUID `json:"product_id"`
	Rating    float64   `json:"rating"`
	Comment   string    `json:"comment"`
}
type CreateReviewResponse struct {
}

// DeleteReviewRequest struct
type DeleteReviewRequest struct {
	ReviewID uuid.UUID `json:"review_id"`
}
type DeleteReviewResponse struct{}

// GetReviewsRequest struct
type GetReviewsRequest struct {
	ProductName *string `json:"product_name"`
	Page        *int    `json:"page"`
	Limit       *int    `json:"limit"`
}
type GetReviewsResponse struct {
	Reviews []entity.Review `json:"reviews"`
	Count   int64           `json:"count"`
}

// GetReviewsSummaryRequest struct
type GetReviewsSummaryRequest struct {
}
type GetReviewsSummaryResponse struct {
	Count int64 `json:"count"`
}
