package logger

import (
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestProvideLogger(t *testing.T) {
	// Create a mock container
	mockContainer := &mockContainer{
		items: make(map[string]interface{}),
	}

	// Test first call - should create new logger
	logger1 := ProvideLogger(mockContainer)

	assert.NotNil(t, logger1, "Expected logger to be created")

	// Verify logger is configured correctly
	assert.Equal(t, logrus.InfoLevel, logger1.Level, "Expected logger level to be InfoLevel")

	// Verify logger was added to container
	assert.NotNil(t, mockContainer.items[diKey], "Expected logger to be added to container")

	// Test second call - should return same logger
	logger2 := ProvideLogger(mockContainer)

	assert.Equal(t, logger1, logger2, "Expected same logger instance to be returned")
}

func TestProvideLogger_ExistingLogger(t *testing.T) {
	// Create a mock container with existing logger
	existingLogger := logrus.New()
	mockContainer := &mockContainer{
		items: map[string]interface{}{
			diKey: existingLogger,
		},
	}

	// Test call with existing logger
	logger := ProvideLogger(mockContainer)

	assert.Equal(t, existingLogger, logger, "Expected existing logger to be returned")
}

// Mock container for testing
type mockContainer struct {
	items map[string]interface{}
}

func (m *mockContainer) Add(key string, value interface{}) error {
	m.items[key] = value
	return nil
}

func (m *mockContainer) Get(key string) interface{} {
	return m.items[key]
}
