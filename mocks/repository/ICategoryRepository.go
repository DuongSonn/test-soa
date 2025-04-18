// Code generated by mockery v2.52.3. DO NOT EDIT.

package mocks

import (
	context "context"
	entity "sondth-test_soa/app/entity"

	gorm "gorm.io/gorm"

	mock "github.com/stretchr/testify/mock"

	repository "sondth-test_soa/app/repository"
)

// ICategoryRepository is an autogenerated mock type for the ICategoryRepository type
type ICategoryRepository struct {
	mock.Mock
}

// Create provides a mock function with given fields: ctx, tx, data
func (_m *ICategoryRepository) Create(ctx context.Context, tx *gorm.DB, data *entity.Category) error {
	ret := _m.Called(ctx, tx, data)

	if len(ret) == 0 {
		panic("no return value specified for Create")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *gorm.DB, *entity.Category) error); ok {
		r0 = rf(ctx, tx, data)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// FindManyByFilter provides a mock function with given fields: ctx, tx, filter
func (_m *ICategoryRepository) FindManyByFilter(ctx context.Context, tx *gorm.DB, filter *repository.FindCategoryByFilter) ([]entity.Category, error) {
	ret := _m.Called(ctx, tx, filter)

	if len(ret) == 0 {
		panic("no return value specified for FindManyByFilter")
	}

	var r0 []entity.Category
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *gorm.DB, *repository.FindCategoryByFilter) ([]entity.Category, error)); ok {
		return rf(ctx, tx, filter)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *gorm.DB, *repository.FindCategoryByFilter) []entity.Category); ok {
		r0 = rf(ctx, tx, filter)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]entity.Category)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *gorm.DB, *repository.FindCategoryByFilter) error); ok {
		r1 = rf(ctx, tx, filter)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FindOneByFilter provides a mock function with given fields: ctx, tx, filter
func (_m *ICategoryRepository) FindOneByFilter(ctx context.Context, tx *gorm.DB, filter *repository.FindCategoryByFilter) (*entity.Category, error) {
	ret := _m.Called(ctx, tx, filter)

	if len(ret) == 0 {
		panic("no return value specified for FindOneByFilter")
	}

	var r0 *entity.Category
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *gorm.DB, *repository.FindCategoryByFilter) (*entity.Category, error)); ok {
		return rf(ctx, tx, filter)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *gorm.DB, *repository.FindCategoryByFilter) *entity.Category); ok {
		r0 = rf(ctx, tx, filter)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*entity.Category)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *gorm.DB, *repository.FindCategoryByFilter) error); ok {
		r1 = rf(ctx, tx, filter)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetCategorySummary provides a mock function with given fields: ctx, tx
func (_m *ICategoryRepository) GetCategorySummary(ctx context.Context, tx *gorm.DB) ([]entity.Category, error) {
	ret := _m.Called(ctx, tx)

	if len(ret) == 0 {
		panic("no return value specified for GetCategorySummary")
	}

	var r0 []entity.Category
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *gorm.DB) ([]entity.Category, error)); ok {
		return rf(ctx, tx)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *gorm.DB) []entity.Category); ok {
		r0 = rf(ctx, tx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]entity.Category)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *gorm.DB) error); ok {
		r1 = rf(ctx, tx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewICategoryRepository creates a new instance of ICategoryRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewICategoryRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *ICategoryRepository {
	mock := &ICategoryRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
