package service

import (
	"context"

	"github.com/arrowls/go-metrics/internal/dto"
	"github.com/arrowls/go-metrics/internal/repository"
)

type Metric interface {
	Create(ctx context.Context, dto *dto.CreateMetric) error
	GetList(ctx context.Context) *map[string]interface{}
	GetItem(ctx context.Context, dto *dto.GetMetric) (string, error)
	CheckConnection(ctx context.Context) bool
}

type Service struct {
	Metric Metric
}

func NewService(repo *repository.Repository) *Service {
	return &Service{
		Metric: NewMetricService(repo),
	}
}
