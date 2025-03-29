package helper

import (
	"context"
	"sondth-test_soa/app/entity"
	"sondth-test_soa/app/repository"
	"sondth-test_soa/package/errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type productHelper struct {
	postgresRepo repository.RepositoryCollections
}

func NewProductHelper(postgresRepo repository.RepositoryCollections) IProductHelper {
	return &productHelper{}
}

func (s *productHelper) ValidateProductID(ctx context.Context, productID uuid.UUID) (*entity.Product, error) {
	product, err := s.postgresRepo.ProductRepo.FindOneByFilter(ctx, nil, &repository.FindProductByFilter{
		Filter: repository.Filter{
			Fields: []string{"id"},
		},
		ID: &productID,
	})
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New(errors.ErrCodeProductNotFound)
		}
		return nil, err
	}

	return product, nil
}
