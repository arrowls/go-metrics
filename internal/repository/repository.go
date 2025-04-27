package repository

import (
	"context"

	"github.com/arrowls/go-metrics/internal/memstorage"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Metric interface {
	AddGaugeValue(ctx context.Context, key string, value float64)
	AddCounterValue(ctx context.Context, key string, value int64)
	GetAll(ctx context.Context) memstorage.MemStorage
	GetCounterItem(ctx context.Context, name string) (int64, error)
	GetGaugeItem(ctx context.Context, name string) (float64, error)
	CheckConnection(ctx context.Context) bool
}

type Repository struct {
	Metric Metric
}

func NewRepository(storage *memstorage.MemStorage, db *pgxpool.Pool) *Repository {
	metricRepo := NewMetricRepository(storage)
	return &Repository{
		Metric: NewPostgresRepository(db, metricRepo),
	}
}
