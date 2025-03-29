package service

import (
	"context"

	"sondth-test_soa/app/entity"
	"sondth-test_soa/app/helper"
	"sondth-test_soa/app/model"
	"sondth-test_soa/app/repository"
	"sondth-test_soa/package/errors"
	"sondth-test_soa/utils"

	"golang.org/x/sync/errgroup"
)

type reviewService struct {
	postgresRepo repository.RepositoryCollections
	helper       helper.HelperCollections
}

func NewReviewService(
	postgresRepo repository.RepositoryCollections,
	helper helper.HelperCollections,
) IReviewService {
	return &reviewService{
		postgresRepo: postgresRepo,
		helper:       helper,
	}
}

func (s *reviewService) Create(
	ctx context.Context,
	req *model.CreateReviewRequest,
) (*model.CreateReviewResponse, error) {
	// Get user from context
	user, ok := ctx.Value(string(utils.USER_CONTEXT_KEY)).(*entity.User)
	if !ok {
		return nil, errors.New(errors.ErrCodeInternalServerError)
	}

	// Check if product exists
	product, err := s.postgresRepo.ProductRepo.FindOneByFilter(ctx, nil, &repository.FindProductByFilter{
		ID: &req.ProductID,
	})
	if err != nil {
		return nil, errors.New(errors.ErrCodeProductNotFound)
	}

	// Check if user already reviewed this product
	_, err = s.postgresRepo.ReviewRepo.FindOneByFilter(ctx, nil, &repository.FindReviewByFilter{
		UserID:    &user.ID,
		ProductID: &product.ID,
	})
	if err == nil {
		return nil, errors.New(errors.ErrCodeReviewAlreadyExists)
	}

	// Create review
	review := entity.NewReview()
	review.UserID = user.ID
	review.ProductID = product.ID
	review.Rating = req.Rating
	review.Comment = req.Comment

	if err := s.postgresRepo.ReviewRepo.Create(ctx, nil, review); err != nil {
		return nil, errors.New(errors.ErrCodeInternalServerError)
	}

	return &model.CreateReviewResponse{}, nil
}

func (s *reviewService) Delete(
	ctx context.Context,
	req *model.DeleteReviewRequest,
) (*model.DeleteReviewResponse, error) {
	// Get user from context
	user, ok := ctx.Value(string(utils.USER_CONTEXT_KEY)).(*entity.User)
	if !ok {
		return nil, errors.New(errors.ErrCodeInternalServerError)
	}

	// Find review
	review, err := s.postgresRepo.ReviewRepo.FindOneByFilter(ctx, nil, &repository.FindReviewByFilter{
		ID:     &req.ReviewID,
		UserID: &user.ID,
	})
	if err != nil {
		return nil, errors.New(errors.ErrCodeReviewNotFound)
	}

	// Delete review
	if err := s.postgresRepo.ReviewRepo.Delete(ctx, nil, review); err != nil {
		return nil, errors.New(errors.ErrCodeInternalServerError)
	}

	return &model.DeleteReviewResponse{}, nil
}

func (s *reviewService) GetReviews(
	ctx context.Context,
	req *model.GetReviewsRequest,
) (*model.GetReviewsResponse, error) {
	filter := &repository.FindReviewByFilter{
		ProductName: req.ProductName,
		Page:        req.Page,
		Limit:       req.Limit,
	}

	errGroup, errCtx := errgroup.WithContext(ctx)

	var reviews []entity.Review
	errGroup.Go(func() error {
		var err error
		reviews, err = s.postgresRepo.ReviewRepo.FindManyByFilter(errCtx, nil, filter)
		if err != nil {
			return err
		}
		return nil
	})

	var count int64
	errGroup.Go(func() error {
		var err error
		count, err = s.postgresRepo.ReviewRepo.CountByFilter(errCtx, nil, filter)
		if err != nil {
			return err
		}
		return nil
	})

	if err := errGroup.Wait(); err != nil {
		return nil, errors.New(errors.ErrCodeInternalServerError)
	}

	return &model.GetReviewsResponse{
		Reviews: reviews,
		Count:   count,
	}, nil
}

func (s *reviewService) GetReviewsSummary(
	ctx context.Context,
	req *model.GetReviewsSummaryRequest,
) (*model.GetReviewsSummaryResponse, error) {
	count, err := s.postgresRepo.ReviewRepo.CountByFilter(ctx, nil, &repository.FindReviewByFilter{})
	if err != nil {
		return nil, errors.New(errors.ErrCodeInternalServerError)
	}

	return &model.GetReviewsSummaryResponse{
		Count: count,
	}, nil
}
