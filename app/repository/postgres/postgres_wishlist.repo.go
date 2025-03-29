package postgres

import (
	"context"

	"gorm.io/gorm"

	"sondth-test_soa/app/entity"
	"sondth-test_soa/app/repository"
)

type wishlistRepository struct {
	db *gorm.DB
}

func NewPostgresWishlistRepository(db *gorm.DB) repository.IWishlistRepository {
	return &wishlistRepository{
		db,
	}
}

func (r *wishlistRepository) Create(
	ctx context.Context,
	tx *gorm.DB,
	data *entity.Wishlist,
) error {
	if tx != nil {
		return tx.WithContext(ctx).Create(&data).Error
	}

	return r.db.WithContext(ctx).Create(&data).Error
}

func (r *wishlistRepository) Delete(
	ctx context.Context,
	tx *gorm.DB,
	data *entity.Wishlist,
) error {
	if tx != nil {
		return tx.WithContext(ctx).Delete(&data).Error
	}

	return r.db.WithContext(ctx).Delete(&data).Error
}

func (r *wishlistRepository) FindManyByFilter(
	ctx context.Context,
	tx *gorm.DB,
	filter *repository.FindWishlistByFilter,
) ([]entity.Wishlist, error) {
	var wishlists []entity.Wishlist

	query := r.buildFilter(ctx, tx, filter)
	if filter.Page != nil && filter.Limit != nil {
		offset := (*filter.Page - 1) * *filter.Limit
		query = query.Offset(offset).Limit(*filter.Limit)
	}

	err := query.Find(&wishlists).Error
	return wishlists, err
}

func (r *wishlistRepository) CountByFilter(
	ctx context.Context,
	tx *gorm.DB,
	filter *repository.FindWishlistByFilter,
) (int64, error) {
	var count int64
	err := r.buildFilter(ctx, tx, filter).Count(&count).Error
	return count, err
}

func (r *wishlistRepository) FindOneByFilter(
	ctx context.Context,
	tx *gorm.DB,
	filter *repository.FindWishlistByFilter,
) (*entity.Wishlist, error) {
	var wishlist entity.Wishlist
	err := r.buildFilter(ctx, tx, filter).First(&wishlist).Error
	return &wishlist, err
}

// -------------------------------------------------------------------------------
func (r *wishlistRepository) buildFilter(
	ctx context.Context,
	tx *gorm.DB,
	filter *repository.FindWishlistByFilter,
) *gorm.DB {
	query := r.db.WithContext(ctx)
	if tx != nil {
		query = tx
	}

	if len(filter.OmitFields) > 0 {
		query = query.Omit(filter.OmitFields...)
	} else {
		query = query.Select(filter.Fields)
	}

	if filter.UserID != nil {
		query = query.Where("user_id = ?", filter.UserID)
	}

	if len(filter.ProductFields) > 0 {
		query = query.Model(&entity.Wishlist{}).Preload("Product", func(db *gorm.DB) *gorm.DB {
			return db.Select(filter.ProductFields)
		})
	}

	if len(filter.UserFields) > 0 {
		query = query.Model(&entity.Wishlist{}).Preload("User", func(db *gorm.DB) *gorm.DB {
			return db.Select(filter.UserFields)
		})
	}

	if filter.ProductID != nil {
		query = query.Where("product_id = ?", filter.ProductID)
	}

	query = query.Order("created_at DESC")

	return query
}
