//package audit
//
//import (
//	"context"
//	"io"
//	"log/slog"
//	"sync/atomic"
//	"testing"
//
//	"github.com/ArtShib/urlshortener/internal/model"
//)
//
//// -----------------------------------------------------------------------------
//// MOCK SERVICE (Заглушка)
//// -----------------------------------------------------------------------------
//// Мы тестируем Пул, а не базу данных. Заглушка должна работать мгновенно.
//
//type mockEventService struct {
//	processedCount atomic.Int64
//}
//
//func (m *mockEventService) SendAuditRecord(ctx context.Context, record *model.Event) error {
//	// Имитация полезной нагрузки (опционально, сейчас отключено для теста overhead пула)
//	// time.Sleep(1 * time.Millisecond)
//	m.processedCount.Add(1)
//	return nil
//}
//
//func (m *mockEventService) Close() error {
//	return nil
//}
//
//// -----------------------------------------------------------------------------
//// BENCHMARKS
//// -----------------------------------------------------------------------------
//
//// BenchmarkWorkerPool_Throughput проверяет, как быстро пул может "проглатывать" задачи.
//func BenchmarkWorkerPool_Throughput(b *testing.B) {
//	// 1. Setup (отключаем логи, чтобы не тормозили тест)
//	logger := slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{Level: slog.LevelError}))
//
//	svc := &mockEventService{}
//	cfg := &model.WorkerPoolEvent{
//		EventChainSize: 10000, // Большой буфер, чтобы не блокироваться на вставке
//		CountWorkers:   10,    // Фиксированное число воркеров для теста
//	}
//
//	pool := New(svc, logger, cfg)
//
//	// Запускаем пул
//	ctx, cancel := context.WithCancel(context.Background())
//	pool.Start(ctx)
//	defer func() {
//		cancel()
//		pool.Stop()
//	}()
//
//	// Данные для отправки
//	event := &model.Event{
//		UserID: "user-1",
//		Action: "benchmark_test",
//	}
//
//	b.ReportAllocs()
//	b.ResetTimer()
//
//	// 2. Тело бенчмарка (Main Loop)
//	for i := 0; i < b.N; i++ {
//		// Мы тестируем скорость добавления в канал и обработки
//		pool.AddEventRecord(event)
//	}
//
//	// Ждем, пока воркеры разгребут очередь (не обязательно для Throughput, но полезно для точности)
//	// В реальном бенчмарке Throughput мы измеряем скорость *подачи* задач.
//}
//
//// BenchmarkWorkerPool_Parallel проверяет работу под конкурентной нагрузкой (много горутин пишут в пул)
//func BenchmarkWorkerPool_Parallel(b *testing.B) {
//	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
//	svc := &mockEventService{}
//	cfg := &model.WorkerPoolEvent{
//		EventChainSize: 100000,
//		CountWorkers:   50,
//	}
//
//	pool := New(svc, logger, cfg)
//	ctx, cancel := context.WithCancel(context.Background())
//	pool.Start(ctx)
//	defer func() {
//		cancel()
//		pool.Stop()
//	}()
//
//	event := &model.Event{UserID: "user-parallel", Action: "test"}
//
//	b.ReportAllocs()
//	b.ResetTimer()
//
//	b.RunParallel(func(pb *testing.PB) {
//		for pb.Next() {
//			pool.AddEventRecord(event)
//		}
//	})
//}

package audit

import (
	"context"
	"time"

	//"fmt"
	"io"
	"log/slog"
	"runtime"
	"sync"
	"testing"

	"github.com/ArtShib/urlshortener/internal/model"
)

type MockEventService struct {
	wg sync.WaitGroup
}

func (m *MockEventService) SendAuditRecord(ctx context.Context, record *model.Event) error {
	_ = record
	m.wg.Done()
	return nil
}

func (m *MockEventService) Close() error {
	return nil
}

func BenchmarkWorkerPool_Base(b *testing.B) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	cfg := &model.WorkerPoolEvent{
		EventChainSize: b.N + 1000,
		CountWorkers:   10,
	}

	mockSvc := &MockEventService{}
	pool := New(mockSvc, logger, cfg)

	ctx, cancel := context.WithCancel(context.Background())

	pool.Start(ctx)

	time.Sleep(50 * time.Millisecond)

	for i := 1; i < int(cfg.CountWorkers); i++ {
		pool.addWorker(ctx)
	}

	mockSvc.wg.Add(b.N)
	dummyEvent := &model.Event{UserID: "bench_user_id", Action: "bench_action"}

	runtime.GC()
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		pool.AddEventRecord(dummyEvent)
	}

	mockSvc.wg.Wait()

	b.StopTimer()

	cancel()
	pool.Stop()
}
