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
