package requestdeletion

import (
	"context"
	"io"
	"log/slog"
	"runtime"
	"sync"
	"testing"

	"github.com/ArtShib/urlshortener/internal/model"
)

type MockURLService struct {
	wg sync.WaitGroup
}

func (m *MockURLService) DeleteBatch(ctx context.Context, batch model.URLUserRequestArray) error {
	//for range batch {
	//	m.wg.Done()
	//}
	m.wg.Add(-len(batch))
	return nil
}

func BenchmarkWorkerPool_Delete(b *testing.B) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	cfg := &model.WorkerPoolDelete{
		CountWorkers:   10,
		InputChainSize: 10000,
		BufferSize:     100,
		BatchSize:      100,
	}

	mockSvc := &MockURLService{}
	pool := NewWorkerPool(mockSvc, logger, cfg)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	pool.Start(ctx)

	req := model.DeleteRequest{
		UserID: "user-123",
		UUIDs: []string{
			"uuid-1", "uuid-2", "uuid-3", "uuid-4", "uuid-5",
			"uuid-6", "uuid-7", "uuid-8", "uuid-9", "uuid-10",
		},
	}

	mockSvc.wg.Add(b.N * 10)

	runtime.GC()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		pool.AddRequest(req)
	}

	mockSvc.wg.Wait()

	b.StopTimer()
	//cancel()
	pool.Stop()
}
