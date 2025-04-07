package service

import (
	"fmt"
	"testing"

	"github.com/arrowls/go-metrics/internal/memstorage"
	"github.com/arrowls/go-metrics/internal/repository"
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

func (m *MockMetric) GetAll() memstorage.MemStorage {
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
func (m *MockMetric) GetCounterItem(name string) (int64, error) {
	if name != "" {
		return int64(123), nil
	}
	return 0, fmt.Errorf("error")
}
func (m *MockMetric) GetGaugeItem(name string) (float64, error) {
	if name != "" {
		return float64(123), nil
	}
	return 0, fmt.Errorf("error")
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

	t.Run("it should return list", func(t *testing.T) {
		metricMap := service.GetList()

		assert.Equal(t, *metricMap, map[string]interface{}{
			"key1.1": float64(1.1),
			"key2.2": float64(2.2),
			"key1":   int64(1),
			"key2":   int64(2),
		})
	})

	t.Run("it should return item", func(t *testing.T) {
		value, err := service.GetItem("gauge", "key")

		assert.Nil(t, err)
		assert.Equal(t, "123", value)

		value, err = service.GetItem("counter", "key")

		assert.Nil(t, err)
		assert.Equal(t, "123", value)
	})

	t.Run("it should handle return item error", func(t *testing.T) {
		value, err := service.GetItem("", "key")

		assert.NotNil(t, err)
		assert.Equal(t, "", value)

		value, err = service.GetItem("gauge", "")
		assert.NotNil(t, err)
		assert.Equal(t, "", value)

		value, err = service.GetItem("counter", "")
		assert.NotNil(t, err)
		assert.Equal(t, "", value)
	})
}
