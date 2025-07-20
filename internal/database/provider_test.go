package database

import (
	"testing"

	"github.com/arrowls/go-metrics/internal/config"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestProvideDatabasePool(t *testing.T) {
	// Create a mock container
	mockContainer := &mockContainer{
		items: make(map[string]interface{}),
	}

	// Mock config and logger
	mockConfig := config.ServerConfig{
		DatabaseDSN: "postgres://user:pass@localhost/testdb",
	}
	mockLogger := logrus.New()

	// Add mock config and logger to container
	mockContainer.items["server_config"] = mockConfig
	mockContainer.items["logger"] = mockLogger

	// Test first call - should attempt to create pool
	// The function calls logger.Fatal() which exits the program, so we can't easily test this
	// We'll skip this test for now since it requires a real database connection
	t.Skip("Skipping test that requires real database connection")
}

func TestProvideDatabasePool_ExistingPool(t *testing.T) {
	// Create a mock container with existing pool
	mockPool := &mockPool{}
	mockContainer := &mockContainer{
		items: map[string]interface{}{
			diKey: mockPool,
		},
	}

	// Mock config and logger
	mockConfig := config.ServerConfig{
		DatabaseDSN: "postgres://user:pass@localhost/testdb",
	}
	mockLogger := logrus.New()

	// Add mock config and logger to container
	mockContainer.items["server_config"] = mockConfig
	mockContainer.items["logger"] = mockLogger

	// Test call with existing pool
	pool, err := ProvideDatabasePool(mockContainer)

	assert.NoError(t, err, "Expected no error")
	assert.NotNil(t, pool, "Expected pool to be returned")
}

func TestIsConnectionException(t *testing.T) {
	tests := []struct {
		name     string
		code     string
		expected bool
	}{
		{"ConnectionException", ConnectionException, true},
		{"ConnectionDoesNotExist", ConnectionDoesNotExist, true},
		{"ConnectionFailure", ConnectionFailure, true},
		{"SQLClientUnableToEstablishSQLConnection", SQLClientUnableToEstablishSQLConnection, true},
		{"SQLServerRejectedEstablishmentOfSQLConnection", SQLServerRejectedEstablishmentOfSQLConnection, true},
		{"TransactionResolutionUnknown", TransactionResolutionUnknown, true},
		{"ProtocolViolation", ProtocolViolation, true},
		{"UnknownCode", "99999", false},
		{"EmptyCode", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsConnectionException(tt.code)
			assert.Equal(t, tt.expected, result, "IsConnectionException(%s) = %v, expected %v", tt.code, result, tt.expected)
		})
	}
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

// Mock pool for testing
type mockPool struct{}

func (m *mockPool) Close() {}
