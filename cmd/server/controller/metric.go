package controller

import (
	"net/http"

	"github.com/arrowls/go-metrics/cmd/server/service"
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

	if value == "" {
		http.Error(rw, "No metric value specified", http.StatusBadRequest)
		return
	}

	if name == "" {
		http.Error(rw, "No metric name specified", http.StatusNotFound)
		return
	}

	err := c.service.Metric.CreateByType(metricType, name, value)

	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
	}
}
