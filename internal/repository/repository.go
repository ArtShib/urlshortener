package repository

import (
	"context"
	"log"

	"github.com/ArtShib/urlshortener/internal/model"
	"github.com/ArtShib/urlshortener/internal/repository/memory"
	"github.com/ArtShib/urlshortener/internal/repository/postgres"
)

type URLRepository interface {
	Save(ctx context.Context, url *model.URL) (*model.URL, error)
	Get(ctx context.Context, uuid string) (*model.URL, error)
	Close() error
	Ping(context.Context) error
	GetBatch(ctx context.Context, userId string) (*model.URLUserBatch, error)
}

func NewRepository(ctx context.Context, repoType string, dsnORpath string) (URLRepository, error) {
	switch repoType {
	case "db":
		return postgres.NewPostgresRepository(ctx, dsnORpath)
	case "file":
		return memory.NewMemoryRepository(ctx, dsnORpath)
	}
	log.Fatal("not loaded repo")
	return nil, nil
}
