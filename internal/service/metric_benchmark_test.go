package service

import (
	"context"
	"testing"

	"github.com/arrowls/go-metrics/internal/dto"
	"github.com/arrowls/go-metrics/internal/memstorage"
	"github.com/arrowls/go-metrics/internal/repository"
)

func BenchmarkMetricService_Create(b *testing.B) {
	storage := memstorage.GetInstance()
	repo := repository.NewRepository(storage)
	service := NewMetricService(repo)
	ctx := context.Background()

	createMetric := &dto.CreateMetric{
		Type:  "gauge",
		Name:  "test_metric",
		Value: "123.45",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		service.Create(ctx, createMetric)
	}
}

func BenchmarkMetricService_GetItem(b *testing.B) {
	storage := memstorage.GetInstance()
	repo := repository.NewRepository(storage)
	service := NewMetricService(repo)
	ctx := context.Background()

	// Подготавливаем данные
	createMetric := &dto.CreateMetric{
		Type:  "gauge",
		Name:  "test_metric",
		Value: "123.45",
	}
	service.Create(ctx, createMetric)

	getMetric := &dto.GetMetric{
		Type: "gauge",
		Name: "test_metric",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		service.GetItem(ctx, getMetric)
	}
}

func BenchmarkMetricService_GetList(b *testing.B) {
	storage := memstorage.GetInstance()
	repo := repository.NewRepository(storage)
	service := NewMetricService(repo)
	ctx := context.Background()

	// Подготавливаем данные
	for i := 0; i < 100; i++ {
		createMetric := &dto.CreateMetric{
			Type:  "gauge",
			Name:  "test_metric_" + string(rune(i)),
			Value: "123.45",
		}
		service.Create(ctx, createMetric)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		service.GetList(ctx)
	}
}

func BenchmarkMetricService_CreateBatch(b *testing.B) {
	storage := memstorage.GetInstance()
	repo := repository.NewRepository(storage)
	service := NewMetricService(repo)
	ctx := context.Background()

	// Подготавливаем batch данных
	createMetrics := make([]dto.CreateMetric, 10)
	for i := 0; i < 10; i++ {
		createMetrics[i] = dto.CreateMetric{
			Type:  "gauge",
			Name:  "test_metric_" + string(rune(i)),
			Value: "123.45",
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		service.CreateBatch(ctx, createMetrics)
	}
}

func BenchmarkMetricService_ConcurrentCreate(b *testing.B) {
	storage := memstorage.GetInstance()
	repo := repository.NewRepository(storage)
	service := NewMetricService(repo)
	ctx := context.Background()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			createMetric := &dto.CreateMetric{
				Type:  "gauge",
				Name:  "test_metric_" + string(rune(i)),
				Value: "123.45",
			}
			service.Create(ctx, createMetric)
			i++
		}
	})
}

func BenchmarkMetricService_ConcurrentGetItem(b *testing.B) {
	storage := memstorage.GetInstance()
	repo := repository.NewRepository(storage)
	service := NewMetricService(repo)
	ctx := context.Background()

	// Подготавливаем данные
	createMetric := &dto.CreateMetric{
		Type:  "gauge",
		Name:  "test_metric",
		Value: "123.45",
	}
	service.Create(ctx, createMetric)

	getMetric := &dto.GetMetric{
		Type: "gauge",
		Name: "test_metric",
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			service.GetItem(ctx, getMetric)
		}
	})
}

func BenchmarkNewMetricService(b *testing.B) {
	storage := memstorage.GetInstance()
	repo := repository.NewRepository(storage)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		NewMetricService(repo)
	}
}
