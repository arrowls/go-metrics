package controller

import (
	"log"

	"github.com/arrowls/go-metrics/internal/apperrors"
	"github.com/arrowls/go-metrics/internal/di"
	"github.com/arrowls/go-metrics/internal/middleware"
	"github.com/arrowls/go-metrics/internal/service"
	"github.com/go-chi/chi/v5"
)

const diKey = "http_controller"

func ProvideHTTPController(container di.ContainerInterface) *Controller {
	controllerInst := container.Get(diKey)
	if controller, ok := controllerInst.(*Controller); ok {
		return controller
	}

	httpHandler := apperrors.ProvideHTTPErrorHandler(container)
	services := service.ProvideMetricService(container)

	controller := &Controller{
		NewMetricController(services, httpHandler),
		NewPublicController(services),
	}

	if err := container.Add(diKey, controller); err != nil {
		log.Fatal(err)
	}

	return controller
}

const routerDiKey = "router"

func ProvideRouter(container di.ContainerInterface) *chi.Mux {
	routerInst := container.Get(routerDiKey)
	if router, ok := routerInst.(*chi.Mux); ok {
		return router
	}

	router := chi.NewRouter()
	c := ProvideHTTPController(container)

	router.Use(middleware.NewCompressionMiddleware)
	router.Use(middleware.ProvideLoggingMiddleware(container))
	router.Use(middleware.ProvideHashingMiddleware(container))

	router.Get("/assets/*", c.Public.HandlePublic)
	router.Get("/", c.Public.HandleIndex)
	router.Post("/update/{type}/{name}/{value}", c.Metric.HandleNew)
	router.Post("/update", c.Metric.HandleNewFromBody)
	router.Post("/updates", c.Metric.HandleCreateBatch)
	router.Post("/value", c.Metric.HandleGetItemFromBody)
	router.Get("/value/{type}/{name}", c.Metric.HandleItem)
	router.Get("/ping", c.Public.Ping)

	if err := container.Add(routerDiKey, router); err != nil {
		log.Fatal(err)
	}

	return router
}
