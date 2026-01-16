package service

import (
	"context"
	"io"
	"log/slog"
	"runtime"
	"testing"

	"github.com/ArtShib/urlshortener/internal/model"
)

type mockURLRepo struct{}

func (m *mockURLRepo) Save(ctx context.Context, url *model.URL) (*model.URL, error) {
	// Имитируем успешное сохранение и возврат
	return url, nil
}
func (m *mockURLRepo) Get(ctx context.Context, shortCode string) (*model.URL, error) {
	return &model.URL{OriginalURL: "http://yandex.ru", ShortURL: shortCode}, nil
}
func (m *mockURLRepo) Ping(ctx context.Context) error { return nil }
func (m *mockURLRepo) GetBatch(ctx context.Context, userID string) (model.URLUserBatch, error) {
	return nil, nil
}
func (m *mockURLRepo) DeleteBatch(ctx context.Context, deleteRequest model.URLUserRequestArray) error {
	return nil
}

type mockShortener struct{}

func (m *mockShortener) GenerateUUID() (string, error) {
	// Возвращаем статику, чтобы не бенчмаркать генератор UUID (если это не цель)
	return "a7v4M9PY", nil
}
func (m *mockShortener) GenerateShortURL(url string, uuid string) string {
	return url + "/" + uuid
}

func BenchmarkURLService_Shorten(b *testing.B) {

	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	repo := &mockURLRepo{}
	shortener := &mockShortener{}
	cfg := &model.ShortServiceConfig{BaseURL: "http://localhost:8080"}

	svc := NewURLService(repo, cfg, shortener, logger)
	ctx := context.Background()
	originalURL := "https://yandex.ru/test"

	runtime.GC()
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _ = svc.Shorten(ctx, originalURL)
	}
}

func BenchmarkURLService_ShortenJSONBatch(b *testing.B) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	repo := &mockURLRepo{}
	shortener := &mockShortener{}
	cfg := &model.ShortServiceConfig{BaseURL: "http://localhost:8080"}

	svc := NewURLService(repo, cfg, shortener, logger)
	ctx := context.Background()

	batchSize := 100
	batch := make(model.RequestShortenerBatchArray, batchSize)
	for i := 0; i < batchSize; i++ {
		batch[i] = model.RequestShortenerBatch{
			CorrelationID: "corr-id",
			OriginalURL:   "http://yandex.ru/test",
		}
	}

	runtime.GC()
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _ = svc.ShortenJSONBatch(ctx, batch)
	}
}
