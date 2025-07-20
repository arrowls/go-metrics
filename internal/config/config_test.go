package config

import (
	"os"
	"testing"

	"github.com/caarlos0/env/v7"
	"github.com/stretchr/testify/assert"
)

func TestAgentConfig_EnvironmentParsing(t *testing.T) {
	// Save original environment variables
	originalReportInterval := os.Getenv("REPORT_INTERVAL")
	originalPollInterval := os.Getenv("POLL_INTERVAL")
	originalAddress := os.Getenv("ADDRESS")
	originalKey := os.Getenv("KEY")
	originalRateLimit := os.Getenv("RATE_LIMIT")

	// Restore after test
	defer func() {
		if originalReportInterval != "" {
			os.Setenv("REPORT_INTERVAL", originalReportInterval)
		} else {
			os.Unsetenv("REPORT_INTERVAL")
		}
		if originalPollInterval != "" {
			os.Setenv("POLL_INTERVAL", originalPollInterval)
		} else {
			os.Unsetenv("POLL_INTERVAL")
		}
		if originalAddress != "" {
			os.Setenv("ADDRESS", originalAddress)
		} else {
			os.Unsetenv("ADDRESS")
		}
		if originalKey != "" {
			os.Setenv("KEY", originalKey)
		} else {
			os.Unsetenv("KEY")
		}
		if originalRateLimit != "" {
			os.Setenv("RATE_LIMIT", originalRateLimit)
		} else {
			os.Unsetenv("RATE_LIMIT")
		}
	}()

	// Test with environment variables
	os.Setenv("REPORT_INTERVAL", "15")
	os.Setenv("POLL_INTERVAL", "5")
	os.Setenv("ADDRESS", "localhost:9090")
	os.Setenv("KEY", "test-key")
	os.Setenv("RATE_LIMIT", "100")

	// Test environment parsing
	config := AgentConfig{}
	err := env.Parse(&config)
	assert.NoError(t, err, "Failed to parse env")

	assert.Equal(t, 15, config.ReportInterval, "Expected ReportInterval to be 15")
	assert.Equal(t, 5, config.PollInterval, "Expected PollInterval to be 5")
	assert.Equal(t, "localhost:9090", config.ServerEndpoint, "Expected ServerEndpoint to be localhost:9090")
	assert.Equal(t, "test-key", config.Key, "Expected Key to be test-key")
	assert.Equal(t, 100, config.RateLimit, "Expected RateLimit to be 100")
}

func TestAgentConfig_Defaults(t *testing.T) {
	// Clear environment variables
	os.Unsetenv("REPORT_INTERVAL")
	os.Unsetenv("POLL_INTERVAL")
	os.Unsetenv("ADDRESS")
	os.Unsetenv("KEY")
	os.Unsetenv("RATE_LIMIT")

	// Test default values
	config := AgentConfig{}
	err := env.Parse(&config)
	assert.NoError(t, err, "Failed to parse env")

	// Check defaults (should be zero values when no env vars are set)
	assert.Equal(t, 0, config.ReportInterval, "Expected ReportInterval to be 0")
	assert.Equal(t, 0, config.PollInterval, "Expected PollInterval to be 0")
	assert.Equal(t, "", config.ServerEndpoint, "Expected ServerEndpoint to be empty")
}

func TestServerConfig_EnvironmentParsing(t *testing.T) {
	// Save original environment variables
	originalAddress := os.Getenv("ADDRESS")
	originalStoreInterval := os.Getenv("STORE_INTERVAL")
	originalFilePath := os.Getenv("FILE_STORAGE_PATH")
	originalRestore := os.Getenv("RESTORE")
	originalDatabaseDSN := os.Getenv("DATABASE_DSN")
	originalKey := os.Getenv("KEY")

	// Restore after test
	defer func() {
		if originalAddress != "" {
			os.Setenv("ADDRESS", originalAddress)
		} else {
			os.Unsetenv("ADDRESS")
		}
		if originalStoreInterval != "" {
			os.Setenv("STORE_INTERVAL", originalStoreInterval)
		} else {
			os.Unsetenv("STORE_INTERVAL")
		}
		if originalFilePath != "" {
			os.Setenv("FILE_STORAGE_PATH", originalFilePath)
		} else {
			os.Unsetenv("FILE_STORAGE_PATH")
		}
		if originalRestore != "" {
			os.Setenv("RESTORE", originalRestore)
		} else {
			os.Unsetenv("RESTORE")
		}
		if originalDatabaseDSN != "" {
			os.Setenv("DATABASE_DSN", originalDatabaseDSN)
		} else {
			os.Unsetenv("DATABASE_DSN")
		}
		if originalKey != "" {
			os.Setenv("KEY", originalKey)
		} else {
			os.Unsetenv("KEY")
		}
	}()

	// Test with environment variables
	os.Setenv("ADDRESS", "localhost:8080")
	os.Setenv("STORE_INTERVAL", "600")
	os.Setenv("FILE_STORAGE_PATH", "/tmp/metrics.json")
	os.Setenv("RESTORE", "true")
	os.Setenv("DATABASE_DSN", "postgres://user:pass@localhost/db")
	os.Setenv("KEY", "server-key")

	// Test environment parsing
	config := ServerConfig{}
	err := env.Parse(&config)
	assert.NoError(t, err, "Failed to parse env")

	assert.Equal(t, "localhost:8080", config.ServerEndpoint, "Expected ServerEndpoint to be localhost:8080")
	assert.Equal(t, 600, config.StoreInterval, "Expected StoreInterval to be 600")
	assert.Equal(t, "/tmp/metrics.json", config.StorageFilePath, "Expected StorageFilePath to be /tmp/metrics.json")
	assert.True(t, config.Restore, "Expected Restore to be true")
	assert.Equal(t, "postgres://user:pass@localhost/db", config.DatabaseDSN, "Expected DatabaseDSN to be postgres://user:pass@localhost/db")
	assert.Equal(t, "server-key", config.Key, "Expected Key to be server-key")
}

func TestServerConfig_Defaults(t *testing.T) {
	// Clear environment variables
	os.Unsetenv("ADDRESS")
	os.Unsetenv("STORE_INTERVAL")
	os.Unsetenv("FILE_STORAGE_PATH")
	os.Unsetenv("RESTORE")
	os.Unsetenv("DATABASE_DSN")
	os.Unsetenv("KEY")

	// Test default values
	config := ServerConfig{}
	err := env.Parse(&config)
	assert.NoError(t, err, "Failed to parse env")

	// Check defaults (should be zero values when no env vars are set)
	assert.Equal(t, "", config.ServerEndpoint, "Expected ServerEndpoint to be empty")
	assert.Equal(t, 0, config.StoreInterval, "Expected StoreInterval to be 0")
	assert.Equal(t, "", config.StorageFilePath, "Expected StorageFilePath to be empty")
	assert.False(t, config.Restore, "Expected Restore to be false")
}

func TestProvideServerConfig(t *testing.T) {
	// Create a mock container
	mockContainer := &mockContainer{
		items: make(map[string]interface{}),
	}

	config := ProvideServerConfig(mockContainer)

	// Verify config was created
	assert.NotEmpty(t, config.ServerEndpoint, "Expected ServerEndpoint to be set")

	// Verify config was added to container
	assert.NotNil(t, mockContainer.items[diKey], "Expected config to be added to container")
}

func TestConstants(t *testing.T) {
	assert.Equal(t, "gauge", GaugeType, "Expected GaugeType to be 'gauge'")
	assert.Equal(t, "counter", CounterType, "Expected CounterType to be 'counter'")
	assert.Equal(t, "HashSHA256", HashHeaderName, "Expected HashHeaderName to be 'HashSHA256'")
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
