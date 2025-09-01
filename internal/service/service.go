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
	config *model.ShortServiceConfig
}

func NewURLService(repo URLRepository, cfg *model.ShortServiceConfig) *URLService {
	return &URLService{
		repo: repo,
		config: cfg,
	}
}

func (s *URLService) Shorten(url string) (string, error) {
	if url == "" {
		return "", errors.New("empty URL")
	}
	uuid, err := shortener.GenerateUUID()
	if err != nil {
		return "", err
	}
	shortURL := s.config.BaseURL 
	urlModel := &model.URL{
		UUID: uuid,
		ShortURL: shortener.GenerateShortURL(shortURL, uuid),
		OriginalURL: url,
	}

	if err := s.repo.Store(urlModel); err != nil {
		return "", err
	}
	return urlModel.ShortURL, nil
}

func (s *URLService) GetID(shortCode string) (string, error) {
	if shortCode == "" {
		return "", errors.New("empty short code")
	}

	url, err := s.repo.FindByShortCode(shortCode)
	if err != nil {
		return "", err
	}

	return url.OriginalURL, nil
}

func (s *URLService) ShortenJson(url string) (*model.ResponseShortener, error) {
	if url == "" {
		return nil, errors.New("empty URL")
	}
	uuid, err := shortener.GenerateUUID()
	if err != nil {
		return nil, err
	}
	shortURL := s.config.BaseURL 
	urlModel := &model.URL{
		UUID: uuid,
		ShortURL: shortener.GenerateShortURL(shortURL, uuid),
		OriginalURL: url,
	}

	if err := s.repo.Store(urlModel); err != nil {
		return nil, err
	}

	return &model.ResponseShortener{
		Result: urlModel.ShortURL,
	}, nil
}

