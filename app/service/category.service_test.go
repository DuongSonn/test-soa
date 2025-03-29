package service

import (
	"context"
	"reflect"
	"testing"

	"errors"

	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"

	"sondth-test_soa/app/entity"
	"sondth-test_soa/app/helper"
	"sondth-test_soa/app/model"
	"sondth-test_soa/app/repository"
	helper_mocks "sondth-test_soa/mocks/helper"
	repo_mocks "sondth-test_soa/mocks/repository"
)

var (
	testCategoryName        = "Test Category"
	testCategoryDesc        = "Test Description"
	testExistingCategory    = "Existing Category"
	testDefaultProductCount = int64(10)
	testPage                = 1
	testLimit               = 10
)

func Test_categoryService_Create(t *testing.T) {
	type args struct {
		ctx context.Context
		req *model.CreateCategoryRequest
	}
	type testCase struct {
		name    string
		s       *categoryService
		args    args
		want    *model.CreateCategoryResponse
		wantErr bool
		mock    func(repo *repo_mocks.ICategoryRepository)
	}

	ctx := context.Background()
	tests := []testCase{
		{
			name: "Create Success",
			s: &categoryService{
				postgresRepo: repository.RepositoryCollections{
					CategoryRepo: repo_mocks.NewICategoryRepository(t),
				},
				helper: helper_mocks.InitMockHelper(t),
			},
			args: args{
				ctx: ctx,
				req: &model.CreateCategoryRequest{
					Name: testCategoryName,
				},
			},
			want:    &model.CreateCategoryResponse{},
			wantErr: false,
			mock: func(repo *repo_mocks.ICategoryRepository) {
				repo.On("FindOneByFilter", ctx, mock.Anything, mock.MatchedBy(func(filter *repository.FindCategoryByFilter) bool {
					return filter.Name != nil && *filter.Name == testCategoryName
				})).Return(nil, gorm.ErrRecordNotFound)

				repo.On("Create", ctx, mock.Anything, mock.MatchedBy(func(category *entity.Category) bool {
					return category.Name == testCategoryName
				})).Return(nil)
			},
		},
		{
			name: "Category Already Exists",
			s: &categoryService{
				postgresRepo: repository.RepositoryCollections{
					CategoryRepo: repo_mocks.NewICategoryRepository(t),
				},
				helper: helper_mocks.InitMockHelper(t),
			},
			args: args{
				ctx: ctx,
				req: &model.CreateCategoryRequest{
					Name: testExistingCategory,
				},
			},
			want:    nil,
			wantErr: true,
			mock: func(repo *repo_mocks.ICategoryRepository) {
				repo.On("FindOneByFilter", ctx, mock.Anything, mock.MatchedBy(func(filter *repository.FindCategoryByFilter) bool {
					return filter.Name != nil && *filter.Name == testExistingCategory
				})).Return(&entity.Category{}, nil)
			},
		},
		{
			name: "FindOneByFilter Error",
			s: &categoryService{
				postgresRepo: repository.RepositoryCollections{
					CategoryRepo: repo_mocks.NewICategoryRepository(t),
				},
				helper: helper_mocks.InitMockHelper(t),
			},
			args: args{
				ctx: ctx,
				req: &model.CreateCategoryRequest{
					Name: testCategoryName,
				},
			},
			want:    nil,
			wantErr: true,
			mock: func(repo *repo_mocks.ICategoryRepository) {
				repo.On("FindOneByFilter", ctx, mock.Anything, mock.MatchedBy(func(filter *repository.FindCategoryByFilter) bool {
					return filter.Name != nil && *filter.Name == testCategoryName
				})).Return(nil, errors.New("database error"))
			},
		},
		{
			name: "Create Error",
			s: &categoryService{
				postgresRepo: repository.RepositoryCollections{
					CategoryRepo: repo_mocks.NewICategoryRepository(t),
				},
				helper: helper_mocks.InitMockHelper(t),
			},
			args: args{
				ctx: ctx,
				req: &model.CreateCategoryRequest{
					Name: testCategoryName,
				},
			},
			want:    nil,
			wantErr: true,
			mock: func(repo *repo_mocks.ICategoryRepository) {
				repo.On("FindOneByFilter", ctx, mock.Anything, mock.MatchedBy(func(filter *repository.FindCategoryByFilter) bool {
					return filter.Name != nil && *filter.Name == testCategoryName
				})).Return(nil, gorm.ErrRecordNotFound)

				repo.On("Create", ctx, mock.Anything, mock.MatchedBy(func(category *entity.Category) bool {
					return category.Name == testCategoryName
				})).Return(errors.New("database error"))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set up mocks before each test
			tt.mock(tt.s.postgresRepo.CategoryRepo.(*repo_mocks.ICategoryRepository))

			got, err := tt.s.Create(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("categoryService.Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("categoryService.Create() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_categoryService_GetCategories(t *testing.T) {
	type args struct {
		ctx context.Context
		req *model.GetCategoriesRequest
	}
	type testCase struct {
		name    string
		s       *categoryService
		args    args
		want    *model.GetCategoriesResponse
		wantErr bool
		mock    func(repo *repo_mocks.ICategoryRepository)
	}

	ctx := context.Background()
	tests := []testCase{
		{
			name: "Get Categories Success",
			s: &categoryService{
				postgresRepo: repository.RepositoryCollections{
					CategoryRepo: repo_mocks.NewICategoryRepository(t),
				},
				helper: helper_mocks.InitMockHelper(t),
			},
			args: args{
				ctx: ctx,
				req: &model.GetCategoriesRequest{
					Page:  &testPage,
					Limit: &testLimit,
				},
			},
			want: &model.GetCategoriesResponse{
				Count: 1,
				Result: []entity.Category{
					{
						Name:         "Category 1",
						Description:  &testCategoryDesc,
						ProductCount: testDefaultProductCount,
					},
				},
			},
			wantErr: false,
			mock: func(repo *repo_mocks.ICategoryRepository) {
				repo.On("FindManyByFilter", ctx, mock.Anything, mock.MatchedBy(func(filter *repository.FindCategoryByFilter) bool {
					return filter.Page != nil && *filter.Page == testPage && filter.Limit != nil && *filter.Limit == testLimit
				})).Return([]entity.Category{
					{
						Name:         "Category 1",
						Description:  &testCategoryDesc,
						ProductCount: testDefaultProductCount,
					},
				}, nil)
			},
		},
		{
			name: "Get Categories Without Pagination",
			s: &categoryService{
				postgresRepo: repository.RepositoryCollections{
					CategoryRepo: repo_mocks.NewICategoryRepository(t),
				},
				helper: helper_mocks.InitMockHelper(t),
			},
			args: args{
				ctx: ctx,
				req: &model.GetCategoriesRequest{},
			},
			want: &model.GetCategoriesResponse{
				Count: 1,
				Result: []entity.Category{
					{
						Name:         "Category 1",
						Description:  &testCategoryDesc,
						ProductCount: testDefaultProductCount,
					},
				},
			},
			wantErr: false,
			mock: func(repo *repo_mocks.ICategoryRepository) {
				repo.On("FindManyByFilter",
					mock.Anything,
					mock.Anything,
					mock.Anything).Return([]entity.Category{
					{
						Name:         "Category 1",
						Description:  &testCategoryDesc,
						ProductCount: testDefaultProductCount,
					},
				}, nil)
			},
		},
		{
			name: "FindManyByFilter Error",
			s: &categoryService{
				postgresRepo: repository.RepositoryCollections{
					CategoryRepo: repo_mocks.NewICategoryRepository(t),
				},
				helper: helper_mocks.InitMockHelper(t),
			},
			args: args{
				ctx: ctx,
				req: &model.GetCategoriesRequest{},
			},
			want:    nil,
			wantErr: true,
			mock: func(repo *repo_mocks.ICategoryRepository) {
				repo.On("FindManyByFilter", ctx, mock.Anything, mock.Anything).Return(nil, errors.New("database error"))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set up mocks before each test
			tt.mock(tt.s.postgresRepo.CategoryRepo.(*repo_mocks.ICategoryRepository))

			got, err := tt.s.GetCategories(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("categoryService.GetCategories() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("categoryService.GetCategories() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_categoryService_GetCategoriesSummary(t *testing.T) {
	type args struct {
		ctx context.Context
		req *model.GetCategoriesSummaryRequest
	}
	type testCase struct {
		name    string
		s       *categoryService
		args    args
		want    *model.GetCategoriesSummaryResponse
		wantErr bool
		mock    func(repo *repo_mocks.ICategoryRepository)
	}

	ctx := context.Background()
	tests := []testCase{
		{
			name: "Get Categories Summary Success",
			s: &categoryService{
				postgresRepo: repository.RepositoryCollections{
					CategoryRepo: repo_mocks.NewICategoryRepository(t),
				},
				helper: helper_mocks.InitMockHelper(t),
			},
			args: args{
				ctx: ctx,
				req: &model.GetCategoriesSummaryRequest{},
			},
			want: &model.GetCategoriesSummaryResponse{
				Categories: []entity.Category{
					{
						Name:         "Category 1",
						Description:  &testCategoryDesc,
						ProductCount: testDefaultProductCount,
					},
				},
			},
			wantErr: false,
			mock: func(repo *repo_mocks.ICategoryRepository) {
				repo.On("GetCategorySummary", ctx, mock.Anything).Return([]entity.Category{
					{
						Name:         "Category 1",
						Description:  &testCategoryDesc,
						ProductCount: testDefaultProductCount,
					},
				}, nil)
			},
		},
		{
			name: "GetCategorySummary Error",
			s: &categoryService{
				postgresRepo: repository.RepositoryCollections{
					CategoryRepo: repo_mocks.NewICategoryRepository(t),
				},
				helper: helper_mocks.InitMockHelper(t),
			},
			args: args{
				ctx: ctx,
				req: &model.GetCategoriesSummaryRequest{},
			},
			want:    nil,
			wantErr: true,
			mock: func(repo *repo_mocks.ICategoryRepository) {
				repo.On("GetCategorySummary", ctx, mock.Anything).Return(nil, errors.New("database error"))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set up mocks before each test
			tt.mock(tt.s.postgresRepo.CategoryRepo.(*repo_mocks.ICategoryRepository))

			got, err := tt.s.GetCategoriesSummary(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("categoryService.GetCategoriesSummary() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("categoryService.GetCategoriesSummary() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewCategoryService(t *testing.T) {
	type args struct {
		postgresRepo repository.RepositoryCollections
		helper       helper.HelperCollections
	}
	tests := []struct {
		name string
		args args
		want ICategoryService
	}{
		{
			name: "Init Success",
			args: args{
				postgresRepo: repo_mocks.InitMockRepository(t),
				helper:       helper_mocks.InitMockHelper(t),
			},
			want: &categoryService{
				postgresRepo: repo_mocks.InitMockRepository(t),
				helper:       helper_mocks.InitMockHelper(t),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewCategoryService(tt.args.postgresRepo, tt.args.helper); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewCategoryService() = %v, want %v", got, tt.want)
			}
		})
	}
}
