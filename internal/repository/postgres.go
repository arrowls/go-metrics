package repository

import (
	"context"
	"fmt"

	"github.com/arrowls/go-metrics/internal/memstorage"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresRepository struct {
	db         *pgxpool.Pool
	metricRepo Metric
}

func NewPostgresRepository(db *pgxpool.Pool, metricRepo Metric) Metric {
	return &PostgresRepository{
		db,
		metricRepo,
	}
}
func (m *PostgresRepository) AddGaugeValue(ctx context.Context, name string, value float64) {
	m.metricRepo.AddGaugeValue(ctx, name, value)
}

func (m *PostgresRepository) AddCounterValue(ctx context.Context, name string, value int64) {
	m.metricRepo.AddCounterValue(ctx, name, value)
}

func (m *PostgresRepository) GetAll(ctx context.Context) memstorage.MemStorage {
	return m.metricRepo.GetAll(ctx)
}

func (m *PostgresRepository) GetGaugeItem(ctx context.Context, name string) (float64, error) {
	return m.metricRepo.GetGaugeItem(ctx, name)
}

func (m *PostgresRepository) GetCounterItem(ctx context.Context, name string) (int64, error) {
	return m.metricRepo.GetCounterItem(ctx, name)
}

func (m *PostgresRepository) CheckConnection(ctx context.Context) bool {
	if err := m.db.Ping(ctx); err != nil {
		fmt.Println(err)
		return false
	}
	return true
}
