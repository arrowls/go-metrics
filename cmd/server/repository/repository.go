package repository

import (
	"github.com/arrowls/go-metrics/internal/memstorage"
)

type Metric interface {
	AddGaugeValue(key string, value float64)
	AddCounterValue(key string, value int64)
}

type Repository struct {
	Metric Metric
}

func NewRepository(storage *memstorage.MemStorage) *Repository {
	return &Repository{
		Metric: NewMetricRepository(storage),
	}
}
