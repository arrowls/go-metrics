// Package repository содержит слой доступа к данным для метрик.
// Предоставляет абстракцию над различными хранилищами (память, файл, база данных).
package repository

import (
	"context"

	"github.com/arrowls/go-metrics/internal/dto"
	"github.com/arrowls/go-metrics/internal/memstorage"
)

// Metric определяет интерфейс для работы с хранилищем метрик.
// Включает методы для добавления, получения и проверки метрик.
type Metric interface {
	// AddGaugeValue добавляет или обновляет gauge метрику.
	// Принимает имя метрики и значение типа float64.
	AddGaugeValue(ctx context.Context, key string, value float64) error

	// AddCounterValue добавляет значение к counter метрике.
	// Counter метрики накапливают значения, поэтому происходит сложение.
	AddCounterValue(ctx context.Context, key string, value int64) error

	// GetAll возвращает все сохраненные метрики.
	// Возвращает структуру MemStorage с gauge и counter метриками.
	GetAll(ctx context.Context) (memstorage.MemStorage, error)

	// GetCounterItem возвращает значение counter метрики по имени.
	// Возвращает int64 значение и ошибку если метрика не найдена.
	GetCounterItem(ctx context.Context, name string) (int64, error)

	// GetGaugeItem возвращает значение gauge метрики по имени.
	// Возвращает float64 значение и ошибку если метрика не найдена.
	GetGaugeItem(ctx context.Context, name string) (float64, error)

	// CheckConnection проверяет доступность хранилища.
	// Возвращает true если хранилище доступно.
	CheckConnection(ctx context.Context) bool

	// CreateBatch создает несколько метрик за одну операцию.
	// Принимает слайс DTO метрик и обрабатывает их пакетно.
	CreateBatch(ctx context.Context, batch []dto.CreateMetric) error
}

// Repository содержит все репозитории приложения.
// Предоставляет централизованный доступ к слою данных.
type Repository struct {
	// Metric - репозиторий для работы с метриками
	Metric Metric
}

// NewRepository создает новый экземпляр репозитория.
// Принимает хранилище метрик и возвращает настроенный репозиторий.
func NewRepository(storage *memstorage.MemStorage) *Repository {
	return &Repository{
		Metric: NewMetricRepository(storage),
	}
}
