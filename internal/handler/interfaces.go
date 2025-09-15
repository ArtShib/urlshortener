package handler

import "github.com/ArtShib/urlshortener/internal/model"


type URLService interface {
	Shorten(url string) (string, error)
	GetID(shortCode string) (string, error)
	ShortenJSON(url string) (*model.ResponseShortener, error)
}
