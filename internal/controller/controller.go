package controller

import (
	"net/http"

	"github.com/arrowls/go-metrics/internal/middleware"
	"github.com/arrowls/go-metrics/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
)

type Metric interface {
	HandleNew(rw http.ResponseWriter, r *http.Request)
	HandleItem(rw http.ResponseWriter, r *http.Request)
	HandleNewFromBody(rw http.ResponseWriter, r *http.Request)
	HandleGetItemFromBody(rw http.ResponseWriter, r *http.Request)
}

type Public interface {
	HandlePublic(rw http.ResponseWriter, r *http.Request)
	HandleIndex(rw http.ResponseWriter, r *http.Request)
}

type ErrorHandler interface {
	Handle(w http.ResponseWriter, err error)
}

type Controller struct {
	Metric Metric
	Public Public
}

func NewController(services *service.Service, handler ErrorHandler) *Controller {
	return &Controller{
		NewMetricController(services, handler),
		NewPublicController(services),
	}
}

func (c *Controller) InitRoutes(loggerInst *logrus.Logger) *chi.Mux {
	router := chi.NewRouter()
	router.Use(middleware.NewCompressionMiddleware)
	router.Use(middleware.NewLoggingMiddleware(loggerInst))

	router.Get("/assets/*", c.Public.HandlePublic)
	router.Get("/", c.Public.HandleIndex)
	router.Post("/update/{type}/{name}/{value}", c.Metric.HandleNew)
	router.Post("/update", c.Metric.HandleNewFromBody)
	router.Post("/value", c.Metric.HandleGetItemFromBody)
	router.Get("/value/{type}/{name}", c.Metric.HandleItem)

	return router
}
