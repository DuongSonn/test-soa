package service

import (
	"context"

	"sondth-test_soa/app/entity"
	"sondth-test_soa/app/helper"
	"sondth-test_soa/app/model"
	"sondth-test_soa/app/repository"
	"sondth-test_soa/package/errors"
	logger "sondth-test_soa/package/log"

	"golang.org/x/sync/errgroup"
)

type productService struct {
	postgresRepo repository.RepositoryCollections
	helper       helper.HelperCollections
}

func NewProductService(
	postgresRepo repository.RepositoryCollections,
	helper helper.HelperCollections,
) IProductService {
	return &productService{
		postgresRepo: postgresRepo,
		helper:       helper,
	}
}

func (s *productService) Create(
	ctx context.Context,
	req *model.CreateProductRequest,
) (*model.CreateProductResponse, error) {
	// Check if category exists
	if _, err := s.helper.CategoryHelper.ValidateCategoryID(ctx, req.CategoryID); err != nil {
		return nil, err
	}

	// Check if product name already exists
	if _, err := s.postgresRepo.ProductRepo.FindOneByFilter(ctx, nil, &repository.FindProductByFilter{
		Name: &req.Name,
	}); err == nil {
		return nil, errors.New(errors.ErrCodeProductExisted)
	}

	product := entity.NewProduct()
	product.Name = req.Name
	product.Description = req.Description
	product.Price = req.Price
	product.Quantity = req.Quantity
	product.CategoryID = req.CategoryID

	if err := s.postgresRepo.ProductRepo.Create(ctx, nil, product); err != nil {
		return nil, err
	}

	return &model.CreateProductResponse{}, nil
}

func (s *productService) Update(
	ctx context.Context,
	req *model.UpdateProductRequest,
) (*model.UpdateProductResponse, error) {
	// Check if product exists
	product, err := s.helper.ProductHelper.ValidateProductID(ctx, req.ID)
	if err != nil {
		return nil, err
	}

	// Check if category exists if it's being updated
	if req.CategoryID != product.CategoryID {
		if _, err := s.helper.CategoryHelper.ValidateCategoryID(ctx, req.CategoryID); err != nil {
			return nil, err
		}
		product.CategoryID = req.CategoryID
	}

	product.Name = req.Name
	product.Description = req.Description
	product.Price = req.Price
	product.Quantity = req.Quantity

	if err := s.postgresRepo.ProductRepo.Update(ctx, nil, product); err != nil {
		return nil, err
	}

	return &model.UpdateProductResponse{}, nil
}

func (s *productService) Delete(
	ctx context.Context,
	req *model.DeleteProductRequest,
) (*model.DeleteProductResponse, error) {
	// Check if product exists
	product, err := s.helper.ProductHelper.ValidateProductID(ctx, req.ID)
	if err != nil {
		return nil, err
	}

	if err := s.postgresRepo.ProductRepo.Delete(ctx, nil, product); err != nil {
		return nil, err
	}

	return &model.DeleteProductResponse{}, nil
}

func (s *productService) GetProducts(
	ctx context.Context,
	req *model.GetProductRequest,
) (*model.GetProductResponse, error) {
	results := &model.GetProductResponse{}
	filter := &repository.FindProductByFilter{
		Filter: repository.Filter{
			Fields: []string{"products.id", "products.name", "products.description", "products.price", "products.quantity", "products.category_id"},
		},
		Name:           req.Name,
		CategoryIDs:    req.CategoryIDs,
		Page:           req.Page,
		Limit:          req.Limit,
		Status:         req.Status,
		Order:          req.Order,
		CategoryFields: []string{"categories.id", "categories.name"},
	}

	errGroup, errCtx := errgroup.WithContext(ctx)
	errGroup.Go(func() error {
		count, err := s.postgresRepo.ProductRepo.CountByFilter(errCtx, nil, filter)
		if err != nil {
			return err
		}

		results.Count = count
		return nil
	})
	errGroup.Go(func() error {
		products, err := s.postgresRepo.ProductRepo.FindManyByFilter(errCtx, nil, filter)
		if err != nil {
			return err
		}
		// Set status for each product based on quantity
		for i := range products {
			if products[i].Quantity > 0 {
				products[i].Status = entity.PRODUCT_STATUS_IN_STOCK
			} else {
				products[i].Status = entity.PRODUCT_STATUS_OUT_OF_STOCK
			}
		}

		results.Result = products
		return nil
	})
	if err := errGroup.Wait(); err != nil {
		logger.WithCtx(ctx).Error("GetProducts", err)
		return nil, err
	}

	return results, nil
}
