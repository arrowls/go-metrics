package service

import (
	"errors"
	"strconv"

	"github.com/arrowls/go-metrics/cmd/server/repository"
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
