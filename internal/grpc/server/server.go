package server

import (
	"context"

	"github.com/ArtShib/urlshortener/gen/go/shortener"
	"github.com/ArtShib/urlshortener/internal/model"
)

// URLService интерфейс сервиса для получения оригинального url.
type URLService interface {
	Shorten(ctx context.Context, url string) (string, error)
	GetID(ctx context.Context, shortCode string) (*model.URL, error)
	GetJSONBatch(ctx context.Context, userID string) (model.URLUserBatch, error)
}

type serverAPI struct {
	shortener.UnimplementedShortenerServiceServer
	service URLService
}
