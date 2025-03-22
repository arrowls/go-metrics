package controller

import (
	"net/http"
	"strings"

	"github.com/arrowls/go-metrics/cmd/server/service"
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
	if r.Method != http.MethodPost {
		rw.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// /update/{type}/{name}/{value}
	urlParts := strings.Split(r.URL.Path, "/")

	if len(urlParts) < 5 {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	var (
		metricType = urlParts[2]
		name       = urlParts[3]
		value      = urlParts[4]
	)

	if value == "" {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	if name == "" {
		rw.WriteHeader(http.StatusNotFound)
		return
	}

	err := c.service.Metric.CreateByType(metricType, name, value)

	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
	}
}
