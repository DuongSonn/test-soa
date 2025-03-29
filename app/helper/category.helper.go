package helper

import (
	"context"
	"sondth-test_soa/app/entity"
	"sondth-test_soa/app/repository"
	"sondth-test_soa/package/errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type categoryHelper struct {
	postgresRepo repository.RepositoryCollections
}

func NewCategoryHelper(postgresRepo repository.RepositoryCollections) ICategoryHelper {
	return &categoryHelper{
		postgresRepo: postgresRepo,
	}
}

func (s *categoryHelper) ValidateCategoryID(ctx context.Context, categoryID uuid.UUID) (*entity.Category, error) {
	category, err := s.postgresRepo.CategoryRepo.FindOneByFilter(ctx, nil, &repository.FindCategoryByFilter{
		Filter: repository.Filter{
			Fields: []string{"id"},
		},
		ID: &categoryID,
	})
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New(errors.ErrCodeCategoryNotFound)
		}
		return nil, err
	}

	return category, nil
}
