package controller

import (
	"net/http"

	"github.com/arrowls/go-metrics/cmd/server/service"
	"github.com/go-chi/chi/v5"
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

func (c *Controller) InitRoutes() *chi.Mux {
	router := chi.NewRouter()

	router.Post("/update/{type}/{name}/{value}", c.Metric.HandleNew)

	return router
}
