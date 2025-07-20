package repository

import (
	"context"
	"testing"

	"github.com/arrowls/go-metrics/internal/dto"
	"github.com/arrowls/go-metrics/internal/memstorage"
	"github.com/stretchr/testify/assert"
)

func TestMetricRepository_AddGaugeValue(t *testing.T) {
	storage := memstorage.GetInstance()
	repo := NewMetricRepository(storage)

	ctx := context.Background()
	repo.AddGaugeValue(ctx, "test name", 1.23)
	assert.Equal(t, 1.23, storage.Gauge["test name"])

	repo.AddGaugeValue(ctx, "test name", 4.56)

	assert.Equal(t, 4.56, storage.Gauge["test name"])
}

func TestMetricRepository_AddCounterValue(t *testing.T) {
	storage := memstorage.GetInstance()
	repo := NewMetricRepository(storage)

	ctx := context.Background()
	repo.AddCounterValue(ctx, "test name", 123)
	assert.Equal(t, int64(123), storage.Counter["test name"])
}

func TestMetricRepository_GetAll(t *testing.T) {
	storage := memstorage.GetInstance()
	repo := NewMetricRepository(storage)

	ctx := context.Background()
	value, _ := repo.GetAll(ctx)
	assert.Equal(t, *storage, value)
}

func TestMetricRepository_GetGaugeItem(t *testing.T) {
	storage := memstorage.GetInstance()
	repo := NewMetricRepository(storage)

	t.Run("success", func(t *testing.T) {
		ctx := context.Background()
		storage.Gauge["test name"] = 1.23
		value, err := repo.GetGaugeItem(ctx, "test name")
		assert.NoError(t, err)
		assert.Equal(t, 1.23, value)
	})

	t.Run("fail", func(t *testing.T) {
		ctx := context.Background()
		storage.Gauge["test name"] = 4.56
		value, err := repo.GetGaugeItem(ctx, "undefined name")
		assert.Error(t, err)
		assert.Equal(t, float64(0), value)
	})
}

func TestMetricRepository_GetCounterItem(t *testing.T) {
	storage := memstorage.GetInstance()
	repo := NewMetricRepository(storage)

	t.Run("success", func(t *testing.T) {
		ctx := context.Background()
		storage.Counter["test name"] = 123
		value, err := repo.GetCounterItem(ctx, "test name")
		assert.NoError(t, err)
		assert.Equal(t, int64(123), value)
	})

	t.Run("fail", func(t *testing.T) {
		ctx := context.Background()
		storage.Counter["test name"] = 456
		value, err := repo.GetCounterItem(ctx, "undefined name")
		assert.Error(t, err)
		assert.Equal(t, int64(0), value)
	})
}

func TestMetricRepository_CheckConnection(t *testing.T) {
	storage := memstorage.GetInstance()
	repo := NewMetricRepository(storage)

	ctx := context.Background()
	result := repo.CheckConnection(ctx)
	assert.True(t, result)
}

func TestMetricRepository_CreateBatch(t *testing.T) {
	storage := memstorage.GetInstance()
	repo := NewMetricRepository(storage)

	ctx := context.Background()

	t.Run("success with gauge and counter", func(t *testing.T) {
		batch := []dto.CreateMetric{
			{Name: "test_gauge", Type: "gauge", Value: "1.23"},
			{Name: "test_counter", Type: "counter", Value: "123"},
		}

		err := repo.CreateBatch(ctx, batch)
		assert.NoError(t, err)

		gaugeValue, err := repo.GetGaugeItem(ctx, "test_gauge")
		assert.NoError(t, err)
		assert.Equal(t, 1.23, gaugeValue)

		counterValue, err := repo.GetCounterItem(ctx, "test_counter")
		assert.NoError(t, err)
		assert.Equal(t, int64(123), counterValue)
	})

	t.Run("invalid gauge value", func(t *testing.T) {
		batch := []dto.CreateMetric{
			{Name: "test_gauge", Type: "gauge", Value: "invalid"},
		}

		err := repo.CreateBatch(ctx, batch)
		assert.Error(t, err)
	})

	t.Run("invalid counter value", func(t *testing.T) {
		batch := []dto.CreateMetric{
			{Name: "test_counter", Type: "counter", Value: "invalid"},
		}

		err := repo.CreateBatch(ctx, batch)
		assert.Error(t, err)
	})

	t.Run("unknown metric type", func(t *testing.T) {
		batch := []dto.CreateMetric{
			{Name: "test_unknown", Type: "unknown", Value: "123"},
		}

		err := repo.CreateBatch(ctx, batch)
		assert.NoError(t, err)
	})
}
