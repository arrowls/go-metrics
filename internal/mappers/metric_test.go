package mappers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/arrowls/go-metrics/internal/dto"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)

func TestHTTPToCreateMetric(t *testing.T) {
	tests := []struct {
		name        string
		metricType  string
		metricName  string
		metricValue string
		expectError bool
	}{
		{
			name:        "valid gauge metric",
			metricType:  "gauge",
			metricName:  "testMetric",
			metricValue: "123.45",
			expectError: false,
		},
		{
			name:        "valid counter metric",
			metricType:  "counter",
			metricName:  "testCounter",
			metricValue: "100",
			expectError: false,
		},
		{
			name:        "missing value",
			metricType:  "gauge",
			metricName:  "testMetric",
			metricValue: "",
			expectError: true,
		},
		{
			name:        "missing name",
			metricType:  "gauge",
			metricName:  "",
			metricValue: "123.45",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("POST", "/update/"+tt.metricType+"/"+tt.metricName+"/"+tt.metricValue, nil)
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("type", tt.metricType)
			rctx.URLParams.Add("name", tt.metricName)
			rctx.URLParams.Add("value", tt.metricValue)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

			result, err := HTTPToCreateMetric(req)

			if tt.expectError {
				assert.Error(t, err, "Expected error but got none")
			} else {
				assert.NoError(t, err, "Unexpected error")
				assert.Equal(t, tt.metricType, result.Type, "Expected type %s, got %s", tt.metricType, result.Type)
				assert.Equal(t, tt.metricName, result.Name, "Expected name %s, got %s", tt.metricName, result.Name)
				assert.Equal(t, tt.metricValue, result.Value, "Expected value %s, got %s", tt.metricValue, result.Value)
			}
		})
	}
}

func TestHTTPToGetMetric(t *testing.T) {
	tests := []struct {
		name        string
		metricType  string
		metricName  string
		expectError bool
	}{
		{
			name:        "valid gauge metric",
			metricType:  "gauge",
			metricName:  "testMetric",
			expectError: false,
		},
		{
			name:        "valid counter metric",
			metricType:  "counter",
			metricName:  "testCounter",
			expectError: false,
		},
		{
			name:        "missing name",
			metricType:  "gauge",
			metricName:  "",
			expectError: true,
		},
		{
			name:        "invalid metric type",
			metricType:  "invalid",
			metricName:  "testMetric",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/value/"+tt.metricType+"/"+tt.metricName, nil)
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("type", tt.metricType)
			rctx.URLParams.Add("name", tt.metricName)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

			result, err := HTTPToGetMetric(req)

			if tt.expectError {
				assert.Error(t, err, "Expected error but got none")
			} else {
				assert.NoError(t, err, "Unexpected error")
				assert.Equal(t, tt.metricType, result.Type, "Expected type %s, got %s", tt.metricType, result.Type)
				assert.Equal(t, tt.metricName, result.Name, "Expected name %s, got %s", tt.metricName, result.Name)
			}
		})
	}
}

func TestHTTPWithBodyToCreateMetric(t *testing.T) {
	tests := []struct {
		name        string
		requestBody dto.Metrics
		expectError bool
	}{
		{
			name: "valid gauge metric",
			requestBody: dto.Metrics{
				ID:    "testMetric",
				MType: "gauge",
				Value: func() *float64 { v := 123.45; return &v }(),
			},
			expectError: false,
		},
		{
			name: "valid counter metric",
			requestBody: dto.Metrics{
				ID:    "testCounter",
				MType: "counter",
				Delta: func() *int64 { v := int64(100); return &v }(),
			},
			expectError: false,
		},
		{
			name: "missing value for gauge",
			requestBody: dto.Metrics{
				ID:    "testMetric",
				MType: "gauge",
			},
			expectError: true,
		},
		{
			name: "missing delta for counter",
			requestBody: dto.Metrics{
				ID:    "testCounter",
				MType: "counter",
			},
			expectError: true,
		},
		{
			name: "invalid metric type",
			requestBody: dto.Metrics{
				ID:    "testMetric",
				MType: "invalid",
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bodyBytes, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest("POST", "/update/", bytes.NewReader(bodyBytes))

			result, err := HTTPWithBodyToCreateMetric(req)

			if tt.expectError {
				assert.Error(t, err, "Expected error but got none")
			} else {
				assert.NoError(t, err, "Unexpected error")
				assert.Equal(t, tt.requestBody.MType, result.Type, "Expected type %s, got %s", tt.requestBody.MType, result.Type)
				assert.Equal(t, tt.requestBody.ID, result.Name, "Expected name %s, got %s", tt.requestBody.ID, result.Name)
			}
		})
	}
}

func TestCreateDTOToHTTPRes(t *testing.T) {
	tests := []struct {
		name        string
		createDto   *dto.CreateMetric
		value       string
		expectError bool
	}{
		{
			name: "valid gauge response",
			createDto: &dto.CreateMetric{
				Type:  "gauge",
				Name:  "testMetric",
				Value: "123.45",
			},
			value:       "123.45",
			expectError: false,
		},
		{
			name: "valid counter response",
			createDto: &dto.CreateMetric{
				Type:  "counter",
				Name:  "testCounter",
				Value: "100",
			},
			value:       "100",
			expectError: false,
		},
		{
			name: "invalid gauge value",
			createDto: &dto.CreateMetric{
				Type:  "gauge",
				Name:  "testMetric",
				Value: "invalid",
			},
			value:       "invalid",
			expectError: true,
		},
		{
			name: "invalid counter value",
			createDto: &dto.CreateMetric{
				Type:  "counter",
				Name:  "testCounter",
				Value: "invalid",
			},
			value:       "invalid",
			expectError: true,
		},
		{
			name: "invalid metric type",
			createDto: &dto.CreateMetric{
				Type:  "invalid",
				Name:  "testMetric",
				Value: "123.45",
			},
			value:       "123.45",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := CreateDTOToHTTPRes(tt.createDto, tt.value)

			if tt.expectError {
				assert.Error(t, err, "Expected error but got none")
			} else {
				assert.NoError(t, err, "Unexpected error")
				assert.NotNil(t, result, "Expected response but got nil")
			}
		})
	}
}

func TestMetricToDTO(t *testing.T) {
	tests := []struct {
		name        string
		metricName  string
		value       interface{}
		expectError bool
	}{
		{
			name:        "float64 value",
			metricName:  "testMetric",
			value:       float64(123.45),
			expectError: false,
		},
		{
			name:        "int64 value",
			metricName:  "testCounter",
			value:       int64(100),
			expectError: false,
		},
		{
			name:        "invalid value type",
			metricName:  "testMetric",
			value:       "string",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := MetricToDTO(tt.metricName, tt.value)

			if tt.expectError {
				assert.Error(t, err, "Expected error but got none")
			} else {
				assert.NoError(t, err, "Unexpected error")
				assert.Equal(t, tt.metricName, result.ID, "Expected ID %s, got %s", tt.metricName, result.ID)
			}
		})
	}
}
