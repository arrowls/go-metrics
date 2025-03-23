package collector

import (
	"math/rand"
	"runtime"
	"sync"
)

type Collector struct {
	Alloc         float64
	BuckHashSys   float64
	Frees         float64
	GCCPUFraction float64
	GCSys         float64
	HeapAlloc     float64
	HeapIdle      float64
	HeapInuse     float64
	HeapObjects   float64
	HeapReleased  float64
	HeapSys       float64
	LastGC        float64
	Lookups       float64
	MCacheInuse   float64
	MCacheSys     float64
	MSpanInuse    float64
	MSpanSys      float64
	Mallocs       float64
	NextGC        float64
	NumForcedGC   float64
	NumGC         float64
	OtherSys      float64
	PauseTotalNs  float64
	StackInuse    float64
	StackSys      float64
	Sys           float64
	TotalAlloc    float64
	PollCount     int64
	RandomValue   float64
	Mutex         sync.Mutex
}

func New() MetricProvider {
	return &Collector{}
}

func (c *Collector) Collect() {
	memstats := new(runtime.MemStats)

	c.Mutex.Lock()
	defer c.Mutex.Unlock()

	c.Alloc = float64(memstats.Alloc)
	c.BuckHashSys = float64(memstats.BuckHashSys)
	c.Frees = float64(memstats.Frees)
	c.GCCPUFraction = memstats.GCCPUFraction
	c.GCSys = float64(memstats.GCSys)
	c.HeapAlloc = float64(memstats.HeapAlloc)
	c.HeapIdle = float64(memstats.HeapIdle)
	c.HeapInuse = float64(memstats.HeapInuse)
	c.HeapObjects = float64(memstats.HeapObjects)
	c.HeapReleased = float64(memstats.HeapReleased)
	c.HeapSys = float64(memstats.HeapSys)
	c.LastGC = float64(memstats.LastGC)
	c.Lookups = float64(memstats.Lookups)
	c.MCacheInuse = float64(memstats.MCacheInuse)
	c.MCacheSys = float64(memstats.MCacheSys)
	c.MSpanInuse = float64(memstats.MSpanInuse)
	c.MSpanSys = float64(memstats.MSpanSys)
	c.Mallocs = float64(memstats.Mallocs)
	c.NextGC = float64(memstats.NextGC)
	c.NumForcedGC = float64(memstats.NumForcedGC)
	c.NumGC = float64(memstats.NumGC)
	c.OtherSys = float64(memstats.OtherSys)
	c.PauseTotalNs = float64(memstats.PauseTotalNs)
	c.StackInuse = float64(memstats.StackInuse)
	c.StackSys = float64(memstats.StackSys)
	c.Sys = float64(memstats.Sys)
	c.TotalAlloc = float64(memstats.TotalAlloc)
	c.PollCount++
	c.RandomValue = rand.Float64()
}

func (c *Collector) AsMap() *map[string]interface{} {
	mapCollector := make(map[string]interface{})

	c.Mutex.Lock()
	defer c.Mutex.Unlock()

	mapCollector["Alloc"] = c.Alloc
	mapCollector["BuckHashSys"] = c.BuckHashSys
	mapCollector["Frees"] = c.Frees
	mapCollector["GCCPUFraction"] = c.GCCPUFraction
	mapCollector["GCSys"] = c.GCSys
	mapCollector["HeapAlloc"] = c.HeapAlloc
	mapCollector["HeapIdle"] = c.HeapIdle
	mapCollector["HeapInuse"] = c.HeapInuse
	mapCollector["HeapObjects"] = c.HeapObjects
	mapCollector["HeapReleased"] = c.HeapReleased
	mapCollector["HeapSys"] = c.HeapSys
	mapCollector["LastGC"] = c.LastGC
	mapCollector["Lookups"] = c.Lookups
	mapCollector["MCacheInuse"] = c.MCacheInuse
	mapCollector["MCacheSys"] = c.MCacheSys
	mapCollector["MSpanInuse"] = c.MSpanInuse
	mapCollector["MSpanSys"] = c.MSpanSys
	mapCollector["Mallocs"] = c.Mallocs
	mapCollector["NextGC"] = c.NextGC
	mapCollector["NumForcedGC"] = c.NumForcedGC
	mapCollector["NumGC"] = c.NumGC
	mapCollector["OtherSys"] = c.OtherSys
	mapCollector["PauseTotalNs"] = c.PauseTotalNs
	mapCollector["StackInuse"] = c.StackInuse
	mapCollector["StackSys"] = c.StackSys
	mapCollector["Sys"] = c.Sys
	mapCollector["TotalAlloc"] = c.TotalAlloc
	mapCollector["PollCount"] = c.PollCount
	mapCollector["RandomValue"] = c.RandomValue

	return &mapCollector
}
