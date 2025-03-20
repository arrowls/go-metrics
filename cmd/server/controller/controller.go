package controller

import (
	"net/http"

	"github.com/arrowls/go-metrics/cmd/server/middleware"
	"github.com/arrowls/go-metrics/cmd/server/service"
)

type Metric interface {
	HandleNew(rw http.ResponseWriter, r *http.Request)
}

type Controller struct {
	Metric Metric
}

func NewController(services *service.Service) *Controller {
	return &Controller{
		NewMetricController(services),
	}
}

func (c *Controller) InitRoutes() *http.ServeMux {
	mux := http.NewServeMux()

	mux.Handle("/update/", middleware.Wrap(
		http.HandlerFunc(c.Metric.HandleNew),
		[]middleware.Middleware{
			middleware.Logger, // example
		},
	))

	return mux
}
