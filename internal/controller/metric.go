package controller

import (
	"net/http"

	"github.com/arrowls/go-metrics/internal/service"
	"github.com/arrowls/go-metrics/internal/validator"
	"github.com/go-chi/chi/v5"
)

type MetricController struct {
	service service.Service
}

func NewMetricController(service *service.Service) *MetricController {
	return &MetricController{
		service: *service,
	}
}

func (c *MetricController) HandleNew(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Add("Content-Type", "text/plain")

	metricType := chi.URLParam(r, "type")
	name := chi.URLParam(r, "name")
	value := chi.URLParam(r, "value")

	if valid := validator.ValidateHandleNewRequest(rw, r); !valid {
		return
	}

	err := c.service.Metric.CreateByType(metricType, name, value)

	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
	}
}

func (c *MetricController) HandleItem(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Add("Content-Type", "text/plain")
	metricType := chi.URLParam(r, "type")
	name := chi.URLParam(r, "name")

	if valid := validator.ValidateHandleItemRequest(rw, r); !valid {
		return
	}

	value, err := c.service.Metric.GetItem(metricType, name)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusNotFound)
		return
	}

	_, err = rw.Write([]byte(value))

	if err != nil {
		http.Error(rw, "Произошла ошибка", http.StatusInternalServerError)
		return
	}
}
