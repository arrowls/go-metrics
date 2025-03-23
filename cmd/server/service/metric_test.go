package service

import (
	"testing"

	"github.com/arrowls/go-metrics/cmd/server/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockMetric struct {
	mock.Mock
}

func (m *MockMetric) AddGaugeValue(key string, value float64) {
	m.Mock.Called(key, value)
}

func (m *MockMetric) AddCounterValue(key string, value int64) {
	m.Mock.Called(key, value)
}

func TestMetricService_CreateByType(t *testing.T) {
	metric := &MockMetric{}
	repo := &repository.Repository{Metric: metric}

	service := NewMetricService(repo)

	t.Run("it should handle invalid type", func(t *testing.T) {
		err := service.CreateByType("non-existent type", "", "")

		assert.NotNil(t, err)
	})

	t.Run("it should create a metric with Gauge type", func(t *testing.T) {
		metric.Mock.On("AddGaugeValue", "MetricName", 1.23)
		err := service.CreateByType("gauge", "MetricName", "1.23")

		assert.Nil(t, err)

		metric.Mock.AssertCalled(t, "AddGaugeValue", "MetricName", 1.23)
	})

	t.Run("it should create a metric with Counter type", func(t *testing.T) {
		metric.Mock.On("AddCounterValue", "MetricName", int64(123))
		err := service.CreateByType("counter", "MetricName", "123")

		assert.Nil(t, err)

		metric.Mock.AssertCalled(t, "AddCounterValue", "MetricName", int64(123))
	})

	t.Run("it should handle invalid value", func(t *testing.T) {
		err := service.CreateByType("counter", "123", "definitely not a number")

		assert.NotNil(t, err)

		err2 := service.CreateByType("gauge", "123", "definitely not a number")

		assert.NotNil(t, err2)
	})
}
