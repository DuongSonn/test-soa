package service

import (
	"context"

	"sondth-test_soa/app/model"
)

type IProductService interface {
	Create(ctx context.Context, req *model.CreateProductRequest) (*model.CreateProductResponse, error)
	Update(ctx context.Context, req *model.UpdateProductRequest) (*model.UpdateProductResponse, error)
	Delete(ctx context.Context, req *model.DeleteProductRequest) (*model.DeleteProductResponse, error)
	GetProducts(ctx context.Context, req *model.GetProductRequest) (*model.GetProductResponse, error)
}

type ICategoryService interface {
	Create(ctx context.Context, req *model.CreateCategoryRequest) (*model.CreateCategoryResponse, error)
	GetCategories(ctx context.Context, req *model.GetCategoriesRequest) (*model.GetCategoriesResponse, error)
	GetCategoriesSummary(ctx context.Context, req *model.GetCategoriesSummaryRequest) (*model.GetCategoriesSummaryResponse, error)
}

type IUserService interface {
	Register(ctx context.Context, req *model.UserRegisterRequest) (*model.UserRegisterResponse, error)
	Login(ctx context.Context, req *model.UserLoginRequest) (*model.UserLoginResponse, error)
	ForgetPassword(ctx context.Context, req *model.ForgetPasswordRequest) (*model.ForgetPasswordResponse, error)
	ChangePassword(ctx context.Context, req *model.ChangeUserPasswordRequest) (*model.ChangeUserPasswordResponse, error)
	UpdateUser(ctx context.Context, req *model.UpdateUserRequest) (*model.UpdateUserResponse, error)
	GetUsers(ctx context.Context, req *model.GetUsersRequest) (*model.GetUsersResponse, error)
}

type IWishlistService interface {
	AddToWishlist(ctx context.Context, req *model.AddToWishlistRequest) (*model.AddToWishlistResponse, error)
	GetWishlists(ctx context.Context, req *model.GetWishlistsRequest) (*model.GetWishlistsResponse, error)
	DeleteFromWishlist(ctx context.Context, req *model.DeleteFromWishlistRequest) (*model.DeleteFromWishlistResponse, error)
	GetWishlistsSummary(ctx context.Context, req *model.GetWishlistsSummaryRequest) (*model.GetWishlistsSummaryResponse, error)
}

type IReviewService interface {
	Create(ctx context.Context, req *model.CreateReviewRequest) (*model.CreateReviewResponse, error)
	GetReviews(ctx context.Context, req *model.GetReviewsRequest) (*model.GetReviewsResponse, error)
	Delete(ctx context.Context, req *model.DeleteReviewRequest) (*model.DeleteReviewResponse, error)
	GetReviewsSummary(ctx context.Context, req *model.GetReviewsSummaryRequest) (*model.GetReviewsSummaryResponse, error)
}
