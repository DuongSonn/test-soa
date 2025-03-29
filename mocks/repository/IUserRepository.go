// Code generated by mockery v2.52.3. DO NOT EDIT.

package mocks

import (
	context "context"
	entity "sondth-test_soa/app/entity"

	gorm "gorm.io/gorm"

	mock "github.com/stretchr/testify/mock"

	repository "sondth-test_soa/app/repository"
)

// IUserRepository is an autogenerated mock type for the IUserRepository type
type IUserRepository struct {
	mock.Mock
}

// CountByFilter provides a mock function with given fields: ctx, tx, filter
func (_m *IUserRepository) CountByFilter(ctx context.Context, tx *gorm.DB, filter *repository.FindUserByFilter) (int64, error) {
	ret := _m.Called(ctx, tx, filter)

	if len(ret) == 0 {
		panic("no return value specified for CountByFilter")
	}

	var r0 int64
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *gorm.DB, *repository.FindUserByFilter) (int64, error)); ok {
		return rf(ctx, tx, filter)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *gorm.DB, *repository.FindUserByFilter) int64); ok {
		r0 = rf(ctx, tx, filter)
	} else {
		r0 = ret.Get(0).(int64)
	}

	if rf, ok := ret.Get(1).(func(context.Context, *gorm.DB, *repository.FindUserByFilter) error); ok {
		r1 = rf(ctx, tx, filter)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Create provides a mock function with given fields: ctx, tx, data
func (_m *IUserRepository) Create(ctx context.Context, tx *gorm.DB, data *entity.User) error {
	ret := _m.Called(ctx, tx, data)

	if len(ret) == 0 {
		panic("no return value specified for Create")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *gorm.DB, *entity.User) error); ok {
		r0 = rf(ctx, tx, data)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// FindManyByFilter provides a mock function with given fields: ctx, tx, filter
func (_m *IUserRepository) FindManyByFilter(ctx context.Context, tx *gorm.DB, filter *repository.FindUserByFilter) ([]entity.User, error) {
	ret := _m.Called(ctx, tx, filter)

	if len(ret) == 0 {
		panic("no return value specified for FindManyByFilter")
	}

	var r0 []entity.User
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *gorm.DB, *repository.FindUserByFilter) ([]entity.User, error)); ok {
		return rf(ctx, tx, filter)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *gorm.DB, *repository.FindUserByFilter) []entity.User); ok {
		r0 = rf(ctx, tx, filter)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]entity.User)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *gorm.DB, *repository.FindUserByFilter) error); ok {
		r1 = rf(ctx, tx, filter)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FindOneByFilter provides a mock function with given fields: ctx, tx, filter
func (_m *IUserRepository) FindOneByFilter(ctx context.Context, tx *gorm.DB, filter *repository.FindUserByFilter) (*entity.User, error) {
	ret := _m.Called(ctx, tx, filter)

	if len(ret) == 0 {
		panic("no return value specified for FindOneByFilter")
	}

	var r0 *entity.User
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *gorm.DB, *repository.FindUserByFilter) (*entity.User, error)); ok {
		return rf(ctx, tx, filter)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *gorm.DB, *repository.FindUserByFilter) *entity.User); ok {
		r0 = rf(ctx, tx, filter)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*entity.User)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *gorm.DB, *repository.FindUserByFilter) error); ok {
		r1 = rf(ctx, tx, filter)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Update provides a mock function with given fields: ctx, tx, data
func (_m *IUserRepository) Update(ctx context.Context, tx *gorm.DB, data *entity.User) error {
	ret := _m.Called(ctx, tx, data)

	if len(ret) == 0 {
		panic("no return value specified for Update")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *gorm.DB, *entity.User) error); ok {
		r0 = rf(ctx, tx, data)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewIUserRepository creates a new instance of IUserRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewIUserRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *IUserRepository {
	mock := &IUserRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
