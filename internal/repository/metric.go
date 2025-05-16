package repository

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/arrowls/go-metrics/internal/apperrors"
	"github.com/arrowls/go-metrics/internal/config"
	"github.com/arrowls/go-metrics/internal/dto"
	"github.com/arrowls/go-metrics/internal/memstorage"
)

type MetricRepository struct {
	storage *memstorage.MemStorage
}

func NewMetricRepository(storage *memstorage.MemStorage) Metric {
	return &MetricRepository{
		storage: storage,
	}
}

func (m *MetricRepository) AddGaugeValue(_ context.Context, name string, value float64) error {
	m.storage.Lock()
	defer m.storage.Unlock()

	m.storage.Gauge[name] = value
	return nil
}

func (m *MetricRepository) AddCounterValue(_ context.Context, name string, value int64) error {
	m.storage.Lock()
	defer m.storage.Unlock()

	m.storage.Counter[name] += value
	return nil
}

func (m *MetricRepository) GetAll(_ context.Context) (memstorage.MemStorage, error) {
	return *m.storage, nil
}

func (m *MetricRepository) GetGaugeItem(_ context.Context, name string) (float64, error) {
	m.storage.Lock()
	defer m.storage.Unlock()
	item, ok := m.storage.Gauge[name]

	if !ok {
		return 0, errors.Join(apperrors.ErrNotFound, fmt.Errorf("gauge %s not found", name))
	}

	return item, nil
}

func (m *MetricRepository) GetCounterItem(_ context.Context, name string) (int64, error) {
	m.storage.Lock()
	defer m.storage.Unlock()
	item, ok := m.storage.Counter[name]

	if !ok {
		return 0, errors.Join(apperrors.ErrNotFound, fmt.Errorf("counter %s not found", name))
	}

	return item, nil
}

func (m *MetricRepository) CheckConnection(_ context.Context) bool {
	return true
}

func (m *MetricRepository) CreateBatch(ctx context.Context, batch []dto.CreateMetric) error {
	for _, metric := range batch {
		switch metric.Type {
		case config.GaugeType:
			value, err := strconv.ParseFloat(metric.Value, 64)
			if err != nil {
				return errors.Join(apperrors.ErrBadRequest, fmt.Errorf("cannot read [%s] as gauge", metric.Value))
			}

			if err = m.AddGaugeValue(ctx, metric.Name, value); err != nil {
				return fmt.Errorf("error creating batch for gauge %s: %w", metric.Name, err)
			}
		case config.CounterType:
			value, err := strconv.ParseInt(metric.Value, 10, 64)
			if err != nil {
				return errors.Join(apperrors.ErrBadRequest, fmt.Errorf("cannot read [%s] as counter", metric.Value))
			}

			if err = m.AddCounterValue(ctx, metric.Name, value); err != nil {
				return fmt.Errorf("error creating batch for counter %s: %w", metric.Name, err)
			}
		}
	}
	return nil
}
