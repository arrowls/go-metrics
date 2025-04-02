package memstorage

import (
	"sync"
)

type MemStorage struct {
	Gauge   map[string]float64
	Counter map[string]int64
	*sync.RWMutex
}

var instance *MemStorage

var once sync.Once

func GetInstance() *MemStorage {
	once.Do(func() {
		instance = &MemStorage{
			Gauge:   make(map[string]float64),
			Counter: make(map[string]int64),
		}
	})
	return instance
}
