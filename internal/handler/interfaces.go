package handler

import (
	"context"

	"github.com/ArtShib/urlshortener/internal/model"
)


type URLService interface {
	Shorten(url string) (string, error)
	GetID(shortCode string) (string, error)
	ShortenJSON(url string) (*model.ResponseShortener, error)
	Ping(ctx context.Context) error
	ShortenJSONBatch(urls model.RequestShortenerBatchArray) (model.ResponseShortenerBatchArray, error)
}
