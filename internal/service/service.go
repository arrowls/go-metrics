// Package service содержит бизнес-логику для работы с метриками.
// Предоставляет интерфейсы и реализации для создания, получения и управления метриками.
package service

import (
	"context"

	"github.com/arrowls/go-metrics/internal/dto"
)

// Metric определяет интерфейс для работы с метриками.
// Включает методы для создания, получения и проверки метрик.
type Metric interface {
	// Create создает новую метрику или обновляет существующую.
	// Принимает DTO с данными метрики и возвращает ошибку в случае неудачи.
	Create(ctx context.Context, dto *dto.CreateMetric) error

	// GetList возвращает все доступные метрики в виде map.
	// Ключи - имена метрик, значения - их текущие значения.
	GetList(ctx context.Context) *map[string]interface{}

	// GetItem возвращает значение конкретной метрики по имени и типу.
	// Возвращает строковое представление значения и ошибку.
	GetItem(ctx context.Context, dto *dto.GetMetric) (string, error)

	// CheckConnection проверяет доступность хранилища метрик.
	// Возвращает true если соединение активно, false в противном случае.
	CheckConnection(ctx context.Context) bool

	// CreateBatch создает несколько метрик за одну операцию.
	// Принимает слайс DTO метрик и возвращает ошибку в случае неудачи.
	CreateBatch(ctx context.Context, batch []dto.CreateMetric) error
}

// Service содержит все сервисы приложения.
// Предоставляет централизованный доступ к бизнес-логике.
type Service struct {
	// Metric - сервис для работы с метриками
	Metric Metric
}
