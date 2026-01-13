package repository

import (
	"context"
	"log"
	"log/slog"

	"github.com/ArtShib/urlshortener/internal/model"
	"github.com/ArtShib/urlshortener/internal/repository/memory"
	"github.com/ArtShib/urlshortener/internal/repository/postgres"
)

// URLRepository описывает интерфейс для работы с репозиторием данных urlshort
type URLRepository interface {
	Save(ctx context.Context, url *model.URL) (*model.URL, error)
	Get(ctx context.Context, uuid string) (*model.URL, error)
	Close() error
	Ping(context.Context) error
	GetBatch(ctx context.Context, userID string) (model.URLUserBatch, error)
	DeleteBatch(ctx context.Context, deleteRequest model.URLUserRequestArray) error
}

// NewURLRepository конструктор создания репозитория
func NewURLRepository(ctx context.Context, repoType string, dsnORpath string, logger *slog.Logger) (URLRepository, error) {
	switch repoType {
	case "db":
		return postgres.NewPostgresRepository(ctx, dsnORpath, logger)
	case "file":
		return memory.NewMemoryRepository(ctx, dsnORpath)
	}
	log.Fatal("not loaded repo")
	return nil, nil
}
