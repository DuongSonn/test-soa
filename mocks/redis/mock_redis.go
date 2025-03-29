// Code generated by MockGen. DO NOT EDIT.
// Source: redis.go

// Package mockredis is a generated GoMock package.
package mockredis

import (
	context "context"
	reflect "reflect"
	time "time"

	gomock "github.com/golang/mock/gomock"
)

// MockIRedisClient is a mock of IRedisClient interface.
type MockIRedisClient struct {
	ctrl     *gomock.Controller
	recorder *MockIRedisClientMockRecorder
}

// MockIRedisClientMockRecorder is the mock recorder for MockIRedisClient.
type MockIRedisClientMockRecorder struct {
	mock *MockIRedisClient
}

// NewMockIRedisClient creates a new mock instance.
func NewMockIRedisClient(ctrl *gomock.Controller) *MockIRedisClient {
	mock := &MockIRedisClient{ctrl: ctrl}
	mock.recorder = &MockIRedisClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIRedisClient) EXPECT() *MockIRedisClientMockRecorder {
	return m.recorder
}

// Delete mocks base method.
func (m *MockIRedisClient) Delete(ctx context.Context, key string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", ctx, key)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockIRedisClientMockRecorder) Delete(ctx, key interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockIRedisClient)(nil).Delete), ctx, key)
}

// Exists mocks base method.
func (m *MockIRedisClient) Exists(ctx context.Context, key string) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Exists", ctx, key)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Exists indicates an expected call of Exists.
func (mr *MockIRedisClientMockRecorder) Exists(ctx, key interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Exists", reflect.TypeOf((*MockIRedisClient)(nil).Exists), ctx, key)
}

// Get mocks base method.
func (m *MockIRedisClient) Get(ctx context.Context, key string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", ctx, key)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockIRedisClientMockRecorder) Get(ctx, key interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockIRedisClient)(nil).Get), ctx, key)
}

// Incr mocks base method.
func (m *MockIRedisClient) Incr(ctx context.Context, key string) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Incr", ctx, key)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Incr indicates an expected call of Incr.
func (mr *MockIRedisClientMockRecorder) Incr(ctx, key interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Incr", reflect.TypeOf((*MockIRedisClient)(nil).Incr), ctx, key)
}

// Set mocks base method.
func (m *MockIRedisClient) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Set", ctx, key, value, expiration)
	ret0, _ := ret[0].(error)
	return ret0
}

// Set indicates an expected call of Set.
func (mr *MockIRedisClientMockRecorder) Set(ctx, key, value, expiration interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Set", reflect.TypeOf((*MockIRedisClient)(nil).Set), ctx, key, value, expiration)
}

// SetNX mocks base method.
func (m *MockIRedisClient) SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetNX", ctx, key, value, expiration)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SetNX indicates an expected call of SetNX.
func (mr *MockIRedisClientMockRecorder) SetNX(ctx, key, value, expiration interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetNX", reflect.TypeOf((*MockIRedisClient)(nil).SetNX), ctx, key, value, expiration)
}
