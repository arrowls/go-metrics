package repository

import (
	"github.com/arrowls/go-metrics/internal/memstorage"
)

type Metric interface {
	AddGaugeValue(key string, value float64)
	AddCounterValue(key string, value int64)
	GetAll() memstorage.MemStorage
	GetCounterItem(name string) (int64, error)
	GetGaugeItem(name string) (float64, error)
}

type Repository struct {
	Metric Metric
}

func NewRepository(storage *memstorage.MemStorage) *Repository {
	return &Repository{
		Metric: NewMetricRepository(storage),
	}
}
