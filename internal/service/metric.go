package service

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/arrowls/go-metrics/internal/apperrors"
	"github.com/arrowls/go-metrics/internal/config"
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

func (m *MetricService) Create(ctx context.Context, dto *dto.CreateMetric) error {
	if dto.Type == config.GaugeType {
		parsedValue, err := strconv.ParseFloat(dto.Value, 64)
		if err != nil {
			return errors.Join(apperrors.ErrBadRequest, err)
		}

		err = m.repository.Metric.AddGaugeValue(ctx, dto.Name, parsedValue)
		return err
	}

	if dto.Type == config.CounterType {
		parsedValue, err := strconv.ParseInt(dto.Value, 10, 64)
		if err != nil {
			return errors.Join(apperrors.ErrBadRequest, err)
		}

		err = m.repository.Metric.AddCounterValue(ctx, dto.Name, parsedValue)
		return err
	}
	return errors.Join(apperrors.ErrBadRequest, fmt.Errorf("unknown metric type: %s", dto.Type))
}

func (m *MetricService) GetList(ctx context.Context) *map[string]interface{} {
	storage, err := m.repository.Metric.GetAll(ctx)

	returnMap := make(map[string]interface{})

	if err != nil {
		return &returnMap
	}

	for k, v := range storage.Gauge {
		returnMap[k] = v
	}

	for k, v := range storage.Counter {
		returnMap[k] = v
	}

	return &returnMap
}

func (m *MetricService) GetItem(ctx context.Context, dto *dto.GetMetric) (string, error) {
	if dto.Type == config.GaugeType {
		value, err := m.repository.Metric.GetGaugeItem(ctx, dto.Name)
		if err != nil {
			return "", err
		}

		return strconv.FormatFloat(value, 'f', -1, 64), nil
	}

	if dto.Type == config.CounterType {
		value, err := m.repository.Metric.GetCounterItem(ctx, dto.Name)

		if err != nil {
			return "", err
		}

		return strconv.FormatInt(value, 10), nil
	}

	return "", errors.Join(apperrors.ErrBadRequest, fmt.Errorf("unknown metric typef: %s", dto.Type))
}

func (m *MetricService) CheckConnection(ctx context.Context) bool {
	return m.repository.Metric.CheckConnection(ctx)
}

func (m *MetricService) CreateBatch(ctx context.Context, batch []dto.CreateMetric) error {
	return m.repository.Metric.CreateBatch(ctx, batch)
}
