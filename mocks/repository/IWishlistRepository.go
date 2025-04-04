// Code generated by mockery v2.52.3. DO NOT EDIT.

package mocks

import (
	context "context"
	entity "sondth-test_soa/app/entity"

	gorm "gorm.io/gorm"

	mock "github.com/stretchr/testify/mock"

	repository "sondth-test_soa/app/repository"
)

// IWishlistRepository is an autogenerated mock type for the IWishlistRepository type
type IWishlistRepository struct {
	mock.Mock
}

// CountByFilter provides a mock function with given fields: ctx, tx, filter
func (_m *IWishlistRepository) CountByFilter(ctx context.Context, tx *gorm.DB, filter *repository.FindWishlistByFilter) (int64, error) {
	ret := _m.Called(ctx, tx, filter)

	if len(ret) == 0 {
		panic("no return value specified for CountByFilter")
	}

	var r0 int64
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *gorm.DB, *repository.FindWishlistByFilter) (int64, error)); ok {
		return rf(ctx, tx, filter)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *gorm.DB, *repository.FindWishlistByFilter) int64); ok {
		r0 = rf(ctx, tx, filter)
	} else {
		r0 = ret.Get(0).(int64)
	}

	if rf, ok := ret.Get(1).(func(context.Context, *gorm.DB, *repository.FindWishlistByFilter) error); ok {
		r1 = rf(ctx, tx, filter)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Create provides a mock function with given fields: ctx, tx, data
func (_m *IWishlistRepository) Create(ctx context.Context, tx *gorm.DB, data *entity.Wishlist) error {
	ret := _m.Called(ctx, tx, data)

	if len(ret) == 0 {
		panic("no return value specified for Create")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *gorm.DB, *entity.Wishlist) error); ok {
		r0 = rf(ctx, tx, data)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Delete provides a mock function with given fields: ctx, tx, data
func (_m *IWishlistRepository) Delete(ctx context.Context, tx *gorm.DB, data *entity.Wishlist) error {
	ret := _m.Called(ctx, tx, data)

	if len(ret) == 0 {
		panic("no return value specified for Delete")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *gorm.DB, *entity.Wishlist) error); ok {
		r0 = rf(ctx, tx, data)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// FindManyByFilter provides a mock function with given fields: ctx, tx, filter
func (_m *IWishlistRepository) FindManyByFilter(ctx context.Context, tx *gorm.DB, filter *repository.FindWishlistByFilter) ([]entity.Wishlist, error) {
	ret := _m.Called(ctx, tx, filter)

	if len(ret) == 0 {
		panic("no return value specified for FindManyByFilter")
	}

	var r0 []entity.Wishlist
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *gorm.DB, *repository.FindWishlistByFilter) ([]entity.Wishlist, error)); ok {
		return rf(ctx, tx, filter)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *gorm.DB, *repository.FindWishlistByFilter) []entity.Wishlist); ok {
		r0 = rf(ctx, tx, filter)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]entity.Wishlist)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *gorm.DB, *repository.FindWishlistByFilter) error); ok {
		r1 = rf(ctx, tx, filter)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FindOneByFilter provides a mock function with given fields: ctx, tx, filter
func (_m *IWishlistRepository) FindOneByFilter(ctx context.Context, tx *gorm.DB, filter *repository.FindWishlistByFilter) (*entity.Wishlist, error) {
	ret := _m.Called(ctx, tx, filter)

	if len(ret) == 0 {
		panic("no return value specified for FindOneByFilter")
	}

	var r0 *entity.Wishlist
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *gorm.DB, *repository.FindWishlistByFilter) (*entity.Wishlist, error)); ok {
		return rf(ctx, tx, filter)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *gorm.DB, *repository.FindWishlistByFilter) *entity.Wishlist); ok {
		r0 = rf(ctx, tx, filter)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*entity.Wishlist)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *gorm.DB, *repository.FindWishlistByFilter) error); ok {
		r1 = rf(ctx, tx, filter)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewIWishlistRepository creates a new instance of IWishlistRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewIWishlistRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *IWishlistRepository {
	mock := &IWishlistRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
