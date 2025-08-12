package service

import (
	"errors"

	"github.com/ArtShib/urlshortener/internal/lib/shortener"
	"github.com/ArtShib/urlshortener/internal/model"
)

type URLRepository interface {
	Store(url *model.URL) error
	FindByShortCode(shortCode string) (*model.URL, error)
}

type URLService struct{
	repo URLRepository
	config *model.HTTPServerConfig
}

func NewURLService(repo URLRepository, cfg *model.HTTPServerConfig) *URLService {
	return &URLService{
		repo: repo,
		config: cfg,
	}
}

func (s *URLService) Shorten(url string) (string, error) {
	if (url == "") {
		return "", errors.New("empty URL")
	}
	shortCode := shortener.GenerateShortCode(url)
	shortURL := s.config.Port 
	urlModel := &model.URL{
		LongURL: url,
		ShortCode: shortCode,
		ShortURL: shortener.GenerateShortURL(shortURL, shortCode),
	}

	if err := s.repo.Store(urlModel); err != nil {
		return "", err
	}
	return urlModel.ShortURL, nil
}

func (s *URLService) GetID(shortCode string) (string, error) {
	if (shortCode == "") {
		return "", errors.New("empty short code")
	}

	url, err := s.repo.FindByShortCode(shortCode)
	if err != nil {
		return "", err
	}

	return url.LongURL, nil
}
