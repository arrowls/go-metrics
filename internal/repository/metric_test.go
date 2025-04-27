package repository

import (
	"testing"

	"github.com/arrowls/go-metrics/internal/memstorage"
	"github.com/stretchr/testify/assert"
)

func TestMetricRepository_AddGaugeValue(t *testing.T) {
	storage := memstorage.GetInstance()
	repo := NewMetricRepository(storage)

	repo.AddGaugeValue("test name", 1.23)

	assert.Equal(t, 1.23, storage.Gauge["test name"])

	repo.AddGaugeValue("test name", 4.56)

	assert.Equal(t, 4.56, storage.Gauge["test name"])
}

func TestMetricRepository_AddCounterValue(t *testing.T) {
	storage := memstorage.GetInstance()
	repo := NewMetricRepository(storage)

	repo.AddCounterValue("test name", 123)

	assert.Equal(t, int64(123), storage.Counter["test name"])

	repo.AddCounterValue("test name", 456)

	assert.Equal(t, int64(123+456), storage.Counter["test name"])
}

func TestMetricRepository_GetAll(t *testing.T) {
	storage := memstorage.GetInstance()
	repo := NewMetricRepository(storage)

	value := repo.GetAll()

	assert.Equal(t, *storage, value)
}

func TestMetricRepository_GetGaugeItem(t *testing.T) {
	storage := memstorage.GetInstance()
	repo := NewMetricRepository(storage)

	t.Run("success", func(t *testing.T) {
		storage.Gauge["test name"] = 1.23

		value, err := repo.GetGaugeItem("test name")

		assert.NoError(t, err)
		assert.Equal(t, 1.23, value)
	})

	t.Run("fail", func(t *testing.T) {
		storage.Gauge["test name"] = 4.56
		value, err := repo.GetGaugeItem("undefined name")

		assert.Error(t, err)
		assert.Equal(t, float64(0), value)
	})
}

func TestMetricRepository_GetCounterItem(t *testing.T) {
	storage := memstorage.GetInstance()
	repo := NewMetricRepository(storage)

	t.Run("success", func(t *testing.T) {
		storage.Counter["test name"] = 123

		value, err := repo.GetCounterItem("test name")

		assert.NoError(t, err)
		assert.Equal(t, int64(123), value)
	})

	t.Run("fail", func(t *testing.T) {
		storage.Counter["test name"] = 456
		value, err := repo.GetCounterItem("undefined name")

		assert.Error(t, err)
		assert.Equal(t, int64(0), value)
	})
}
