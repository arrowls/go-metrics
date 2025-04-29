package repository

import (
	"context"

	"github.com/arrowls/go-metrics/internal/memstorage"
)

type Metric interface {
	AddGaugeValue(ctx context.Context, key string, value float64) error
	AddCounterValue(ctx context.Context, key string, value int64) error
	GetAll(ctx context.Context) (memstorage.MemStorage, error)
	GetCounterItem(ctx context.Context, name string) (int64, error)
	GetGaugeItem(ctx context.Context, name string) (float64, error)
	CheckConnection(ctx context.Context) bool
}

type Repository struct {
	Metric Metric
}

func NewRepository(storage *memstorage.MemStorage) *Repository {
	return &Repository{
		Metric: NewMetricRepository(storage),
	}
}
