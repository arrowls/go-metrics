package repository

import (
	"encoding/json"
	"fmt"

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

	if m.storage.Gauge == nil {
		m.storage.Gauge = make(map[string]float64)
	}

	m.storage.Gauge[name] = value

	printValue, _ := json.MarshalIndent(m.storage.Gauge, " ", " ")
	fmt.Printf("Current gauge values: %s\n", printValue)
}

func (m *MetricRepository) AddCounterValue(name string, value int64) {
	if m.storage.Counter == nil {
		m.storage.Counter = make(map[string]int64)
	}

	m.storage.Counter[name] += value
	printValue, _ := json.MarshalIndent(m.storage.Counter, " ", " ")
	fmt.Printf("Current counter values: %s\n", printValue)
}
