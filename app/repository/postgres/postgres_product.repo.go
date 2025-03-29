package postgres

import (
	"context"

	"gorm.io/gorm"

	"sondth-test_soa/app/entity"
	"sondth-test_soa/app/repository"
)

type productRepository struct {
	db *gorm.DB
}

func NewPostgresProductRepository(db *gorm.DB) repository.IProductRepository {
	return &productRepository{
		db,
	}
}

func (r *productRepository) Create(
	ctx context.Context,
	tx *gorm.DB,
	data *entity.Product,
) error {
	if tx != nil {
		return tx.WithContext(ctx).Create(&data).Error
	}

	return r.db.WithContext(ctx).Create(&data).Error
}

func (r *productRepository) Update(
	ctx context.Context,
	tx *gorm.DB,
	data *entity.Product,
) error {
	if tx != nil {
		return tx.WithContext(ctx).Save(&data).Error
	}

	return r.db.WithContext(ctx).Save(&data).Error
}

func (r *productRepository) Delete(
	ctx context.Context,
	tx *gorm.DB,
	data *entity.Product,
) error {
	if tx != nil {
		return tx.WithContext(ctx).Delete(&data).Error
	}

	return r.db.WithContext(ctx).Delete(&data).Error
}

func (r *productRepository) FindOneByFilter(
	ctx context.Context,
	tx *gorm.DB,
	filter *repository.FindProductByFilter,
) (*entity.Product, error) {
	var product entity.Product
	err := r.buildFilter(ctx, tx, filter).First(&product).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *productRepository) FindManyByFilter(
	ctx context.Context,
	tx *gorm.DB,
	filter *repository.FindProductByFilter,
) ([]entity.Product, error) {
	var products []entity.Product

	query := r.buildFilter(ctx, tx, filter)
	if filter.Page != nil && filter.Limit != nil {
		offset := (*filter.Page - 1) * *filter.Limit
		query = query.Offset(offset).Limit(*filter.Limit)
	}

	err := query.Find(&products).Error
	return products, err
}

func (r *productRepository) CountByFilter(
	ctx context.Context,
	tx *gorm.DB,
	filter *repository.FindProductByFilter,
) (int64, error) {
	var count int64
	err := r.buildFilter(ctx, tx, filter).Count(&count).Error
	return count, err
}

// -------------------------------------------------------------------------------
func (r *productRepository) buildFilter(
	ctx context.Context,
	tx *gorm.DB,
	filter *repository.FindProductByFilter,
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

	if filter.ID != nil {
		query = query.Where("id = ?", filter.ID)
	}

	if filter.Name != nil {
		query = query.Where("name ILIKE ?", "%"+*filter.Name+"%")
	}

	if len(filter.CategoryIDs) > 0 {
		query = query.Where("category_id IN ?", filter.CategoryIDs)
	}

	if filter.Status != nil {
		query = query.Where("status = ?", filter.Status)
	}

	if filter.Order.Field != "" {
		query = query.Order(filter.Order.Field + " " + filter.Order.Order)
	}

	if len(filter.CategoryFields) > 0 {
		query = query.Model(&entity.Product{}).Preload("Category", func(db *gorm.DB) *gorm.DB {
			return db.Select(filter.CategoryFields)
		})
	}

	return query
}
