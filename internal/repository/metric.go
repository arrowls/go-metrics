package repository

import (
	"errors"
	"fmt"

	"github.com/arrowls/go-metrics/internal/apperrors"
	"github.com/arrowls/go-metrics/internal/memstorage"
)

type MetricRepository struct {
	storage *memstorage.MemStorage
}

func NewMetricRepository(storage *memstorage.MemStorage) *MetricRepository {
	return &MetricRepository{
		storage: storage,
	}
}

func (m *MetricRepository) AddGaugeValue(name string, value float64) {
	m.storage.Lock()
	defer m.storage.Unlock()

	m.storage.Gauge[name] = value
}

func (m *MetricRepository) AddCounterValue(name string, value int64) {
	m.storage.Lock()
	defer m.storage.Unlock()

	m.storage.Counter[name] += value
}

func (m *MetricRepository) GetAll() memstorage.MemStorage {
	return *m.storage
}

func (m *MetricRepository) GetGaugeItem(name string) (float64, error) {
	m.storage.Lock()
	defer m.storage.Unlock()
	item, ok := m.storage.Gauge[name]

	if !ok {
		return 0, errors.Join(apperrors.ErrNotFound, fmt.Errorf("gauge %s not found", name))
	}

	return item, nil
}

func (m *MetricRepository) GetCounterItem(name string) (int64, error) {
	m.storage.Lock()
	defer m.storage.Unlock()
	item, ok := m.storage.Counter[name]

	if !ok {
		return 0, errors.Join(apperrors.ErrNotFound, fmt.Errorf("counter %s not found", name))
	}

	return item, nil
}
