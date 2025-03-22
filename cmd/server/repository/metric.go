package repository

import (
	"github.com/arrowls/go-metrics/internal/memstorage"
)

type MetricRepository struct{}

func NewMetricRepository() *MetricRepository {
	return &MetricRepository{}
}

func (m *MetricRepository) AddGaugeValue(name string, value float64) {
	storage := *memstorage.GetInstance()

	if storage.Gauge == nil {
		storage.Gauge = make(map[string]float64)
	}

	storage.Gauge[name] = value
}

func (m *MetricRepository) AddCounterValue(name string, value int64) {
	storage := *memstorage.GetInstance()

	if storage.Counter == nil {
		storage.Counter = make(map[string]int64)
	}

	storage.Counter[name] += value
}
