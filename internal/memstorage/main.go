// Package memstorage содержит хранилище метрик в памяти.
// Реализует паттерн Singleton для обеспечения единственного экземпляра хранилища.
package memstorage

import (
	"sync"
)

// MemStorage представляет хранилище метрик в памяти.
// Содержит отдельные map для gauge и counter метрик.
// Использует RWMutex для потокобезопасного доступа.
type MemStorage struct {
	// Gauge - map для хранения gauge метрик
	// Ключ - имя метрики, значение - float64
	Gauge map[string]float64
	// Counter - map для хранения counter метрик
	// Ключ - имя метрики, значение - int64
	Counter map[string]int64
	// RWMutex - мьютекс для потокобезопасного доступа к хранилищу
	*sync.RWMutex
}

var instance *MemStorage

var once sync.Once

// GetInstance возвращает единственный экземпляр хранилища метрик.
// Реализует паттерн Singleton для обеспечения глобального доступа к хранилищу.
// При первом вызове инициализирует хранилище с пустыми map.
func GetInstance() *MemStorage {
	once.Do(func() {
		instance = &MemStorage{
			Gauge:   make(map[string]float64),
			Counter: make(map[string]int64),
			RWMutex: &sync.RWMutex{},
		}
	})
	return instance
}
