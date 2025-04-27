package service

import (
	"context"
	"fmt"
	"testing"

	"github.com/arrowls/go-metrics/internal/dto"
	"github.com/arrowls/go-metrics/internal/memstorage"
	"github.com/arrowls/go-metrics/internal/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockMetric struct {
	mock.Mock
}

func (m *MockMetric) AddGaugeValue(_ context.Context, key string, value float64) {
	m.Mock.Called(key, value)
}

func (m *MockMetric) AddCounterValue(_ context.Context, key string, value int64) {
	m.Mock.Called(key, value)
}

func (m *MockMetric) GetAll(_ context.Context) memstorage.MemStorage {
	inst := *memstorage.GetInstance()

	inst.Gauge = map[string]float64{
		"key1.1": 1.1,
		"key2.2": 2.2,
	}

	inst.Counter = map[string]int64{
		"key1": 1,
		"key2": 2,
	}

	return inst
}
func (m *MockMetric) GetCounterItem(_ context.Context, name string) (int64, error) {
	if name != "" {
		return int64(123), nil
	}
	return 0, fmt.Errorf("error")
}
func (m *MockMetric) GetGaugeItem(_ context.Context, name string) (float64, error) {
	if name != "" {
		return float64(123), nil
	}
	return 0, fmt.Errorf("error")
}

func (m *MockMetric) CheckConnection(_ context.Context) bool {
	return true
}

func TestMetricService_CreateByType(t *testing.T) {
	ctx := context.Background()

	metric := &MockMetric{}
	repo := &repository.Repository{Metric: metric}

	service := NewMetricService(repo)

	t.Run("it should handle invalid type", func(t *testing.T) {
		err := service.Create(ctx, &dto.CreateMetric{
			Type:  "non-existent",
			Name:  "",
			Value: "",
		})

		assert.NotNil(t, err)
	})

	t.Run("it should create a metric with Gauge type", func(t *testing.T) {
		metric.Mock.On("AddGaugeValue", "MetricName", 1.23)
		err := service.Create(ctx, &dto.CreateMetric{
			Type:  "gauge",
			Name:  "MetricName",
			Value: "1.23",
		})

		assert.Nil(t, err)

		metric.Mock.AssertCalled(t, "AddGaugeValue", "MetricName", 1.23)
	})

	t.Run("it should create a metric with Counter type", func(t *testing.T) {
		metric.Mock.On("AddCounterValue", "MetricName", int64(123))
		err := service.Create(ctx, &dto.CreateMetric{
			Type:  "counter",
			Name:  "MetricName",
			Value: "123",
		})

		assert.Nil(t, err)

		metric.Mock.AssertCalled(t, "AddCounterValue", "MetricName", int64(123))
	})

	t.Run("it should handle invalid value", func(t *testing.T) {
		err := service.Create(ctx, &dto.CreateMetric{
			Type:  "counter",
			Name:  "123",
			Value: "definitely not a number",
		})

		assert.NotNil(t, err)

		err2 := service.Create(ctx, &dto.CreateMetric{
			Type:  "gauge",
			Name:  "123",
			Value: "definitely not a number",
		})

		assert.NotNil(t, err2)
	})

	t.Run("it should return list", func(t *testing.T) {
		metricMap := service.GetList(ctx)

		assert.Equal(t, *metricMap, map[string]interface{}{
			"key1.1": 1.1,
			"key2.2": 2.2,
			"key1":   int64(1),
			"key2":   int64(2),
		})
	})

	t.Run("it should return item", func(t *testing.T) {
		value, err := service.GetItem(ctx, &dto.GetMetric{
			Type: "gauge",
			Name: "key",
		})

		assert.Nil(t, err)
		assert.Equal(t, "123", value)

		value, err = service.GetItem(ctx, &dto.GetMetric{
			Type: "counter",
			Name: "key",
		})

		assert.Nil(t, err)
		assert.Equal(t, "123", value)
	})

	t.Run("it should handle return item error", func(t *testing.T) {
		value, err := service.GetItem(ctx, &dto.GetMetric{
			Type: "",
			Name: "key",
		})

		assert.NotNil(t, err)
		assert.Equal(t, "", value)

		value, err = service.GetItem(ctx, &dto.GetMetric{
			Type: "gauge",
			Name: "",
		})
		assert.NotNil(t, err)
		assert.Equal(t, "", value)

		value, err = service.GetItem(ctx, &dto.GetMetric{
			Type: "counter",
			Name: "",
		})
		assert.NotNil(t, err)
		assert.Equal(t, "", value)
	})
}
