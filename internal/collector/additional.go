package collector

import (
	"sync"

	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/mem"
)

type AdditionalCollector struct {
	TotalMemory     float64
	FreeMemory      float64
	CPUutilization1 float64
	Mutex           sync.Mutex
}

func NewAdditionalCollector() MetricProvider {
	return &AdditionalCollector{}
}

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
