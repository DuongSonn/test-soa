package repository

import (
	"context"
	"sondth-test_soa/app/entity"

	"gorm.io/gorm"
)

type RepositoryCollections struct {
	ProductRepo  IProductRepository
	CategoryRepo ICategoryRepository
	ReviewRepo   IReviewRepository
	WishlistRepo IWishlistRepository
	UserRepo     IUserRepository
}

type IProductRepository interface {
	Create(ctx context.Context, tx *gorm.DB, data *entity.Product) error
	Update(ctx context.Context, tx *gorm.DB, data *entity.Product) error
	Delete(ctx context.Context, tx *gorm.DB, data *entity.Product) error
	FindOneByFilter(ctx context.Context, tx *gorm.DB, filter *FindProductByFilter) (*entity.Product, error)
	FindManyByFilter(ctx context.Context, tx *gorm.DB, filter *FindProductByFilter) ([]entity.Product, error)
	CountByFilter(ctx context.Context, tx *gorm.DB, filter *FindProductByFilter) (int64, error)
}

type ICategoryRepository interface {
	Create(ctx context.Context, tx *gorm.DB, data *entity.Category) error
	FindManyByFilter(ctx context.Context, tx *gorm.DB, filter *FindCategoryByFilter) ([]entity.Category, error)
	FindOneByFilter(ctx context.Context, tx *gorm.DB, filter *FindCategoryByFilter) (*entity.Category, error)
	GetCategorySummary(ctx context.Context, tx *gorm.DB) ([]entity.Category, error)
}

type IUserRepository interface {
	Create(ctx context.Context, tx *gorm.DB, data *entity.User) error
	FindOneByFilter(ctx context.Context, tx *gorm.DB, filter *FindUserByFilter) (*entity.User, error)
	FindManyByFilter(ctx context.Context, tx *gorm.DB, filter *FindUserByFilter) ([]entity.User, error)
	CountByFilter(ctx context.Context, tx *gorm.DB, filter *FindUserByFilter) (int64, error)
	Update(ctx context.Context, tx *gorm.DB, data *entity.User) error
}

type IReviewRepository interface {
	Create(ctx context.Context, tx *gorm.DB, data *entity.Review) error
	Delete(ctx context.Context, tx *gorm.DB, data *entity.Review) error
	FindManyByFilter(ctx context.Context, tx *gorm.DB, filter *FindReviewByFilter) ([]entity.Review, error)
	CountByFilter(ctx context.Context, tx *gorm.DB, filter *FindReviewByFilter) (int64, error)
	FindOneByFilter(ctx context.Context, tx *gorm.DB, filter *FindReviewByFilter) (*entity.Review, error)
}

type IWishlistRepository interface {
	Create(ctx context.Context, tx *gorm.DB, data *entity.Wishlist) error
	Delete(ctx context.Context, tx *gorm.DB, data *entity.Wishlist) error
	FindManyByFilter(ctx context.Context, tx *gorm.DB, filter *FindWishlistByFilter) ([]entity.Wishlist, error)
	CountByFilter(ctx context.Context, tx *gorm.DB, filter *FindWishlistByFilter) (int64, error)
	FindOneByFilter(ctx context.Context, tx *gorm.DB, filter *FindWishlistByFilter) (*entity.Wishlist, error)
}
