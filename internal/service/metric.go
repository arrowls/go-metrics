package service

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/arrowls/go-metrics/internal/apperrors"
	"github.com/arrowls/go-metrics/internal/dto"
	"github.com/arrowls/go-metrics/internal/repository"
)

type MetricService struct {
	repository repository.Repository
}

func NewMetricService(repository *repository.Repository) *MetricService {
	return &MetricService{
		*repository,
	}
}

func (m *MetricService) Create(dto *dto.CreateMetric) error {
	if dto.Type == "gauge" {
		parsedValue, err := strconv.ParseFloat(dto.Value, 64)
		if err != nil {
			return errors.Join(apperrors.ErrBadRequest, err)
		}

		m.repository.Metric.AddGaugeValue(dto.Name, parsedValue)
		return nil
	}

	if dto.Type == "counter" {
		parsedValue, err := strconv.ParseInt(dto.Value, 10, 64)
		if err != nil {
			return errors.Join(apperrors.ErrBadRequest, err)
		}

		m.repository.Metric.AddCounterValue(dto.Name, parsedValue)
		return nil
	}

	return errors.Join(apperrors.ErrBadRequest, fmt.Errorf("неизвестный тип метрики: %s", dto.Type))
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

func (m *MetricService) GetItem(dto *dto.GetMetric) (string, error) {
	if dto.Type == "gauge" {
		value, err := m.repository.Metric.GetGaugeItem(dto.Name)

		if err != nil {
			return "", err
		}

		return strconv.FormatFloat(value, 'f', -1, 64), nil
	}

	if dto.Type == "counter" {
		value, err := m.repository.Metric.GetCounterItem(dto.Name)

		if err != nil {
			return "", err
		}

		return strconv.FormatInt(value, 10), nil
	}

	return "", errors.Join(apperrors.ErrBadRequest, fmt.Errorf("неизвестный тип метрики: %s", dto.Type))
}
