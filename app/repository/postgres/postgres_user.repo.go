package postgres

import (
	"context"

	"gorm.io/gorm"

	"sondth-test_soa/app/entity"
	"sondth-test_soa/app/repository"
)

type userRepository struct {
	db *gorm.DB
}

func NewPostgresUserRepository(db *gorm.DB) repository.IUserRepository {
	return &userRepository{
		db,
	}
}

func (r *userRepository) Create(
	ctx context.Context,
	tx *gorm.DB,
	data *entity.User,
) error {
	if tx != nil {
		return tx.WithContext(ctx).Create(&data).Error
	}

	return r.db.WithContext(ctx).Create(&data).Error
}

func (r *userRepository) Update(
	ctx context.Context,
	tx *gorm.DB,
	data *entity.User,
) error {
	if tx != nil {
		return tx.WithContext(ctx).Save(&data).Error
	}

	return r.db.WithContext(ctx).Save(&data).Error
}

func (r *userRepository) FindOneByFilter(
	ctx context.Context,
	tx *gorm.DB,
	filter *repository.FindUserByFilter,
) (*entity.User, error) {
	var user entity.User
	err := r.buildFilter(ctx, tx, filter).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindManyByFilter(
	ctx context.Context,
	tx *gorm.DB,
	filter *repository.FindUserByFilter,
) ([]entity.User, error) {
	var users []entity.User

	query := r.buildFilter(ctx, tx, filter)
	if filter.Page != nil && filter.Limit != nil {
		offset := (*filter.Page - 1) * *filter.Limit
		query = query.Offset(offset).Limit(*filter.Limit)
	}

	err := query.Find(&users).Error
	return users, err
}

func (r *userRepository) CountByFilter(
	ctx context.Context,
	tx *gorm.DB,
	filter *repository.FindUserByFilter,
) (int64, error) {
	var count int64
	err := r.buildFilter(ctx, tx, filter).Count(&count).Error
	return count, err
}

// -------------------------------------------------------------------------------
func (r *userRepository) buildFilter(
	ctx context.Context,
	tx *gorm.DB,
	filter *repository.FindUserByFilter,
) *gorm.DB {
	query := r.db.WithContext(ctx)
	if tx != nil {
		query = tx.WithContext(ctx)
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
		query = query.Where("username ILIKE ? OR email ILIKE ?", "%"+*filter.Name+"%", "%"+*filter.Name+"%")
	}

	if filter.Role != nil {
		query = query.Where("role = ?", filter.Role)
	}

	if filter.Username != nil {
		query = query.Where("username = ?", *filter.Username)
	}

	return query
}
