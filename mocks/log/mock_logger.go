package log

import (
	"log/slog"
	"sync"
)

// MockLoggerClient implements ILoggerClient for testing
type MockLoggerClient struct {
	mu     sync.RWMutex
	logger *slog.Logger
}

// NewMockLoggerClient creates a new MockLoggerClient instance
func NewMockLoggerClient() *MockLoggerClient {
	return &MockLoggerClient{
		logger: slog.New(slog.NewTextHandler(nil, nil)), // Null handler for testing
	}
}

// GetLogger returns the mock logger
func (m *MockLoggerClient) GetLogger() *slog.Logger {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.logger
}

// SetLogger sets the mock logger
func (m *MockLoggerClient) SetLogger(l *slog.Logger) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.logger = l
}
