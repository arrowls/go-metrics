package service

import (
	"errors"
	"strconv"

	"github.com/arrowls/go-metrics/internal/repository"
)

type MetricService struct {
	repository repository.Repository
	// global configs, db access, etc
}

func NewMetricService(repository *repository.Repository) *MetricService {
	return &MetricService{
		*repository,
	}
}

func (m *MetricService) CreateByType(metricType string, name string, stringValue string) error {
	if metricType == "gauge" {
		parsedValue, err := strconv.ParseFloat(stringValue, 64)

		if err != nil {
			return err
		}

		m.repository.Metric.AddGaugeValue(name, parsedValue)
		return nil
	}

	if metricType == "counter" {
		parsedValue, err := strconv.ParseInt(stringValue, 10, 64)
		if err != nil {
			return err
		}

		m.repository.Metric.AddCounterValue(name, parsedValue)
		return nil
	}

	return errors.New("invalid metric type")
}

func (m *MetricService) GetList() *map[string]interface{} {
	storage := m.repository.Metric.GetAll()

	returnMap := make(map[string]interface{})

	for k, v := range storage.Gauge {
		returnMap[k] = v
	}

	for k, v := range storage.Counter {
		returnMap[k] = v
	}

	return &returnMap
}

func (m *MetricService) GetItem(metricType string, name string) (string, error) {
	if metricType == "gauge" {
		value, err := m.repository.Metric.GetGaugeItem(name)

		if err != nil {
			return "", err
		}

		return strconv.FormatFloat(value, 'f', -1, 64), nil
	}

	if metricType == "counter" {
		value, err := m.repository.Metric.GetCounterItem(name)

		if err != nil {
			return "", err
		}

		return strconv.FormatInt(value, 10), nil
	}

	return "", errors.New("invalid metric type")
}
