package service

import (
	"sondth-test_soa/app/helper"
	"sondth-test_soa/app/repository"
)

type ServiceCollections struct {
	CategorySvc ICategoryService
	ProductSvc  IProductService
	UserService IUserService
	ReviewSvc   IReviewService
	WishlistSvc IWishlistService
}

func RegisterServices(helpers helper.HelperCollections, repositories repository.RepositoryCollections) ServiceCollections {
	return ServiceCollections{
		CategorySvc: NewCategoryService(repositories, helpers),
		ProductSvc:  NewProductService(repositories, helpers),
		UserService: NewUserService(repositories, helpers),
		ReviewSvc:   NewReviewService(repositories, helpers),
		WishlistSvc: NewWishlistService(repositories, helpers),
	}
}
