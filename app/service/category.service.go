package service

import (
	"context"

	"gorm.io/gorm"

	"sondth-test_soa/app/entity"
	"sondth-test_soa/app/helper"
	"sondth-test_soa/app/model"
	"sondth-test_soa/app/repository"
	"sondth-test_soa/package/errors"
)

type categoryService struct {
	postgresRepo repository.RepositoryCollections
	helper       helper.HelperCollections
}

func NewCategoryService(
	postgresRepo repository.RepositoryCollections,
	helper helper.HelperCollections,
) ICategoryService {
	return &categoryService{
		postgresRepo: postgresRepo,
		helper:       helper,
	}
}

func (s *categoryService) Create(
	ctx context.Context,
	req *model.CreateCategoryRequest,
) (*model.CreateCategoryResponse, error) {
	existedCategory, err := s.postgresRepo.CategoryRepo.FindOneByFilter(ctx, nil, &repository.FindCategoryByFilter{
		Name: &req.Name,
	})
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	if existedCategory != nil {
		return nil, errors.New(errors.ErrCodeCategoryExisted)
	}

	category := entity.NewCategory()
	category.Name = req.Name
	if err := s.postgresRepo.CategoryRepo.Create(ctx, nil, category); err != nil {
		return nil, err
	}

	return &model.CreateCategoryResponse{}, nil
}

func (s *categoryService) GetCategories(
	ctx context.Context,
	req *model.GetCategoriesRequest,
) (*model.GetCategoriesResponse, error) {
	var (
		defaultPage  = 1
		defaultLimit = 10
		filter       = &repository.FindCategoryByFilter{
			Filter: repository.Filter{
				Fields: []string{"id", "name", "description"},
			},
		}
	)

	if req.Page != nil && req.Limit != nil {
		filter.Page = req.Page
		filter.Limit = req.Limit
	} else {
		filter.Page = &defaultPage
		filter.Limit = &defaultLimit
	}

	categories, err := s.postgresRepo.CategoryRepo.FindManyByFilter(ctx, nil, filter)
	if err != nil {
		return nil, err
	}

	return &model.GetCategoriesResponse{
		Count:  int64(len(categories)),
		Result: categories,
	}, nil
}

func (s *categoryService) GetCategoriesSummary(
	ctx context.Context,
	req *model.GetCategoriesSummaryRequest,
) (*model.GetCategoriesSummaryResponse, error) {
	// Get all categories without pagination
	categories, err := s.postgresRepo.CategoryRepo.GetCategorySummary(ctx, nil)
	if err != nil {
		return nil, err
	}

	return &model.GetCategoriesSummaryResponse{
		Categories: categories,
	}, nil
}
