package dto

import (
	"testing"
)

func TestMetrics_Validate(t *testing.T) {
	tests := []struct {
		name        string
		metric      Metrics
		expectError bool
	}{
		{
			name: "valid gauge metric",
			metric: Metrics{
				ID:    "testMetric",
				MType: "gauge",
				Value: func() *float64 { v := 123.45; return &v }(),
			},
			expectError: false,
		},
		{
			name: "valid counter metric",
			metric: Metrics{
				ID:    "testCounter",
				MType: "counter",
				Delta: func() *int64 { v := int64(100); return &v }(),
			},
			expectError: false,
		},
		{
			name: "missing ID",
			metric: Metrics{
				ID:    "",
				MType: "gauge",
				Value: func() *float64 { v := 123.45; return &v }(),
			},
			expectError: true,
		},
		{
			name: "invalid metric type",
			metric: Metrics{
				ID:    "testMetric",
				MType: "invalid",
				Value: func() *float64 { v := 123.45; return &v }(),
			},
			expectError: true,
		},
		{
			name: "empty metric type",
			metric: Metrics{
				ID:    "testMetric",
				MType: "",
				Value: func() *float64 { v := 123.45; return &v }(),
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.metric.Validate()

			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
			}
		})
	}
}

func TestCreateMetric(t *testing.T) {
	metric := CreateMetric{
		Type:  "gauge",
		Name:  "testMetric",
		Value: "123.45",
	}

	if metric.Type != "gauge" {
		t.Errorf("Expected type 'gauge', got '%s'", metric.Type)
	}

	if metric.Name != "testMetric" {
		t.Errorf("Expected name 'testMetric', got '%s'", metric.Name)
	}

	if metric.Value != "123.45" {
		t.Errorf("Expected value '123.45', got '%s'", metric.Value)
	}
}

func TestGetMetric(t *testing.T) {
	metric := GetMetric{
		Type: "counter",
		Name: "testCounter",
	}

	if metric.Type != "counter" {
		t.Errorf("Expected type 'counter', got '%s'", metric.Type)
	}

	if metric.Name != "testCounter" {
		t.Errorf("Expected name 'testCounter', got '%s'", metric.Name)
	}
}

func TestMetrics_WithGaugeValue(t *testing.T) {
	value := 123.45
	metric := Metrics{
		ID:    "testMetric",
		MType: "gauge",
		Value: &value,
	}

	if metric.ID != "testMetric" {
		t.Errorf("Expected ID 'testMetric', got '%s'", metric.ID)
	}

	if metric.MType != "gauge" {
		t.Errorf("Expected MType 'gauge', got '%s'", metric.MType)
	}

	if metric.Value == nil {
		t.Error("Expected Value to be set")
	}

	if *metric.Value != 123.45 {
		t.Errorf("Expected Value 123.45, got %f", *metric.Value)
	}

	if metric.Delta != nil {
		t.Error("Expected Delta to be nil for gauge metric")
	}
}

func TestMetrics_WithCounterValue(t *testing.T) {
	delta := int64(100)
	metric := Metrics{
		ID:    "testCounter",
		MType: "counter",
		Delta: &delta,
	}

	if metric.ID != "testCounter" {
		t.Errorf("Expected ID 'testCounter', got '%s'", metric.ID)
	}

	if metric.MType != "counter" {
		t.Errorf("Expected MType 'counter', got '%s'", metric.MType)
	}

	if metric.Delta == nil {
		t.Error("Expected Delta to be set")
	}

	if *metric.Delta != 100 {
		t.Errorf("Expected Delta 100, got %d", *metric.Delta)
	}

	if metric.Value != nil {
		t.Error("Expected Value to be nil for counter metric")
	}
}

func TestMetrics_Validate_EdgeCases(t *testing.T) {
	tests := []struct {
		name        string
		metric      Metrics
		expectError bool
	}{
		{
			name: "gauge with nil value",
			metric: Metrics{
				ID:    "testMetric",
				MType: "gauge",
				Value: nil,
			},
			expectError: false, // Validate only checks ID and MType
		},
		{
			name: "counter with nil delta",
			metric: Metrics{
				ID:    "testCounter",
				MType: "counter",
				Delta: nil,
			},
			expectError: false, // Validate only checks ID and MType
		},
		{
			name: "gauge with both value and delta",
			metric: Metrics{
				ID:    "testMetric",
				MType: "gauge",
				Value: func() *float64 { v := 123.45; return &v }(),
				Delta: func() *int64 { v := int64(100); return &v }(),
			},
			expectError: false, // Validate only checks ID and MType
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.metric.Validate()

			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
			}
		})
	}
}
