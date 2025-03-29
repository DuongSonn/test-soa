package service

import (
	"context"
	"testing"

	"sondth-test_soa/app/entity"
	"sondth-test_soa/app/helper"
	"sondth-test_soa/app/model"
	"sondth-test_soa/app/repository"
	helper_mocks "sondth-test_soa/mocks/helper"
	repo_mocks "sondth-test_soa/mocks/repository"
	"sondth-test_soa/package/errors"
	"sondth-test_soa/utils"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

var (
	wishlistID = uuid.New()
)

func Test_wishlistService_AddToWishlist(t *testing.T) {
	type args struct {
		ctx context.Context
		req *model.AddToWishlistRequest
	}

	type testCase struct {
		name    string
		args    args
		want    *model.AddToWishlistResponse
		wantErr bool
		mock    func(wishlistRepo *repo_mocks.IWishlistRepository, productRepo *repo_mocks.IProductRepository)
	}

	ctx := context.WithValue(context.Background(), string(utils.USER_CONTEXT_KEY), &entity.User{
		ID: userID,
	})

	tests := []testCase{
		{
			name: "Add to Wishlist Success",
			args: args{
				ctx: ctx,
				req: &model.AddToWishlistRequest{
					ProductID: productID,
				},
			},
			want:    &model.AddToWishlistResponse{},
			wantErr: false,
			mock: func(wishlistRepo *repo_mocks.IWishlistRepository, productRepo *repo_mocks.IProductRepository) {
				// Mock product exists check
				productRepo.On("FindOneByFilter", mock.MatchedBy(func(c context.Context) bool {
					return true
				}), mock.Anything, mock.MatchedBy(func(filter *repository.FindProductByFilter) bool {
					return filter.ID != nil && *filter.ID == productID
				})).Return(&entity.Product{ID: productID}, nil).Once()

				// Mock wishlist not exists check
				wishlistRepo.On("FindOneByFilter", mock.MatchedBy(func(c context.Context) bool {
					return true
				}), mock.Anything, mock.MatchedBy(func(filter *repository.FindWishlistByFilter) bool {
					return filter.UserID != nil && *filter.UserID == userID &&
						filter.ProductID != nil && *filter.ProductID == productID
				})).Return(nil, gorm.ErrRecordNotFound).Once()

				// Mock wishlist creation
				wishlistRepo.On("Create", mock.MatchedBy(func(c context.Context) bool {
					return true
				}), mock.Anything, mock.MatchedBy(func(wishlist *entity.Wishlist) bool {
					return wishlist.UserID == userID && wishlist.ProductID == productID
				})).Return(nil).Once()
			},
		},
		{
			name: "Add to Wishlist Failed - No User Context",
			args: args{
				ctx: context.Background(),
				req: &model.AddToWishlistRequest{
					ProductID: productID,
				},
			},
			want:    nil,
			wantErr: true,
			mock:    func(wishlistRepo *repo_mocks.IWishlistRepository, productRepo *repo_mocks.IProductRepository) {},
		},
		{
			name: "Add to Wishlist Failed - Product Not Found",
			args: args{
				ctx: ctx,
				req: &model.AddToWishlistRequest{
					ProductID: productID,
				},
			},
			want:    nil,
			wantErr: true,
			mock: func(wishlistRepo *repo_mocks.IWishlistRepository, productRepo *repo_mocks.IProductRepository) {
				// Mock product not found
				productRepo.On("FindOneByFilter", mock.MatchedBy(func(c context.Context) bool {
					return true
				}), mock.Anything, mock.MatchedBy(func(filter *repository.FindProductByFilter) bool {
					return filter.ID != nil && *filter.ID == productID
				})).Return(nil, gorm.ErrRecordNotFound).Once()
			},
		},
		{
			name: "Add to Wishlist Failed - Already in Wishlist",
			args: args{
				ctx: ctx,
				req: &model.AddToWishlistRequest{
					ProductID: productID,
				},
			},
			want:    nil,
			wantErr: true,
			mock: func(wishlistRepo *repo_mocks.IWishlistRepository, productRepo *repo_mocks.IProductRepository) {
				// Mock product exists check
				productRepo.On("FindOneByFilter", mock.MatchedBy(func(c context.Context) bool {
					return true
				}), mock.Anything, mock.MatchedBy(func(filter *repository.FindProductByFilter) bool {
					return filter.ID != nil && *filter.ID == productID
				})).Return(&entity.Product{ID: productID}, nil).Once()

				// Mock wishlist already exists
				wishlistRepo.On("FindOneByFilter", mock.MatchedBy(func(c context.Context) bool {
					return true
				}), mock.Anything, mock.MatchedBy(func(filter *repository.FindWishlistByFilter) bool {
					return filter.UserID != nil && *filter.UserID == userID &&
						filter.ProductID != nil && *filter.ProductID == productID
				})).Return(&entity.Wishlist{}, nil).Once()
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Initialize mocks
			wishlistRepo := repo_mocks.NewIWishlistRepository(t)
			productRepo := repo_mocks.NewIProductRepository(t)

			// Setup mocks
			tt.mock(wishlistRepo, productRepo)

			s := &wishlistService{
				postgresRepo: repository.RepositoryCollections{
					WishlistRepo: wishlistRepo,
					ProductRepo:  productRepo,
				},
			}

			got, err := s.AddToWishlist(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("wishlistService.AddToWishlist() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil && got == nil {
				t.Error("wishlistService.AddToWishlist() got nil response, want non-nil")
			}
		})
	}
}

func Test_wishlistService_DeleteFromWishlist(t *testing.T) {
	type args struct {
		ctx context.Context
		req *model.DeleteFromWishlistRequest
	}

	type testCase struct {
		name    string
		args    args
		want    *model.DeleteFromWishlistResponse
		wantErr bool
		mock    func(repo *repo_mocks.IWishlistRepository)
	}

	ctx := context.WithValue(context.Background(), string(utils.USER_CONTEXT_KEY), &entity.User{
		ID: userID,
	})

	tests := []testCase{
		{
			name: "Delete from Wishlist Success",
			args: args{
				ctx: ctx,
				req: &model.DeleteFromWishlistRequest{
					ID: wishlistID,
				},
			},
			want:    &model.DeleteFromWishlistResponse{},
			wantErr: false,
			mock: func(repo *repo_mocks.IWishlistRepository) {
				// Mock wishlist exists check
				repo.On("FindOneByFilter", mock.MatchedBy(func(c context.Context) bool {
					return true
				}), mock.Anything, mock.MatchedBy(func(filter *repository.FindWishlistByFilter) bool {
					return filter.ID != nil && *filter.ID == wishlistID &&
						filter.UserID != nil && *filter.UserID == userID
				})).Return(&entity.Wishlist{ID: wishlistID}, nil).Once()

				// Mock wishlist deletion
				repo.On("Delete", mock.MatchedBy(func(c context.Context) bool {
					return true
				}), mock.Anything, mock.MatchedBy(func(wishlist *entity.Wishlist) bool {
					return wishlist.ID == wishlistID
				})).Return(nil).Once()
			},
		},
		{
			name: "Delete from Wishlist Failed - No User Context",
			args: args{
				ctx: context.Background(),
				req: &model.DeleteFromWishlistRequest{
					ID: wishlistID,
				},
			},
			want:    nil,
			wantErr: true,
			mock:    func(repo *repo_mocks.IWishlistRepository) {},
		},
		{
			name: "Delete from Wishlist Failed - Not Found",
			args: args{
				ctx: ctx,
				req: &model.DeleteFromWishlistRequest{
					ID: wishlistID,
				},
			},
			want:    nil,
			wantErr: true,
			mock: func(repo *repo_mocks.IWishlistRepository) {
				// Mock wishlist not found
				repo.On("FindOneByFilter", mock.MatchedBy(func(c context.Context) bool {
					return true
				}), mock.Anything, mock.MatchedBy(func(filter *repository.FindWishlistByFilter) bool {
					return filter.ID != nil && *filter.ID == wishlistID &&
						filter.UserID != nil && *filter.UserID == userID
				})).Return(nil, gorm.ErrRecordNotFound).Once()
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Initialize mocks
			repo := repo_mocks.NewIWishlistRepository(t)

			// Setup mocks
			tt.mock(repo)

			s := &wishlistService{
				postgresRepo: repository.RepositoryCollections{
					WishlistRepo: repo,
				},
			}

			got, err := s.DeleteFromWishlist(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("wishlistService.DeleteFromWishlist() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil && got == nil {
				t.Error("wishlistService.DeleteFromWishlist() got nil response, want non-nil")
			}
		})
	}
}

func Test_wishlistService_GetWishlists(t *testing.T) {
	type args struct {
		ctx context.Context
		req *model.GetWishlistsRequest
	}

	type testCase struct {
		name    string
		args    args
		want    *model.GetWishlistsResponse
		wantErr bool
		mock    func(repo *repo_mocks.IWishlistRepository)
	}

	ctx := context.WithValue(context.Background(), string(utils.USER_CONTEXT_KEY), &entity.User{
		ID: userID,
	})

	wishlists := []entity.Wishlist{
		{
			ID:        wishlistID,
			UserID:    userID,
			ProductID: productID,
		},
	}

	tests := []testCase{
		{
			name: "Get Wishlists Success",
			args: args{
				ctx: ctx,
				req: &model.GetWishlistsRequest{
					Page:  page,
					Limit: limit,
				},
			},
			want: &model.GetWishlistsResponse{
				Wishlists: wishlists,
				Count:     1,
			},
			wantErr: false,
			mock: func(repo *repo_mocks.IWishlistRepository) {
				// Mock find wishlists
				repo.On("FindManyByFilter", mock.MatchedBy(func(c context.Context) bool {
					return true
				}), mock.Anything, mock.MatchedBy(func(filter *repository.FindWishlistByFilter) bool {
					return filter.UserID != nil && *filter.UserID == userID &&
						filter.Page != nil && *filter.Page == page &&
						filter.Limit != nil && *filter.Limit == limit
				})).Return(wishlists, nil).Once()

				// Mock count wishlists
				repo.On("CountByFilter", mock.MatchedBy(func(c context.Context) bool {
					return true
				}), mock.Anything, mock.MatchedBy(func(filter *repository.FindWishlistByFilter) bool {
					return filter.UserID != nil && *filter.UserID == userID &&
						filter.Page != nil && *filter.Page == page &&
						filter.Limit != nil && *filter.Limit == limit
				})).Return(int64(1), nil).Once()
			},
		},
		{
			name: "Get Wishlists Failed - No User Context",
			args: args{
				ctx: context.Background(),
				req: &model.GetWishlistsRequest{
					Page:  page,
					Limit: limit,
				},
			},
			want:    nil,
			wantErr: true,
			mock:    func(repo *repo_mocks.IWishlistRepository) {},
		},
		{
			name: "Get Wishlists Failed - Find Error",
			args: args{
				ctx: ctx,
				req: &model.GetWishlistsRequest{
					Page:  page,
					Limit: limit,
				},
			},
			want:    nil,
			wantErr: true,
			mock: func(repo *repo_mocks.IWishlistRepository) {
				// Mock find wishlists error
				repo.On("FindManyByFilter", mock.MatchedBy(func(c context.Context) bool {
					return true
				}), mock.Anything, mock.MatchedBy(func(filter *repository.FindWishlistByFilter) bool {
					return filter.UserID != nil && *filter.UserID == userID &&
						filter.Page != nil && *filter.Page == page &&
						filter.Limit != nil && *filter.Limit == limit
				})).Return(nil, errors.New(errors.ErrCodeInternalServerError)).Once()

				// Mock count wishlists error
				repo.On("CountByFilter", mock.MatchedBy(func(c context.Context) bool {
					return true
				}), mock.Anything, mock.MatchedBy(func(filter *repository.FindWishlistByFilter) bool {
					return filter.UserID != nil && *filter.UserID == userID &&
						filter.Page != nil && *filter.Page == page &&
						filter.Limit != nil && *filter.Limit == limit
				})).Return(int64(0), errors.New(errors.ErrCodeInternalServerError)).Once()
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Initialize mocks
			repo := repo_mocks.NewIWishlistRepository(t)

			// Setup mocks
			tt.mock(repo)

			s := &wishlistService{
				postgresRepo: repository.RepositoryCollections{
					WishlistRepo: repo,
				},
			}

			got, err := s.GetWishlists(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("wishlistService.GetWishlists() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil && got == nil {
				t.Error("wishlistService.GetWishlists() got nil response, want non-nil")
			}
		})
	}
}

func TestNewWishlistService(t *testing.T) {
	type args struct {
		postgresRepo repository.RepositoryCollections
		helper       helper.HelperCollections
	}
	tests := []struct {
		name string
		args args
		want IWishlistService
	}{
		{
			name: "NewWishlistService",
			args: args{
				postgresRepo: repository.RepositoryCollections{
					WishlistRepo: repo_mocks.NewIWishlistRepository(t),
					ProductRepo:  repo_mocks.NewIProductRepository(t),
				},
				helper: helper.HelperCollections{
					UserHelper: helper_mocks.NewIUserHelper(t),
				},
			},
			want: &wishlistService{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewWishlistService(tt.args.postgresRepo, tt.args.helper)
			if got == nil {
				t.Error("NewWishlistService() got nil service")
			}
		})
	}
}
