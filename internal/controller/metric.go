package controller

import (
	"fmt"
	"net/http"

	"github.com/arrowls/go-metrics/internal/apperrors"
	"github.com/arrowls/go-metrics/internal/dto"
	"github.com/arrowls/go-metrics/internal/mappers"
	"github.com/arrowls/go-metrics/internal/service"
)

type MetricController struct {
	service      service.Service
	errorHandler ErrorHandler
}

func NewMetricController(service *service.Service, handler ErrorHandler) *MetricController {
	return &MetricController{
		service:      *service,
		errorHandler: handler,
	}
}

func (c *MetricController) HandleNew(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Add("Content-Type", "text/plain")

	createDto, err := mappers.HTTPToCreateMetric(r)
	if err != nil {
		c.errorHandler.Handle(rw, fmt.Errorf("failed to read request: %w", err))
		return
	}

	err = c.service.Metric.Create(r.Context(), createDto)
	if err != nil {
		c.errorHandler.Handle(rw, fmt.Errorf("failed to create metric: %w", err))
	}
}

func (c *MetricController) HandleItem(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Add("Content-Type", "text/plain")

	getItemDto, err := mappers.HTTPToGetMetric(r)
	if err != nil {
		c.errorHandler.Handle(rw, fmt.Errorf("failed to read request: %w", err))
		return
	}

	value, err := c.service.Metric.GetItem(r.Context(), getItemDto)
	if err != nil {
		c.errorHandler.Handle(rw, fmt.Errorf("failed to get metric: %w", err))
		return
	}

	_, err = rw.Write([]byte(value))
	if err != nil {
		c.errorHandler.Handle(rw, apperrors.ErrUnknown)
	}
}

func (c *MetricController) HandleNewFromBody(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Add("Content-Type", "application/json")

	createDto, err := mappers.HTTPWithBodyToCreateMetric(r)
	if err != nil {
		c.errorHandler.Handle(rw, fmt.Errorf("failed to read request: %w", err))
		return
	}

	err = c.service.Metric.Create(r.Context(), createDto)
	if err != nil {
		c.errorHandler.Handle(rw, fmt.Errorf("failed to create metric: %w", err))
		return
	}

	updatedValue, err := c.service.Metric.GetItem(r.Context(), &dto.GetMetric{
		Type: createDto.Type,
		Name: createDto.Name,
	})
	if err != nil {
		c.errorHandler.Handle(rw, fmt.Errorf("failed to get metric: %w", err))
		return
	}

	response, err := mappers.CreateDTOToHTTPRes(createDto, updatedValue)
	if err != nil {
		c.errorHandler.Handle(rw, apperrors.ErrUnknown)
	}

	_, err = rw.Write(response)
	if err != nil {
		c.errorHandler.Handle(rw, apperrors.ErrUnknown)
	}
}

func (c *MetricController) HandleGetItemFromBody(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Add("Content-Type", "application/json")

	fmt.Printf("123")
	getItemDto, err := mappers.HTTPWithBodyToGetMetric(r)
	if err != nil {
		c.errorHandler.Handle(rw, fmt.Errorf("error reading request: %w ", err))
		return
	}

	value, err := c.service.Metric.GetItem(r.Context(), getItemDto)
	if err != nil {
		c.errorHandler.Handle(rw, fmt.Errorf("failed to fetch metric: %w", err))
		return
	}

	response, err := mappers.CreateDTOToHTTPRes(&dto.CreateMetric{
		Type: getItemDto.Type,
		Name: getItemDto.Name,
	}, value)
	if err != nil {
		c.errorHandler.Handle(rw, apperrors.ErrUnknown)
	}

	_, err = rw.Write(response)
	if err != nil {
		c.errorHandler.Handle(rw, apperrors.ErrUnknown)
	}
}

func (c *MetricController) HandleCreateBatch(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Add("Content-Type", "application/json")

	createBatch, err := mappers.HTTPToCreateMetrics(r)
	if err != nil {
		c.errorHandler.Handle(rw, fmt.Errorf("error reading request: %w", err))
		return
	}

	err = c.service.Metric.CreateBatch(r.Context(), createBatch)
	if err != nil {
		c.errorHandler.Handle(rw, fmt.Errorf("error creating batch: %w", err))
	}
}
