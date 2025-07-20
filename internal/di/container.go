// Package di содержит контейнер зависимостей (Dependency Injection).
// Предоставляет простую реализацию IoC контейнера для управления зависимостями приложения.
package di

import (
	"fmt"
	"sync"
)

// ContainerInterface определяет интерфейс для контейнера зависимостей.
// Предоставляет методы для добавления и получения зависимостей по ключу.
type ContainerInterface interface {
	// Add добавляет зависимость в контейнер по ключу.
	// Возвращает ошибку если зависимость с таким ключом уже существует.
	Add(key string, value any) error

	// Get возвращает зависимость из контейнера по ключу.
	// Возвращает nil если зависимость не найдена.
	Get(key string) any
}

// Container реализует простой контейнер зависимостей.
// Использует map для хранения зависимостей и RWMutex для потокобезопасности.
type Container struct {
	// mu - мьютекс для потокобезопасного доступа к зависимостям
	mu sync.RWMutex
	// deps - map для хранения зависимостей
	// Ключ - строка, значение - любой тип
	deps map[string]any
}

// NewContainer создает новый экземпляр контейнера зависимостей.
// Инициализирует пустой map для хранения зависимостей.
func NewContainer() ContainerInterface {
	return &Container{
		deps: make(map[string]any),
	}
}

// Add добавляет зависимость в контейнер по ключу.
// Проверяет существование зависимости и возвращает ошибку при дублировании.
// Метод потокобезопасен благодаря использованию мьютекса.
func (c *Container) Add(key string, value any) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, ok := c.deps[key]; ok {
		return fmt.Errorf("container already has a dependency %s", key)
	}

	c.deps[key] = value

	return nil
}

func (c *Container) Get(key string) any {
	c.mu.RLock()
	defer c.mu.RUnlock()

	d, ok := c.deps[key]

	if !ok {
		return nil
	}
	return d
}
