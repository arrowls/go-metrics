// Package dto содержит объекты передачи данных (Data Transfer Objects).
// Включает структуры для создания, получения и валидации метрик.
package dto

import (
	"errors"
	"fmt"

	"github.com/arrowls/go-metrics/internal/apperrors"
	"github.com/arrowls/go-metrics/internal/config"
)

// CreateMetric представляет DTO для создания новой метрики.
// Используется для передачи данных от HTTP контроллера к сервису.
type CreateMetric struct {
	// Type - тип метрики (gauge или counter)
	Type string
	// Name - имя метрики
	Name string
	// Value - значение метрики в строковом формате
	Value string
}

// GetMetric представляет DTO для получения метрики.
// Используется для передачи параметров запроса к сервису.
type GetMetric struct {
	// Type - тип метрики (gauge или counter)
	Type string
	// Name - имя метрики
	Name string
}

// Metrics представляет DTO для JSON сериализации метрик.
// Используется для HTTP API с JSON форматом данных.
type Metrics struct {
	// ID - уникальный идентификатор метрики
	ID string `json:"id"`
	// MType - тип метрики (gauge или counter)
	MType string `json:"type"`
	// Delta - значение для counter метрик (опционально)
	Delta *int64 `json:"delta,omitempty"`
	// Value - значение для gauge метрик (опционально)
	Value *float64 `json:"value,omitempty"`
}

// Validate проверяет корректность данных метрики.
// Проверяет наличие имени и корректность типа метрики.
// Возвращает ошибку если данные некорректны.
func (m *Metrics) Validate() error {
	if m.ID == "" {
		return errors.Join(apperrors.ErrBadRequest, fmt.Errorf("metric name not specified"))
	}
	if m.MType != config.GaugeType && m.MType != config.CounterType {
		return errors.Join(apperrors.ErrNotFound, fmt.Errorf("unknown metric type: %s", m.MType))
	}
	return nil
}
