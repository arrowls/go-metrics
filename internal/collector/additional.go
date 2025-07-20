// Package collector содержит компоненты для сбора системных метрик.
package collector

import (
	"sync"

	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/mem"
)

// AdditionalCollector собирает дополнительные системные метрики.
// Использует библиотеку gopsutil для получения информации о памяти и CPU.
type AdditionalCollector struct {
	// TotalMemory - общий объем памяти системы в байтах
	TotalMemory float64
	// FreeMemory - свободный объем памяти в байтах
	FreeMemory float64
	// CPUutilization1 - количество логических CPU ядер
	CPUutilization1 float64
	// Mutex - мьютекс для потокобезопасного доступа к метрикам
	Mutex sync.Mutex
}

// NewAdditionalCollector создает новый экземпляр дополнительного коллектора.
// Возвращает интерфейс MetricProvider для работы с метриками.
func NewAdditionalCollector() MetricProvider {
	return &AdditionalCollector{}
}

// Collect собирает дополнительные системные метрики.
// Читает информацию о памяти и CPU с помощью библиотеки gopsutil.
// Метод потокобезопасен благодаря использованию мьютекса.
func (c *AdditionalCollector) Collect() {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()

	memory, err := mem.VirtualMemory()
	if err != nil {
		return
	}

	c.TotalMemory = float64(memory.Total)
	c.FreeMemory = float64(memory.Free)

	countCPU, err := cpu.Counts(true)
	if err != nil {
		return
	}

	c.CPUutilization1 = float64(countCPU)
}

func (c *AdditionalCollector) AsMap() *map[string]interface{} {
	mapCollector := make(map[string]interface{})

	c.Mutex.Lock()
	defer c.Mutex.Unlock()

	mapCollector["TotalMemory"] = c.TotalMemory
	mapCollector["FreeMemory"] = c.FreeMemory
	mapCollector["CPUutilization1"] = c.CPUutilization1

	return &mapCollector
}
