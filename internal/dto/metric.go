package dto

import (
	"errors"
	"fmt"

	"github.com/arrowls/go-metrics/internal/apperrors"
	"github.com/arrowls/go-metrics/internal/config"
)

type CreateMetric struct {
	Type  string
	Name  string
	Value string
}

type GetMetric struct {
	Type string
	Name string
}

type Metrics struct {
	ID    string   `json:"id"`
	MType string   `json:"type"`
	Delta *int64   `json:"delta,omitempty"`
	Value *float64 `json:"value,omitempty"`
}

func (m *Metrics) Validate() error {
	if m.ID == "" {
		return errors.Join(apperrors.ErrBadRequest, fmt.Errorf("metric name not specified"))
	}
	if m.MType != config.GaugeType && m.MType != config.CounterType {
		return errors.Join(apperrors.ErrNotFound, fmt.Errorf("unknown metric type: %s", m.MType))
	}
	return nil
}
