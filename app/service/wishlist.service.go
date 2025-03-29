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
	"gorm.io/gorm"
)

type wishlistService struct {
	postgresRepo repository.RepositoryCollections
	helper       helper.HelperCollections
}

func NewWishlistService(
	postgresRepo repository.RepositoryCollections,
	helper helper.HelperCollections,
) IWishlistService {
	return &wishlistService{
		postgresRepo: postgresRepo,
		helper:       helper,
	}
}

func (s *wishlistService) AddToWishlist(
	ctx context.Context,
	req *model.AddToWishlistRequest,
) (*model.AddToWishlistResponse, error) {
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

	// Check if already in wishlist
	existedWishList, err := s.postgresRepo.WishlistRepo.FindOneByFilter(ctx, nil, &repository.FindWishlistByFilter{
		UserID:    &user.ID,
		ProductID: &product.ID,
	})
	if err != nil && err != gorm.ErrRecordNotFound || existedWishList != nil {
		return nil, errors.New(errors.ErrCodeProductAlreadyInWishlist)
	}

	// Create wishlist
	wishlist := entity.NewWishlist()
	wishlist.UserID = user.ID
	wishlist.ProductID = product.ID
	if err := s.postgresRepo.WishlistRepo.Create(ctx, nil, wishlist); err != nil {
		return nil, errors.New(errors.ErrCodeInternalServerError)
	}

	return &model.AddToWishlistResponse{}, nil
}

func (s *wishlistService) DeleteFromWishlist(
	ctx context.Context,
	req *model.DeleteFromWishlistRequest,
) (*model.DeleteFromWishlistResponse, error) {
	// Get user from context
	user, ok := ctx.Value(string(utils.USER_CONTEXT_KEY)).(*entity.User)
	if !ok {
		return nil, errors.New(errors.ErrCodeInternalServerError)
	}

	// Find wishlist
	wishlist, err := s.postgresRepo.WishlistRepo.FindOneByFilter(ctx, nil, &repository.FindWishlistByFilter{
		ID:     &req.ID,
		UserID: &user.ID,
	})
	if err != nil {
		return nil, errors.New(errors.ErrCodeWishlistNotFound)
	}

	// Delete wishlist
	if err := s.postgresRepo.WishlistRepo.Delete(ctx, nil, wishlist); err != nil {
		return nil, errors.New(errors.ErrCodeInternalServerError)
	}

	return &model.DeleteFromWishlistResponse{}, nil
}

func (s *wishlistService) GetWishlists(
	ctx context.Context,
	req *model.GetWishlistsRequest,
) (*model.GetWishlistsResponse, error) {
	// Get user from context
	user, ok := ctx.Value(string(utils.USER_CONTEXT_KEY)).(*entity.User)
	if !ok {
		return nil, errors.New(errors.ErrCodeInternalServerError)
	}

	filter := &repository.FindWishlistByFilter{
		UserID: &user.ID,
		Page:   &req.Page,
		Limit:  &req.Limit,
	}

	errGroup, errCtx := errgroup.WithContext(ctx)

	var wishlists []entity.Wishlist
	errGroup.Go(func() error {
		var err error
		wishlists, err = s.postgresRepo.WishlistRepo.FindManyByFilter(errCtx, nil, filter)
		if err != nil {
			return err
		}
		return nil
	})

	var count int64
	errGroup.Go(func() error {
		var err error
		count, err = s.postgresRepo.WishlistRepo.CountByFilter(errCtx, nil, filter)
		if err != nil {
			return err
		}
		return nil
	})

	if err := errGroup.Wait(); err != nil {
		return nil, errors.New(errors.ErrCodeInternalServerError)
	}

	return &model.GetWishlistsResponse{
		Wishlists: wishlists,
		Count:     count,
	}, nil
}

func (s *wishlistService) GetWishlistsSummary(
	ctx context.Context,
	req *model.GetWishlistsSummaryRequest,
) (*model.GetWishlistsSummaryResponse, error) {
	count, err := s.postgresRepo.WishlistRepo.CountByFilter(ctx, nil, &repository.FindWishlistByFilter{})
	if err != nil {
		return nil, errors.New(errors.ErrCodeInternalServerError)
	}

	return &model.GetWishlistsSummaryResponse{
		Count: count,
	}, nil
}
