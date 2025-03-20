package service

import (
	"github.com/arrowls/go-metrics/cmd/server/repository"
)

type Metric interface {
	CreateByType(metricType string, name string, stringValue string) error
}

type Service struct {
	Metric Metric
}

func NewService(repo *repository.Repository) *Service {
	return &Service{
		Metric: NewMetricService(repo),
	}
}
