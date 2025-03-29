package service

import (
	"context"
	"testing"

	"sondth-test_soa/app/entity"
	"sondth-test_soa/app/helper"
	"sondth-test_soa/app/model"
	"sondth-test_soa/app/repository"
	repo_mocks "sondth-test_soa/mocks/repository"
	"sondth-test_soa/package/errors"
	"sondth-test_soa/utils"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

var (
	reviewID    = uuid.New()
	userID      = uuid.New()
	rating      = float64(5)
	comment     = "Great product!"
	productID   = uuid.New()
	productName = "Test Product"
	page        = 1
	limit       = 10
)

func Test_reviewService_Create(t *testing.T) {
	type args struct {
		ctx context.Context
		req *model.CreateReviewRequest
	}

	type testCase struct {
		name    string
		args    args
		want    *model.CreateReviewResponse
		wantErr bool
		mock    func(repo *repo_mocks.IReviewRepository, productRepo *repo_mocks.IProductRepository)
	}

	ctx := context.WithValue(context.Background(), string(utils.USER_CONTEXT_KEY), &entity.User{
		ID: userID,
	})

	tests := []testCase{
		{
			name: "Create Success",
			args: args{
				ctx: ctx,
				req: &model.CreateReviewRequest{
					ProductID: productID,
					Rating:    rating,
					Comment:   comment,
				},
			},
			want:    &model.CreateReviewResponse{},
			wantErr: false,
			mock: func(repo *repo_mocks.IReviewRepository, productRepo *repo_mocks.IProductRepository) {
				// Mock product exists check
				productRepo.On("FindOneByFilter", ctx, mock.Anything, mock.MatchedBy(func(filter *repository.FindProductByFilter) bool {
					return filter.ID != nil && *filter.ID == productID
				})).Return(&entity.Product{ID: productID}, nil).Once()

				// Mock review not exists check
				repo.On("FindOneByFilter", ctx, mock.Anything, mock.MatchedBy(func(filter *repository.FindReviewByFilter) bool {
					return filter.UserID != nil && *filter.UserID == userID &&
						filter.ProductID != nil && *filter.ProductID == productID
				})).Return(nil, gorm.ErrRecordNotFound).Once()

				// Mock review creation
				repo.On("Create", ctx, mock.Anything, mock.MatchedBy(func(review *entity.Review) bool {
					return review.UserID == userID &&
						review.ProductID == productID &&
						review.Rating == float64(rating) &&
						review.Comment == comment
				})).Return(nil).Once()
			},
		},
		{
			name: "Create Failed - No User Context",
			args: args{
				ctx: context.Background(),
				req: &model.CreateReviewRequest{
					ProductID: productID,
					Rating:    rating,
					Comment:   comment,
				},
			},
			want:    nil,
			wantErr: true,
			mock:    func(repo *repo_mocks.IReviewRepository, productRepo *repo_mocks.IProductRepository) {},
		},
		{
			name: "Create Failed - Product Not Found",
			args: args{
				ctx: ctx,
				req: &model.CreateReviewRequest{
					ProductID: productID,
					Rating:    rating,
					Comment:   comment,
				},
			},
			want:    nil,
			wantErr: true,
			mock: func(repo *repo_mocks.IReviewRepository, productRepo *repo_mocks.IProductRepository) {
				// Mock product not found
				productRepo.On("FindOneByFilter", ctx, mock.Anything, mock.MatchedBy(func(filter *repository.FindProductByFilter) bool {
					return filter.ID != nil && *filter.ID == productID
				})).Return(nil, gorm.ErrRecordNotFound).Once()
			},
		},
		{
			name: "Create Failed - Review Already Exists",
			args: args{
				ctx: ctx,
				req: &model.CreateReviewRequest{
					ProductID: productID,
					Rating:    rating,
					Comment:   comment,
				},
			},
			want:    nil,
			wantErr: true,
			mock: func(repo *repo_mocks.IReviewRepository, productRepo *repo_mocks.IProductRepository) {
				// Mock product exists check
				productRepo.On("FindOneByFilter", ctx, mock.Anything, mock.MatchedBy(func(filter *repository.FindProductByFilter) bool {
					return filter.ID != nil && *filter.ID == productID
				})).Return(&entity.Product{ID: productID}, nil).Once()

				// Mock review already exists
				repo.On("FindOneByFilter", ctx, mock.Anything, mock.MatchedBy(func(filter *repository.FindReviewByFilter) bool {
					return filter.UserID != nil && *filter.UserID == userID &&
						filter.ProductID != nil && *filter.ProductID == productID
				})).Return(&entity.Review{}, nil).Once()
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Initialize mocks
			repo := repo_mocks.NewIReviewRepository(t)
			productRepo := repo_mocks.NewIProductRepository(t)

			// Setup mocks
			tt.mock(repo, productRepo)

			s := &reviewService{
				postgresRepo: repository.RepositoryCollections{
					ReviewRepo:  repo,
					ProductRepo: productRepo,
				},
			}

			got, err := s.Create(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("reviewService.Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil && got == nil {
				t.Error("reviewService.Create() got nil response, want non-nil")
			}
		})
	}
}

func Test_reviewService_Delete(t *testing.T) {
	type args struct {
		ctx context.Context
		req *model.DeleteReviewRequest
	}

	type testCase struct {
		name    string
		args    args
		want    *model.DeleteReviewResponse
		wantErr bool
		mock    func(repo *repo_mocks.IReviewRepository)
	}

	ctx := context.WithValue(context.Background(), string(utils.USER_CONTEXT_KEY), &entity.User{
		ID: userID,
	})

	tests := []testCase{
		{
			name: "Delete Success",
			args: args{
				ctx: ctx,
				req: &model.DeleteReviewRequest{
					ReviewID: reviewID,
				},
			},
			want:    &model.DeleteReviewResponse{},
			wantErr: false,
			mock: func(repo *repo_mocks.IReviewRepository) {
				// Mock review exists check
				repo.On("FindOneByFilter", ctx, mock.Anything, mock.MatchedBy(func(filter *repository.FindReviewByFilter) bool {
					return filter.ID != nil && *filter.ID == reviewID &&
						filter.UserID != nil && *filter.UserID == userID
				})).Return(&entity.Review{ID: reviewID}, nil).Once()

				// Mock review deletion
				repo.On("Delete", ctx, mock.Anything, mock.MatchedBy(func(review *entity.Review) bool {
					return review.ID == reviewID
				})).Return(nil).Once()
			},
		},
		{
			name: "Delete Failed - No User Context",
			args: args{
				ctx: context.Background(),
				req: &model.DeleteReviewRequest{
					ReviewID: reviewID,
				},
			},
			want:    nil,
			wantErr: true,
			mock:    func(repo *repo_mocks.IReviewRepository) {},
		},
		{
			name: "Delete Failed - Review Not Found",
			args: args{
				ctx: ctx,
				req: &model.DeleteReviewRequest{
					ReviewID: reviewID,
				},
			},
			want:    nil,
			wantErr: true,
			mock: func(repo *repo_mocks.IReviewRepository) {
				// Mock review not found
				repo.On("FindOneByFilter", ctx, mock.Anything, mock.MatchedBy(func(filter *repository.FindReviewByFilter) bool {
					return filter.ID != nil && *filter.ID == reviewID &&
						filter.UserID != nil && *filter.UserID == userID
				})).Return(nil, gorm.ErrRecordNotFound).Once()
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Initialize mocks
			repo := repo_mocks.NewIReviewRepository(t)

			// Setup mocks
			tt.mock(repo)

			s := &reviewService{
				postgresRepo: repository.RepositoryCollections{
					ReviewRepo: repo,
				},
			}

			got, err := s.Delete(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("reviewService.Delete() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil && got == nil {
				t.Error("reviewService.Delete() got nil response, want non-nil")
			}
		})
	}
}

func Test_reviewService_GetReviews(t *testing.T) {
	type args struct {
		ctx context.Context
		req *model.GetReviewsRequest
	}

	type testCase struct {
		name    string
		args    args
		want    *model.GetReviewsResponse
		wantErr bool
		mock    func(repo *repo_mocks.IReviewRepository)
	}

	ctx := context.WithValue(context.Background(), string(utils.USER_CONTEXT_KEY), &entity.User{
		ID: userID,
	})
	reviews := []entity.Review{
		{
			ID:        reviewID,
			UserID:    userID,
			ProductID: productID,
			Rating:    rating,
			Comment:   comment,
		},
	}

	tests := []testCase{
		{
			name: "Get Reviews Success",
			args: args{
				ctx: ctx,
				req: &model.GetReviewsRequest{
					ProductName: &productName,
					Page:        &page,
					Limit:       &limit,
				},
			},
			want: &model.GetReviewsResponse{
				Reviews: reviews,
				Count:   1,
			},
			wantErr: false,
			mock: func(repo *repo_mocks.IReviewRepository) {
				// Mock find reviews
				repo.On("FindManyByFilter", mock.MatchedBy(func(c context.Context) bool {
					return true
				}), mock.Anything, mock.MatchedBy(func(filter *repository.FindReviewByFilter) bool {
					return filter.ProductName != nil && *filter.ProductName == productName &&
						filter.Page != nil && *filter.Page == page &&
						filter.Limit != nil && *filter.Limit == limit
				})).Return(reviews, nil).Once()

				// Mock count reviews
				repo.On("CountByFilter", mock.MatchedBy(func(c context.Context) bool {
					return true
				}), mock.Anything, mock.MatchedBy(func(filter *repository.FindReviewByFilter) bool {
					return filter.ProductName != nil && *filter.ProductName == productName &&
						filter.Page != nil && *filter.Page == page &&
						filter.Limit != nil && *filter.Limit == limit
				})).Return(int64(1), nil).Once()
			},
		},
		{
			name: "Get Reviews Failed - Find Error",
			args: args{
				ctx: ctx,
				req: &model.GetReviewsRequest{
					ProductName: &productName,
					Page:        &page,
					Limit:       &limit,
				},
			},
			want:    nil,
			wantErr: true,
			mock: func(repo *repo_mocks.IReviewRepository) {
				// Mock find reviews error
				repo.On("FindManyByFilter", mock.MatchedBy(func(c context.Context) bool {
					return true
				}), mock.Anything, mock.MatchedBy(func(filter *repository.FindReviewByFilter) bool {
					return filter.ProductName != nil && *filter.ProductName == productName &&
						filter.Page != nil && *filter.Page == page &&
						filter.Limit != nil && *filter.Limit == limit
				})).Return(nil, errors.New(errors.ErrCodeInternalServerError)).Once()

				// Mock count reviews
				repo.On("CountByFilter", mock.MatchedBy(func(c context.Context) bool {
					return true
				}), mock.Anything, mock.MatchedBy(func(filter *repository.FindReviewByFilter) bool {
					return filter.ProductName != nil && *filter.ProductName == productName &&
						filter.Page != nil && *filter.Page == page &&
						filter.Limit != nil && *filter.Limit == limit
				})).Return(int64(0), errors.New(errors.ErrCodeInternalServerError)).Once()
			},
		},
		{
			name: "Get Reviews Failed - Count Error",
			args: args{
				ctx: ctx,
				req: &model.GetReviewsRequest{
					ProductName: &productName,
					Page:        &page,
					Limit:       &limit,
				},
			},
			want:    nil,
			wantErr: true,
			mock: func(repo *repo_mocks.IReviewRepository) {
				// Mock find reviews
				repo.On("FindManyByFilter", mock.MatchedBy(func(c context.Context) bool {
					return true
				}), mock.Anything, mock.MatchedBy(func(filter *repository.FindReviewByFilter) bool {
					return filter.ProductName != nil && *filter.ProductName == productName &&
						filter.Page != nil && *filter.Page == page &&
						filter.Limit != nil && *filter.Limit == limit
				})).Return(reviews, nil).Once()

				// Mock count reviews error
				repo.On("CountByFilter", mock.MatchedBy(func(c context.Context) bool {
					return true
				}), mock.Anything, mock.MatchedBy(func(filter *repository.FindReviewByFilter) bool {
					return filter.ProductName != nil && *filter.ProductName == productName &&
						filter.Page != nil && *filter.Page == page &&
						filter.Limit != nil && *filter.Limit == limit
				})).Return(int64(0), errors.New(errors.ErrCodeInternalServerError)).Once()
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Initialize mocks
			repo := repo_mocks.NewIReviewRepository(t)

			// Setup mocks
			tt.mock(repo)

			s := &reviewService{
				postgresRepo: repository.RepositoryCollections{
					ReviewRepo: repo,
				},
			}

			got, err := s.GetReviews(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("reviewService.GetReviews() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil && got == nil {
				t.Error("reviewService.GetReviews() got nil response, want non-nil")
			}
		})
	}
}

func TestNewReviewService(t *testing.T) {
	type args struct {
		postgresRepo repository.RepositoryCollections
		helper       helper.HelperCollections
	}
	tests := []struct {
		name string
		args args
		want IReviewService
	}{
		{
			name: "NewReviewService",
			args: args{
				postgresRepo: repository.RepositoryCollections{
					ReviewRepo: repo_mocks.NewIReviewRepository(t),
				},
				helper: helper.HelperCollections{},
			},
			want: &reviewService{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewReviewService(tt.args.postgresRepo, tt.args.helper)
			if got == nil {
				t.Error("NewReviewService() got nil service")
			}
		})
	}
}
