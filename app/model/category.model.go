package model

import "sondth-test_soa/app/entity"

// CreateCategoryRequest struct
type CreateCategoryRequest struct {
	Name string `json:"name" validate:"required"`
}
type CreateCategoryResponse struct{}

// GetCategoriesRequest struct
type GetCategoriesRequest struct {
	Page  *int `json:"page"`
	Limit *int `json:"limit"`
}
type GetCategoriesResponse struct {
	Count  int64             `json:"count"`
	Result []entity.Category `json:"result"`
}

// GetCategoriesSummaryRequest struct
type GetCategoriesSummaryRequest struct{}
type GetCategoriesSummaryResponse struct {
	Categories []entity.Category `json:"categories"`
}
