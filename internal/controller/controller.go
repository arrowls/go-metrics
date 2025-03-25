package controller

import (
	"net/http"

	"github.com/arrowls/go-metrics/internal/service"
	"github.com/go-chi/chi/v5"
)

type Metric interface {
	HandleNew(rw http.ResponseWriter, r *http.Request)
	HandleItem(rw http.ResponseWriter, r *http.Request)
}

type Public interface {
	HandlePublic(rw http.ResponseWriter, r *http.Request)
	HandleIndex(rw http.ResponseWriter, r *http.Request)
}

type Controller struct {
	Metric Metric
	Public Public
}

func NewController(services *service.Service) *Controller {
	return &Controller{
		NewMetricController(services),
		NewPublicController(services),
	}
}

func (c *Controller) InitRoutes() *chi.Mux {
	router := chi.NewRouter()

	router.Get("/*", c.Public.HandlePublic)
	router.Get("/", c.Public.HandleIndex)
	router.Post("/update/{type}/{name}/{value}", c.Metric.HandleNew)
	router.Get("/value/{type}/{name}", c.Metric.HandleItem)

	return router
}
