package mappers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/arrowls/go-metrics/internal/apperrors"
	"github.com/arrowls/go-metrics/internal/config"
	"github.com/arrowls/go-metrics/internal/dto"
	"github.com/go-chi/chi/v5"
)

func HTTPToCreateMetric(r *http.Request) (*dto.CreateMetric, error) {
	metricType := chi.URLParam(r, "type")
	name := chi.URLParam(r, "name")
	value := chi.URLParam(r, "value")
	if value == "" {
		return nil, errors.Join(apperrors.ErrBadRequest, fmt.Errorf("metric value not specified"))
	}
	if name == "" {
		return nil, errors.Join(apperrors.ErrNotFound, fmt.Errorf("metric name not specified"))
	}

	return &dto.CreateMetric{
		Type:  metricType,
		Name:  name,
		Value: value,
	}, nil
}

func HTTPToGetMetric(r *http.Request) (*dto.GetMetric, error) {
	metricType := chi.URLParam(r, "type")
	name := chi.URLParam(r, "name")

	if name == "" {
		return nil, errors.Join(apperrors.ErrNotFound, fmt.Errorf("metric name is not specified"))
	}
	if metricType != config.GaugeType && metricType != config.CounterType {
		return nil, errors.Join(apperrors.ErrNotFound, fmt.Errorf("unknown metric type: %s", metricType))
	}

	return &dto.GetMetric{
		Type: metricType,
		Name: name,
	}, nil
}

func HTTPWithBodyToCreateMetric(r *http.Request) (*dto.CreateMetric, error) {
	var requestBody dto.Metrics
	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		return nil, errors.Join(apperrors.ErrBadRequest, fmt.Errorf("failed to read the body of the request: %w", err))
	}

	err = requestBody.Validate()
	if err != nil {
		return nil, err
	}

	switch requestBody.MType {
	case config.GaugeType:
		if requestBody.Value == nil {
			return nil, errors.Join(apperrors.ErrBadRequest, fmt.Errorf("no value for gauge specified"))
		}

		return &dto.CreateMetric{Type: requestBody.MType, Value: strconv.FormatFloat(*requestBody.Value, 'f', -1, 64), Name: requestBody.ID}, nil
	case config.CounterType:
		if requestBody.Delta == nil {
			return nil, errors.Join(apperrors.ErrBadRequest, fmt.Errorf("no value for counter specified"))
		}
		return &dto.CreateMetric{Type: requestBody.MType, Value: fmt.Sprintf("%d", *requestBody.Delta), Name: requestBody.ID}, nil
	default:
		return nil, errors.Join(apperrors.ErrNotFound, fmt.Errorf("unknown metric type: %s", requestBody.MType))
	}
}

func CreateDTOToHTTPRes(createDto *dto.CreateMetric, value string) ([]byte, error) {
	var metric *dto.Metrics
	switch createDto.Type {
	case config.GaugeType:
		val, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return nil, fmt.Errorf("error while reading the value %s", createDto.Name)
		}

		metric = &dto.Metrics{ID: createDto.Name, MType: createDto.Type, Value: &val}
	case config.CounterType:
		deltaVal, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("error while reading the value of counter %s", createDto.Name)
		}
		metric = &dto.Metrics{ID: createDto.Name, MType: createDto.Type, Delta: &deltaVal}
	default:
		return nil, errors.Join(apperrors.ErrNotFound, fmt.Errorf("unknown metric type: %s", createDto.Type))
	}

	jsonResp, err := json.Marshal(metric)
	if err != nil {
		return nil, fmt.Errorf("error while converting to JSON: %v", err)
	}
	return jsonResp, nil
}

func HTTPWithBodyToGetMetric(r *http.Request) (*dto.GetMetric, error) {
	var requestBody dto.Metrics
	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		return nil, errors.Join(apperrors.ErrBadRequest, fmt.Errorf("could not read the request body"))
	}

	err = requestBody.Validate()
	if err != nil {
		return nil, err
	}

	return &dto.GetMetric{
		Type: requestBody.MType,
		Name: requestBody.ID,
	}, nil
}

func MetricToDTO(name string, value interface{}) (*dto.Metrics, error) {
	metric := &dto.Metrics{
		ID: name,
	}

	switch v := value.(type) {
	case float64:
		metric.Value = &v
		metric.MType = config.GaugeType
	case int64:
		metric.Delta = &v
		metric.MType = config.CounterType
	default:
		return nil, fmt.Errorf("unknow metric value: %d of type %T", value, value)
	}

	return metric, nil
}

func HTTPToCreateMetrics(r *http.Request) ([]dto.CreateMetric, error) {
	var metrics []dto.Metrics

	err := json.NewDecoder(r.Body).Decode(&metrics)
	if err != nil {
		return nil, errors.Join(apperrors.ErrBadRequest, fmt.Errorf("invalid request body"))
	}

	createMetrics := make([]dto.CreateMetric, 0, len(metrics))

	for _, metric := range metrics {
		if err = metric.Validate(); err != nil {
			return nil, errors.Join(apperrors.ErrBadRequest, fmt.Errorf("validate failed on %s", metric.ID))
		}

		createMetric := dto.CreateMetric{
			Type: metric.MType,
			Name: metric.ID,
		}

		switch metric.MType {
		case config.CounterType:
			if metric.Delta == nil {
				return nil, errors.Join(apperrors.ErrBadRequest, fmt.Errorf("value is missing for [%s]", metric.ID))
			}
			createMetric.Value = fmt.Sprintf("%d", *metric.Delta)
		case config.GaugeType:
			if metric.Value == nil {
				return nil, errors.Join(apperrors.ErrBadRequest, fmt.Errorf("value is missing for [%s]", metric.ID))
			}
			createMetric.Value = strconv.FormatFloat(*metric.Value, 'f', -1, 64)
		default:
			return nil, fmt.Errorf("unknow metric value type [%s]", metric.MType)
		}

		createMetrics = append(createMetrics, createMetric)
	}

	return createMetrics, nil
}
