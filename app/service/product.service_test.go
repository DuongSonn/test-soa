package service

import (
	"context"
	"reflect"
	"sondth-test_soa/app/entity"
	"sondth-test_soa/app/helper"
	"sondth-test_soa/app/model"
	"sondth-test_soa/app/repository"
	helper_mocks "sondth-test_soa/mocks/helper"
	repo_mocks "sondth-test_soa/mocks/repository"
	"sondth-test_soa/package/errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

var (
	testProductID       = uuid.New()
	testProductName     = "Test Product"
	testProductDesc     = "Test Description"
	testProductPrice    = float64(100)
	testProductQuantity = uint64(10)
	testCategoryID      = uuid.New()
)

func Test_productService_Create(t *testing.T) {
	type args struct {
		ctx context.Context
		req *model.CreateProductRequest
	}
	type testCase struct {
		name    string
		s       *productService
		args    args
		want    *model.CreateProductResponse
		wantErr bool
		mock    func(repo *repo_mocks.IProductRepository, categoryHelper *helper_mocks.ICategoryHelper)
	}

	ctx := context.Background()
	tests := []testCase{
		{
			name: "Create Success",
			s: &productService{
				postgresRepo: repository.RepositoryCollections{
					ProductRepo: repo_mocks.NewIProductRepository(t),
				},
				helper: helper.HelperCollections{
					CategoryHelper: helper_mocks.NewICategoryHelper(t),
				},
			},
			args: args{
				ctx: ctx,
				req: &model.CreateProductRequest{
					Name:        testProductName,
					Description: &testProductDesc,
					Price:       testProductPrice,
					Quantity:    testProductQuantity,
					CategoryID:  testCategoryID,
				},
			},
			want:    &model.CreateProductResponse{},
			wantErr: false,
			mock: func(repo *repo_mocks.IProductRepository, categoryHelper *helper_mocks.ICategoryHelper) {
				// Mock category validation
				categoryHelper.On("ValidateCategoryID", ctx, testCategoryID).Return(&entity.Category{}, nil).Once()

				// Mock product name check
				repo.On("FindOneByFilter", ctx, mock.Anything, mock.MatchedBy(func(filter *repository.FindProductByFilter) bool {
					return filter.Name != nil && *filter.Name == testProductName
				})).Return(nil, gorm.ErrRecordNotFound).Once()

				// Mock product creation
				repo.On("Create", ctx, mock.Anything, mock.MatchedBy(func(product *entity.Product) bool {
					return product.Name == testProductName &&
						*product.Description == testProductDesc &&
						product.Price == testProductPrice &&
						product.Quantity == testProductQuantity &&
						product.CategoryID == testCategoryID
				})).Return(nil).Once()
			},
		},
		{
			name: "Category Not Found",
			s: &productService{
				postgresRepo: repository.RepositoryCollections{
					ProductRepo: repo_mocks.NewIProductRepository(t),
				},
				helper: helper.HelperCollections{
					CategoryHelper: helper_mocks.NewICategoryHelper(t),
				},
			},
			args: args{
				ctx: ctx,
				req: &model.CreateProductRequest{
					Name:        testProductName,
					Description: &testProductDesc,
					Price:       testProductPrice,
					Quantity:    testProductQuantity,
					CategoryID:  testCategoryID,
				},
			},
			want:    nil,
			wantErr: true,
			mock: func(repo *repo_mocks.IProductRepository, categoryHelper *helper_mocks.ICategoryHelper) {
				// Mock category validation failure
				categoryHelper.On("ValidateCategoryID", ctx, testCategoryID).Return(nil, errors.New(errors.ErrCodeCategoryNotFound)).Once()
			},
		},
		{
			name: "Product Already Exists",
			s: &productService{
				postgresRepo: repository.RepositoryCollections{
					ProductRepo: repo_mocks.NewIProductRepository(t),
				},
				helper: helper.HelperCollections{
					CategoryHelper: helper_mocks.NewICategoryHelper(t),
				},
			},
			args: args{
				ctx: ctx,
				req: &model.CreateProductRequest{
					Name:        testProductName,
					Description: &testProductDesc,
					Price:       testProductPrice,
					Quantity:    testProductQuantity,
					CategoryID:  testCategoryID,
				},
			},
			want:    nil,
			wantErr: true,
			mock: func(repo *repo_mocks.IProductRepository, categoryHelper *helper_mocks.ICategoryHelper) {
				// Mock category validation
				categoryHelper.On("ValidateCategoryID", ctx, testCategoryID).Return(&entity.Category{}, nil).Once()

				// Mock product name check - product exists
				repo.On("FindOneByFilter", ctx, mock.Anything, mock.MatchedBy(func(filter *repository.FindProductByFilter) bool {
					return filter.Name != nil && *filter.Name == testProductName
				})).Return(&entity.Product{}, nil).Once()
			},
		},
		{
			name: "Create Error",
			s: &productService{
				postgresRepo: repository.RepositoryCollections{
					ProductRepo: repo_mocks.NewIProductRepository(t),
				},
				helper: helper.HelperCollections{
					CategoryHelper: helper_mocks.NewICategoryHelper(t),
				},
			},
			args: args{
				ctx: ctx,
				req: &model.CreateProductRequest{
					Name:        testProductName,
					Description: &testProductDesc,
					Price:       testProductPrice,
					Quantity:    testProductQuantity,
					CategoryID:  testCategoryID,
				},
			},
			want:    nil,
			wantErr: true,
			mock: func(repo *repo_mocks.IProductRepository, categoryHelper *helper_mocks.ICategoryHelper) {
				// Mock category validation
				categoryHelper.On("ValidateCategoryID", ctx, testCategoryID).Return(&entity.Category{}, nil).Once()

				// Mock product name check
				repo.On("FindOneByFilter", ctx, mock.Anything, mock.MatchedBy(func(filter *repository.FindProductByFilter) bool {
					return filter.Name != nil && *filter.Name == testProductName
				})).Return(nil, gorm.ErrRecordNotFound).Once()

				// Mock product creation error
				repo.On("Create", ctx, mock.Anything, mock.MatchedBy(func(product *entity.Product) bool {
					return product.Name == testProductName &&
						*product.Description == testProductDesc &&
						product.Price == testProductPrice &&
						product.Quantity == testProductQuantity &&
						product.CategoryID == testCategoryID
				})).Return(errors.New(errors.ErrCodeInternalServerError)).Once()
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set up mocks
			tt.mock(
				tt.s.postgresRepo.ProductRepo.(*repo_mocks.IProductRepository),
				tt.s.helper.CategoryHelper.(*helper_mocks.ICategoryHelper),
			)

			got, err := tt.s.Create(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("productService.Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("productService.Create() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_productService_Update(t *testing.T) {
	type args struct {
		ctx context.Context
		req *model.UpdateProductRequest
	}
	type testCase struct {
		name    string
		s       *productService
		args    args
		want    *model.UpdateProductResponse
		wantErr bool
		mock    func(repo *repo_mocks.IProductRepository, productHelper *helper_mocks.IProductHelper, categoryHelper *helper_mocks.ICategoryHelper)
	}

	ctx := context.Background()
	existingProduct := &entity.Product{
		ID:          testProductID,
		Name:        testProductName,
		Description: &testProductDesc,
		Price:       testProductPrice,
		Quantity:    testProductQuantity,
		CategoryID:  testCategoryID,
	}

	tests := []testCase{
		{
			name: "Update Success - Same Category",
			s: &productService{
				postgresRepo: repository.RepositoryCollections{
					ProductRepo: repo_mocks.NewIProductRepository(t),
				},
				helper: helper.HelperCollections{
					ProductHelper:  helper_mocks.NewIProductHelper(t),
					CategoryHelper: helper_mocks.NewICategoryHelper(t),
				},
			},
			args: args{
				ctx: ctx,
				req: &model.UpdateProductRequest{
					ID: testProductID,
					CreateProductRequest: model.CreateProductRequest{
						Name:        testProductName,
						Description: &testProductDesc,
						Price:       testProductPrice,
						Quantity:    testProductQuantity,
						CategoryID:  testCategoryID,
					},
				},
			},
			want:    &model.UpdateProductResponse{},
			wantErr: false,
			mock: func(repo *repo_mocks.IProductRepository, productHelper *helper_mocks.IProductHelper, categoryHelper *helper_mocks.ICategoryHelper) {
				// Mock product validation
				productHelper.On("ValidateProductID", ctx, testProductID).Return(existingProduct, nil).Once()

				// Mock product update with exact product match
				repo.On("Update", ctx, mock.Anything, mock.MatchedBy(func(product *entity.Product) bool {
					return product.ID == testProductID &&
						product.Name == testProductName &&
						*product.Description == testProductDesc &&
						product.Price == testProductPrice &&
						product.Quantity == testProductQuantity &&
						product.CategoryID == testCategoryID
				})).Return(nil).Once()
			},
		},
		{
			name: "Update Success - Different Category",
			s: &productService{
				postgresRepo: repository.RepositoryCollections{
					ProductRepo: repo_mocks.NewIProductRepository(t),
				},
				helper: helper.HelperCollections{
					ProductHelper:  helper_mocks.NewIProductHelper(t),
					CategoryHelper: helper_mocks.NewICategoryHelper(t),
				},
			},
			args: args{
				ctx: ctx,
				req: &model.UpdateProductRequest{
					ID: testProductID,
					CreateProductRequest: model.CreateProductRequest{
						Name:        testProductName,
						Description: &testProductDesc,
						Price:       testProductPrice,
						Quantity:    testProductQuantity,
						CategoryID:  uuid.New(), // Different category ID
					},
				},
			},
			want:    &model.UpdateProductResponse{},
			wantErr: false,
			mock: func(repo *repo_mocks.IProductRepository, productHelper *helper_mocks.IProductHelper, categoryHelper *helper_mocks.ICategoryHelper) {
				// Mock product validation
				productHelper.On("ValidateProductID", ctx, testProductID).Return(existingProduct, nil).Once()

				// Mock category validation for new category
				categoryHelper.On("ValidateCategoryID", ctx, mock.Anything).Return(&entity.Category{}, nil).Once()

				// Mock product update
				repo.On("Update", ctx, mock.Anything, mock.MatchedBy(func(product *entity.Product) bool {
					return product.Name == testProductName &&
						*product.Description == testProductDesc &&
						product.Price == testProductPrice &&
						product.Quantity == testProductQuantity
				})).Return(nil).Once()
			},
		},
		{
			name: "Product Not Found",
			s: &productService{
				postgresRepo: repository.RepositoryCollections{
					ProductRepo: repo_mocks.NewIProductRepository(t),
				},
				helper: helper.HelperCollections{
					ProductHelper:  helper_mocks.NewIProductHelper(t),
					CategoryHelper: helper_mocks.NewICategoryHelper(t),
				},
			},
			args: args{
				ctx: ctx,
				req: &model.UpdateProductRequest{
					ID: testProductID,
					CreateProductRequest: model.CreateProductRequest{
						Name:        testProductName,
						Description: &testProductDesc,
						Price:       testProductPrice,
						Quantity:    testProductQuantity,
						CategoryID:  testCategoryID,
					},
				},
			},
			want:    nil,
			wantErr: true,
			mock: func(repo *repo_mocks.IProductRepository, productHelper *helper_mocks.IProductHelper, categoryHelper *helper_mocks.ICategoryHelper) {
				// Mock product validation failure
				productHelper.On("ValidateProductID", ctx, testProductID).Return(nil, errors.New(errors.ErrCodeProductNotFound)).Once()
			},
		},
		{
			name: "New Category Not Found",
			s: &productService{
				postgresRepo: repository.RepositoryCollections{
					ProductRepo: repo_mocks.NewIProductRepository(t),
				},
				helper: helper.HelperCollections{
					ProductHelper:  helper_mocks.NewIProductHelper(t),
					CategoryHelper: helper_mocks.NewICategoryHelper(t),
				},
			},
			args: args{
				ctx: ctx,
				req: &model.UpdateProductRequest{
					ID: testProductID,
					CreateProductRequest: model.CreateProductRequest{
						Name:        testProductName,
						Description: &testProductDesc,
						Price:       testProductPrice,
						Quantity:    testProductQuantity,
						CategoryID:  uuid.New(), // Different category ID
					},
				},
			},
			want:    nil,
			wantErr: true,
			mock: func(repo *repo_mocks.IProductRepository, productHelper *helper_mocks.IProductHelper, categoryHelper *helper_mocks.ICategoryHelper) {
				// Mock product validation
				productHelper.On("ValidateProductID", ctx, testProductID).Return(existingProduct, nil).Once()

				// Mock category validation failure
				categoryHelper.On("ValidateCategoryID", ctx, mock.Anything).Return(nil, errors.New(errors.ErrCodeCategoryNotFound)).Once()
			},
		},
		{
			name: "Update Error",
			s: &productService{
				postgresRepo: repository.RepositoryCollections{
					ProductRepo: repo_mocks.NewIProductRepository(t),
				},
				helper: helper.HelperCollections{
					ProductHelper:  helper_mocks.NewIProductHelper(t),
					CategoryHelper: helper_mocks.NewICategoryHelper(t),
				},
			},
			args: args{
				ctx: ctx,
				req: &model.UpdateProductRequest{
					ID: testProductID,
					CreateProductRequest: model.CreateProductRequest{
						Name:        testProductName,
						Description: &testProductDesc,
						Price:       testProductPrice,
						Quantity:    testProductQuantity,
						CategoryID:  testCategoryID,
					},
				},
			},
			want:    nil,
			wantErr: true,
			mock: func(repo *repo_mocks.IProductRepository, productHelper *helper_mocks.IProductHelper, categoryHelper *helper_mocks.ICategoryHelper) {
				// Mock product validation
				productHelper.On("ValidateProductID", ctx, testProductID).Return(existingProduct, nil).Once()
				categoryHelper.On("ValidateCategoryID", ctx, testCategoryID).Return(&entity.Category{}, nil).Once()

				// Mock update error
				repo.On("Update", ctx, mock.Anything, mock.Anything).Return(errors.New(errors.ErrCodeInternalServerError)).Once()
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set up mocks
			tt.mock(
				tt.s.postgresRepo.ProductRepo.(*repo_mocks.IProductRepository),
				tt.s.helper.ProductHelper.(*helper_mocks.IProductHelper),
				tt.s.helper.CategoryHelper.(*helper_mocks.ICategoryHelper),
			)

			got, err := tt.s.Update(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("productService.Update() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("productService.Update() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_productService_Delete(t *testing.T) {
	type args struct {
		ctx context.Context
		req *model.DeleteProductRequest
	}
	type testCase struct {
		name    string
		s       *productService
		args    args
		want    *model.DeleteProductResponse
		wantErr bool
		mock    func(repo *repo_mocks.IProductRepository, productHelper *helper_mocks.IProductHelper)
	}

	ctx := context.Background()
	existingProduct := &entity.Product{
		Name:        testProductName,
		Description: &testProductDesc,
		Price:       testProductPrice,
		Quantity:    testProductQuantity,
		CategoryID:  testCategoryID,
	}

	tests := []testCase{
		{
			name: "Delete Success",
			s: &productService{
				postgresRepo: repository.RepositoryCollections{
					ProductRepo: repo_mocks.NewIProductRepository(t),
				},
				helper: helper.HelperCollections{
					ProductHelper: helper_mocks.NewIProductHelper(t),
				},
			},
			args: args{
				ctx: ctx,
				req: &model.DeleteProductRequest{
					ID: testProductID,
				},
			},
			want:    &model.DeleteProductResponse{},
			wantErr: false,
			mock: func(repo *repo_mocks.IProductRepository, productHelper *helper_mocks.IProductHelper) {
				// Mock product validation
				productHelper.On("ValidateProductID", ctx, testProductID).Return(existingProduct, nil).Once()

				// Mock product deletion
				repo.On("Delete", ctx, mock.Anything, mock.MatchedBy(func(product *entity.Product) bool {
					return product.Name == testProductName &&
						*product.Description == testProductDesc &&
						product.Price == testProductPrice &&
						product.Quantity == testProductQuantity &&
						product.CategoryID == testCategoryID
				})).Return(nil).Once()
			},
		},
		{
			name: "Product Not Found",
			s: &productService{
				postgresRepo: repository.RepositoryCollections{
					ProductRepo: repo_mocks.NewIProductRepository(t),
				},
				helper: helper.HelperCollections{
					ProductHelper: helper_mocks.NewIProductHelper(t),
				},
			},
			args: args{
				ctx: ctx,
				req: &model.DeleteProductRequest{
					ID: testProductID,
				},
			},
			want:    nil,
			wantErr: true,
			mock: func(repo *repo_mocks.IProductRepository, productHelper *helper_mocks.IProductHelper) {
				// Mock product validation failure
				productHelper.On("ValidateProductID", ctx, testProductID).Return(nil, errors.New(errors.ErrCodeProductNotFound)).Once()
			},
		},
		{
			name: "Delete Error",
			s: &productService{
				postgresRepo: repository.RepositoryCollections{
					ProductRepo: repo_mocks.NewIProductRepository(t),
				},
				helper: helper.HelperCollections{
					ProductHelper: helper_mocks.NewIProductHelper(t),
				},
			},
			args: args{
				ctx: ctx,
				req: &model.DeleteProductRequest{
					ID: testProductID,
				},
			},
			want:    nil,
			wantErr: true,
			mock: func(repo *repo_mocks.IProductRepository, productHelper *helper_mocks.IProductHelper) {
				// Mock product validation
				productHelper.On("ValidateProductID", ctx, testProductID).Return(existingProduct, nil).Once()

				// Mock delete error
				repo.On("Delete", ctx, mock.Anything, mock.Anything).Return(errors.New(errors.ErrCodeInternalServerError)).Once()
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set up mocks
			tt.mock(
				tt.s.postgresRepo.ProductRepo.(*repo_mocks.IProductRepository),
				tt.s.helper.ProductHelper.(*helper_mocks.IProductHelper),
			)

			got, err := tt.s.Delete(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("productService.Delete() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("productService.Delete() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_productService_GetProducts(t *testing.T) {
	type args struct {
		ctx     context.Context
		request *model.GetProductRequest
	}

	type testCase struct {
		name    string
		s       *productService
		args    args
		want    *model.GetProductResponse
		wantErr bool
		mock    func(repo *repo_mocks.IProductRepository, ctx context.Context)
	}

	ctx := context.Background()
	products := []entity.Product{
		{
			Name:        testProductName,
			Description: &testProductDesc,
			Price:       testProductPrice,
			Quantity:    testProductQuantity,
			CategoryID:  testCategoryID,
		},
		{
			Name:        testProductName + "2",
			Description: &testProductDesc,
			Price:       testProductPrice * 2,
			Quantity:    testProductQuantity * 2,
			CategoryID:  testCategoryID,
		},
	}

	tests := []testCase{
		{
			name: "Get Products Success",
			s: &productService{
				postgresRepo: repository.RepositoryCollections{
					ProductRepo: repo_mocks.NewIProductRepository(t),
				},
			},
			args: args{
				ctx: ctx,
				request: &model.GetProductRequest{
					Name:        &testProductName,
					CategoryIDs: []uuid.UUID{testCategoryID},
					Page:        &testPage,
					Limit:       &testLimit,
				},
			},
			want: &model.GetProductResponse{
				Count:  int64(len(products)),
				Result: products,
			},
			wantErr: false,
			mock: func(repo *repo_mocks.IProductRepository, ctx context.Context) {
				// Mock find products
				repo.On("FindManyByFilter", mock.MatchedBy(func(c context.Context) bool {
					return true // Accept any context since we just want to verify the call
				}), mock.Anything, mock.MatchedBy(func(filter *repository.FindProductByFilter) bool {
					return filter.Name != nil && *filter.Name == testProductName &&
						len(filter.CategoryIDs) == 1 && filter.CategoryIDs[0] == testCategoryID &&
						filter.Page != nil && *filter.Page == 1 &&
						filter.Limit != nil && *filter.Limit == 10
				})).Return(products, nil).Once()

				// Mock count products
				repo.On("CountByFilter", mock.MatchedBy(func(c context.Context) bool {
					return true // Accept any context since we just want to verify the call
				}), mock.Anything, mock.MatchedBy(func(filter *repository.FindProductByFilter) bool {
					return filter.Name != nil && *filter.Name == testProductName &&
						len(filter.CategoryIDs) == 1 && filter.CategoryIDs[0] == testCategoryID &&
						filter.Page != nil && *filter.Page == 1 &&
						filter.Limit != nil && *filter.Limit == 10
				})).Return(int64(len(products)), nil).Once()
			},
		},
		{
			name: "Get Products Success - Empty Filter",
			s: &productService{
				postgresRepo: repository.RepositoryCollections{
					ProductRepo: repo_mocks.NewIProductRepository(t),
				},
			},
			args: args{
				ctx: ctx,
				request: &model.GetProductRequest{
					Page:  &testPage,
					Limit: &testLimit,
				},
			},
			want: &model.GetProductResponse{
				Count:  int64(len(products)),
				Result: products,
			},
			wantErr: false,
			mock: func(repo *repo_mocks.IProductRepository, ctx context.Context) {
				// Mock find products with nil filter
				repo.On("FindManyByFilter", mock.MatchedBy(func(c context.Context) bool {
					return true
				}), mock.Anything, mock.MatchedBy(func(filter *repository.FindProductByFilter) bool {
					return filter.Name == nil &&
						len(filter.CategoryIDs) == 0 &&
						filter.Page != nil && *filter.Page == 1 &&
						filter.Limit != nil && *filter.Limit == 10
				})).Return(products, nil).Once()

				// Mock count products with nil filter
				repo.On("CountByFilter", mock.MatchedBy(func(c context.Context) bool {
					return true
				}), mock.Anything, mock.MatchedBy(func(filter *repository.FindProductByFilter) bool {
					return filter.Name == nil &&
						len(filter.CategoryIDs) == 0 &&
						filter.Page != nil && *filter.Page == 1 &&
						filter.Limit != nil && *filter.Limit == 10
				})).Return(int64(len(products)), nil).Once()
			},
		},
		{
			name: "Get Products Error",
			s: &productService{
				postgresRepo: repository.RepositoryCollections{
					ProductRepo: repo_mocks.NewIProductRepository(t),
				},
			},
			args: args{
				ctx:     ctx,
				request: &model.GetProductRequest{},
			},
			want:    nil,
			wantErr: true,
			mock: func(repo *repo_mocks.IProductRepository, ctx context.Context) {
				// Mock find products error
				repo.On("FindManyByFilter", mock.MatchedBy(func(c context.Context) bool {
					return true
				}), mock.Anything, mock.Anything).Return(nil, errors.New(errors.ErrCodeInternalServerError)).Once()
				repo.On("CountByFilter", mock.MatchedBy(func(c context.Context) bool {
					return true
				}), mock.Anything, mock.Anything).Return(int64(0), errors.New(errors.ErrCodeInternalServerError)).Once()
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a new context with cancellation for each test
			ctx, cancel := context.WithCancel(tt.args.ctx)
			defer cancel()
			tt.args.ctx = ctx

			// Setup mocks
			tt.mock(tt.s.postgresRepo.ProductRepo.(*repo_mocks.IProductRepository), ctx)

			got, err := tt.s.GetProducts(tt.args.ctx, tt.args.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("productService.GetProducts() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("productService.GetProducts() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewProductService(t *testing.T) {
	type args struct {
		postgresRepo repository.RepositoryCollections
		helper       helper.HelperCollections
	}
	tests := []struct {
		name string
		args args
		want IProductService
	}{
		{
			name: "NewProductService",
			args: args{
				postgresRepo: repo_mocks.InitMockRepository(t),
				helper:       helper_mocks.InitMockHelper(t),
			},
			want: &productService{
				postgresRepo: repo_mocks.InitMockRepository(t),
				helper:       helper_mocks.InitMockHelper(t),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewProductService(tt.args.postgresRepo, tt.args.helper); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewProductService() = %v, want %v", got, tt.want)
			}
		})
	}
}
