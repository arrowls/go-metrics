package service

import (
	"github.com/arrowls/go-metrics/internal/repository"
)

type Metric interface {
	CreateByType(metricType string, name string, stringValue string) error
	GetList() *map[string]interface{}
	GetItem(metricType string, name string) (string, error)
}

type Service struct {
	Metric Metric
}

func NewService(repo *repository.Repository) *Service {
	return &Service{
		Metric: NewMetricService(repo),
	}
}
