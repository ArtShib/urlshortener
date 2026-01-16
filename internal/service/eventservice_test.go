package service

import (
	"context"
	"io"
	"log/slog"
	"runtime"
	"testing"

	"github.com/ArtShib/urlshortener/internal/model"
)

type mockRepo struct{}

func (m *mockRepo) SendAuditRecord(ctx context.Context, record *model.Event) error {
	_ = record
	return nil
}

func (m *mockRepo) Close() error {
	return nil
}

func BenchmarkEventService_SendAuditRecord(b *testing.B) {

	handler := slog.NewTextHandler(io.Discard, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})
	logger := slog.New(handler)

	repo := &mockRepo{}
	svc, err := NewEventService(repo, logger)
	if err != nil {
		b.Fatalf("failed to create service: %v", err)
	}

	record := &model.Event{
		UserID:      "id1234",
		Action:      "create_short_url",
		OriginalURL: "https://yandex.ru/test",
	}
	ctx := context.Background()

	runtime.GC()
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = svc.SendAuditRecord(ctx, record)
	}

}

func BenchmarkEventService_SendAuditRecord_Parallel(b *testing.B) {
	handler := slog.NewTextHandler(io.Discard, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})
	logger := slog.New(handler)
	repo := &mockRepo{}
	svc, _ := NewEventService(repo, logger)

	record := &model.Event{UserID: "1", Action: "test"}
	ctx := context.Background()

	runtime.GC()
	b.ReportAllocs()
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = svc.SendAuditRecord(ctx, record)
		}
	})
}
