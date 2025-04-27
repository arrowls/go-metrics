package service

import (
	"github.com/arrowls/go-metrics/internal/dto"
	"github.com/arrowls/go-metrics/internal/repository"
)

type Metric interface {
	Create(dto *dto.CreateMetric) error
	GetList() *map[string]interface{}
	GetItem(dto *dto.GetMetric) (string, error)
}

type Service struct {
	Metric Metric
}

func NewService(repo *repository.Repository) *Service {
	return &Service{
		Metric: NewMetricService(repo),
	}
}
