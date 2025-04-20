package mappers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/arrowls/go-metrics/internal/apperrors"
	"github.com/arrowls/go-metrics/internal/dto"
	"github.com/go-chi/chi/v5"
)

func HTTPToCreateMetric(r *http.Request) (*dto.CreateMetric, error) {
	metricType := chi.URLParam(r, "type")
	name := chi.URLParam(r, "name")
	value := chi.URLParam(r, "value")

	if value == "" {
		return nil, errors.Join(apperrors.ErrBadRequest, fmt.Errorf("не указано значение метрики"))
	}

	if name == "" {
		return nil, errors.Join(apperrors.ErrNotFound, fmt.Errorf("не указано имя метрики"))
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
		return nil, errors.Join(apperrors.ErrNotFound, fmt.Errorf("не указано имя метрики"))
	}
	if metricType != "gauge" && metricType != "counter" {
		return nil, errors.Join(apperrors.ErrNotFound, fmt.Errorf("неизвестный тип метрики: %s", metricType))
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
		return nil, errors.Join(apperrors.ErrBadRequest, fmt.Errorf("не удалось прочитать тело запроса"))
	}

	err = requestBody.Validate()
	if err != nil {
		return nil, err
	}

	switch requestBody.MType {
	case "gauge":
		if requestBody.Value == nil {
			return nil, errors.Join(apperrors.ErrBadRequest, fmt.Errorf("не указано значение gauge"))
		}
		return &dto.CreateMetric{
			Type:  requestBody.MType,
			Value: strconv.FormatFloat(*requestBody.Value, 'f', -1, 64),
			Name:  requestBody.ID,
		}, nil
	case "counter":
		if requestBody.Delta == nil {
			return nil, errors.Join(apperrors.ErrBadRequest, fmt.Errorf("не указано значение counter"))
		}
		return &dto.CreateMetric{
			Type:  requestBody.MType,
			Value: fmt.Sprintf("%d", *requestBody.Delta),
			Name:  requestBody.ID,
		}, nil

	default:
		return nil, errors.Join(apperrors.ErrNotFound, fmt.Errorf("неизвестный тип метрики: %s", requestBody.MType))
	}
}

func CreateDTOToHTTPRes(createDto *dto.CreateMetric, value string) ([]byte, error) {
	var metric *dto.Metrics
	switch createDto.Type {
	case "gauge":
		val, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return nil, fmt.Errorf("ошибка при чтении значения %s", createDto.Name)
		}

		metric = &dto.Metrics{
			ID:    createDto.Name,
			MType: createDto.Type,
			Value: &val,
		}
	case "counter":
		deltaVal, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("ошибка при чтении значения counter %s", createDto.Name)
		}

		metric = &dto.Metrics{
			ID:    createDto.Name,
			MType: createDto.Type,
			Delta: &deltaVal,
		}

	default:
		return nil, errors.Join(apperrors.ErrNotFound, fmt.Errorf("неизвестный тип метрики: %s", createDto.Type))
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
		return nil, errors.Join(apperrors.ErrBadRequest, fmt.Errorf("не удалось прочитать тело запроса"))
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
		metric.MType = "gauge"
	case int64:
		metric.Delta = &v
		metric.MType = "counter"
	default:
		return nil, fmt.Errorf("unknow metric value: %d of type %T", value, value)
	}

	return metric, nil
}
