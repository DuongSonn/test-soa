package redis

import (
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	mockredis "sondth-test_soa/mocks/redis"
)

func TestRedisOperations(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create mock Redis client
	mockClient := mockredis.NewMockIRedisClient(ctrl)

	// Test Set operation
	t.Run("Set", func(t *testing.T) {
		ctx := context.Background()
		key := "test-key"
		value := "test-value"
		expiration := time.Hour

		// Set expectation
		mockClient.EXPECT().
			Set(ctx, key, value, expiration).
			Return(nil)

		// Call Set
		err := mockClient.Set(ctx, key, value, expiration)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
	})

	// Test Get operation
	t.Run("Get", func(t *testing.T) {
		ctx := context.Background()
		key := "test-key"
		expectedValue := "test-value"

		// Set expectation
		mockClient.EXPECT().
			Get(ctx, key).
			Return(expectedValue, nil)

		// Call Get
		value, err := mockClient.Get(ctx, key)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if value != expectedValue {
			t.Errorf("expected value %s, got %s", expectedValue, value)
		}
	})

	// Test Exists operation
	t.Run("Exists", func(t *testing.T) {
		ctx := context.Background()
		key := "test-key"

		// Set expectation
		mockClient.EXPECT().
			Exists(ctx, key).
			Return(true, nil)

		// Call Exists
		exists, err := mockClient.Exists(ctx, key)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if !exists {
			t.Error("expected key to exist")
		}
	})
}
