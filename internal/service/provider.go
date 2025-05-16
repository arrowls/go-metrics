package service

import (
	"log"

	"github.com/arrowls/go-metrics/internal/di"
	"github.com/arrowls/go-metrics/internal/repository"
)

const diKey = "metric_service"

func ProvideMetricService(container di.ContainerInterface) *Service {
	if instance, ok := container.Get(diKey).(*Service); ok {
		return instance
	}

	repo := repository.ProvideRepository(container)

	service := &Service{
		Metric: NewMetricService(repo),
	}

	if err := container.Add(diKey, service); err != nil {
		log.Fatal(err)
	}

	return service
}
