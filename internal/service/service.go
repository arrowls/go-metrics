package service

import (
	"context"

	"github.com/arrowls/go-metrics/internal/dto"
)

type Metric interface {
	Create(ctx context.Context, dto *dto.CreateMetric) error
	GetList(ctx context.Context) *map[string]interface{}
	GetItem(ctx context.Context, dto *dto.GetMetric) (string, error)
	CheckConnection(ctx context.Context) bool
	CreateBatch(ctx context.Context, batch []dto.CreateMetric) error
}

type Service struct {
	Metric Metric
}
