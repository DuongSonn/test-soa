package postgres

import (
	"context"

	"gorm.io/gorm"

	"sondth-test_soa/app/entity"
	"sondth-test_soa/app/repository"
)

type categoryRepository struct {
	db *gorm.DB
}

func NewPostgresCategoryRepository(db *gorm.DB) repository.ICategoryRepository {
	return &categoryRepository{
		db,
	}
}

func (r *categoryRepository) Create(
	ctx context.Context,
	tx *gorm.DB,
	data *entity.Category,
) error {
	if tx != nil {
		return tx.WithContext(ctx).Create(&data).Error
	}

	return r.db.WithContext(ctx).Create(&data).Error
}

func (r *categoryRepository) FindOneByFilter(
	ctx context.Context,
	tx *gorm.DB,
	filter *repository.FindCategoryByFilter,
) (*entity.Category, error) {
	var category entity.Category
	err := r.buildFilter(ctx, tx, filter).First(&category).Error
	if err != nil {
		return nil, err
	}
	return &category, nil
}

func (r *categoryRepository) FindManyByFilter(
	ctx context.Context,
	tx *gorm.DB,
	filter *repository.FindCategoryByFilter,
) ([]entity.Category, error) {
	var category []entity.Category
	err := r.buildFilter(ctx, tx, filter).Find(&category).Error
	return category, err
}

func (r *categoryRepository) GetCategorySummary(
	ctx context.Context,
	tx *gorm.DB,
) ([]entity.Category, error) {
	query := r.db.WithContext(ctx)
	if tx != nil {
		query = tx.WithContext(ctx)
	}
	query = query.Model(&entity.Category{}).
		Select("categories.id, categories.name, COUNT(products.id) as product_count").
		Joins("LEFT JOIN products ON categories.id = products.category_id").
		Group("categories.id")
	rows, err := query.Rows()
	if err != nil {
		return nil, err
	}

	categories := make([]entity.Category, 0)
	for rows.Next() {
		var category entity.Category
		err := rows.Scan(&category.ID, &category.Name, &category.ProductCount)
		if err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}
	return categories, nil
}

// -------------------------------------------------------------------------------
func (r *categoryRepository) buildFilter(
	ctx context.Context,
	tx *gorm.DB,
	filter *repository.FindCategoryByFilter,
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

	if filter.Page != nil && filter.Limit != nil {
		offset := (*filter.Page - 1) * *filter.Limit
		query = query.Offset(offset).Limit(*filter.Limit)
	}

	if filter.Name != nil {
		query = query.Where("name = ?", *filter.Name)
	}

	return query
}
