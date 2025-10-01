package handler

import (
	"context"

	"github.com/ArtShib/urlshortener/internal/model"
)

type URLService interface {
	Shorten(ctx context.Context, url string) (string, error)
	GetID(ctx context.Context, shortCode string) (*model.URL, error)
	ShortenJSON(ctx context.Context, url string) (*model.ResponseShortener, error)
	Ping(ctx context.Context) error
	ShortenJSONBatch(ctx context.Context, urls model.RequestShortenerBatchArray) (model.ResponseShortenerBatchArray, error)
	GetJSONBatch(ctx context.Context, userID string) (model.URLUserBatch, error)
}
