package collector

import (
	"testing"
)

func BenchmarkCollector_Collect(b *testing.B) {
	collector := New().(*Collector)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		collector.Collect()
	}
}

func BenchmarkCollector_AsMap(b *testing.B) {
	collector := New().(*Collector)
	collector.Collect() // Собираем метрики один раз

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		collector.AsMap()
	}
}

func BenchmarkCollector_CollectAndAsMap(b *testing.B) {
	collector := New().(*Collector)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		collector.Collect()
		collector.AsMap()
	}
}

func BenchmarkCollector_ConcurrentCollect(b *testing.B) {
	collector := New().(*Collector)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			collector.Collect()
		}
	})
}

func BenchmarkCollector_ConcurrentAsMap(b *testing.B) {
	collector := New().(*Collector)
	collector.Collect() // Собираем метрики один раз

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			collector.AsMap()
		}
	})
}

func BenchmarkNew(b *testing.B) {
	for i := 0; i < b.N; i++ {
		New()
	}
}
