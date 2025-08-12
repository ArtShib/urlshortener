package service

import (
	"errors"

	"github.com/ArtShib/urlshortener/internal/lib/shortener"
	"github.com/ArtShib/urlshortener/internal/model"
)

type UrlRepository interface {
	Store(url *model.URL) error
	FindByShortCode(shortCode string) (*model.URL, error)
}

type UrlService struct{
	repo UrlRepository
	config *model.HttpServerConfig
}

func NewUrlService(repo UrlRepository, cfg *model.HttpServerConfig) *UrlService {
	return &UrlService{
		repo: repo,
		config: cfg,
	}
}

func (s *UrlService) Shorten(url string) (string, error) {
	if (url == "") {
		return "", errors.New("empty URL")
	}
	shortCode := shortener.GenerateShortCode(url)
	shortUrl := s.config.Port 
	urlModel := &model.URL{
		LongUrl: url,
		ShortCode: shortCode,
		ShortUrl: shortener.GenerateShortUrl(shortUrl, shortCode),
	}

	if err := s.repo.Store(urlModel); err != nil {
		return "", err
	}
	return urlModel.ShortUrl, nil
}

func (s *UrlService) GetID(shortCode string) (string, error) {
	if (shortCode == "") {
		return "", errors.New("empty short code")
	}

	url, err := s.repo.FindByShortCode(shortCode)
	if err != nil {
		return "", err
	}

	return url.LongUrl, nil
}

