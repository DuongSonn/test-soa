package postgres

import (
	"context"

	"gorm.io/gorm"

	"sondth-test_soa/app/entity"
	"sondth-test_soa/app/repository"
)

type reviewRepository struct {
	db *gorm.DB
}

func NewPostgresReviewRepository(db *gorm.DB) repository.IReviewRepository {
	return &reviewRepository{
		db,
	}
}

func (r *reviewRepository) Create(
	ctx context.Context,
	tx *gorm.DB,
	data *entity.Review,
) error {
	if tx != nil {
		return tx.WithContext(ctx).Create(&data).Error
	}

	return r.db.WithContext(ctx).Create(&data).Error
}

func (r *reviewRepository) Delete(
	ctx context.Context,
	tx *gorm.DB,
	data *entity.Review,
) error {
	if tx != nil {
		return tx.WithContext(ctx).Delete(&data).Error
	}

	return r.db.WithContext(ctx).Delete(&data).Error
}

func (r *reviewRepository) FindManyByFilter(
	ctx context.Context,
	tx *gorm.DB,
	filter *repository.FindReviewByFilter,
) ([]entity.Review, error) {
	var reviews []entity.Review

	query := r.buildFilter(ctx, tx, filter)
	if filter.Page != nil && filter.Limit != nil {
		offset := (*filter.Page - 1) * *filter.Limit
		query = query.Offset(offset).Limit(*filter.Limit)
	}

	err := query.Find(&reviews).Error
	return reviews, err
}

func (r *reviewRepository) CountByFilter(
	ctx context.Context,
	tx *gorm.DB,
	filter *repository.FindReviewByFilter,
) (int64, error) {
	var count int64
	err := r.buildFilter(ctx, tx, filter).Count(&count).Error
	return count, err
}

func (r *reviewRepository) FindOneByFilter(
	ctx context.Context,
	tx *gorm.DB,
	filter *repository.FindReviewByFilter,
) (*entity.Review, error) {
	var review entity.Review
	err := r.buildFilter(ctx, tx, filter).First(&review).Error
	return &review, err
}

// -------------------------------------------------------------------------------
func (r *reviewRepository) buildFilter(
	ctx context.Context,
	tx *gorm.DB,
	filter *repository.FindReviewByFilter,
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

	if len(filter.ProductFields) > 0 {
		query = query.Model(&entity.Review{}).Preload("Product", func(db *gorm.DB) *gorm.DB {
			return db.Select(filter.ProductFields)
		})
	}
	if len(filter.UserFields) > 0 {
		query = query.Model(&entity.Review{}).Preload("User", func(db *gorm.DB) *gorm.DB {
			return db.Select(filter.UserFields)
		})
	}

	if filter.UserID != nil {
		query = query.Where("user_id = ?", filter.UserID)
	}

	if filter.ProductName != nil {
		query = query.Joins("Product").Where("products.name ILIKE ?", *filter.ProductName)
	}

	query = query.Order("reviews.created_at DESC")

	return query
}
