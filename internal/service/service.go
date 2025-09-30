package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/ArtShib/urlshortener/internal/lib/shortener"
	"github.com/ArtShib/urlshortener/internal/model"
)

const (
	defaultRequestTimeout = 3 * time.Second
	longOperationTimeout  = 10 * time.Second
)

type URLRepository interface {
	Save(ctx context.Context, url *model.URL) (*model.URL, error)
	Get(ctx context.Context, shortCode string) (*model.URL, error)
	Ping(ctx context.Context) error
	GetBatch(ctx context.Context, userId string) (*model.URLUserBatch, error)
}

type URLService struct {
	repo   URLRepository
	config *model.ShortServiceConfig
}

func NewURLService(repo URLRepository, cfg *model.ShortServiceConfig) *URLService {
	return &URLService{
		repo:   repo,
		config: cfg,
	}
}

func (s *URLService) Shorten(ctx context.Context, url string) (string, error) {

	if url == "" {
		return "", errors.New("empty URL")
	}
	uuid, err := shortener.GenerateUUID()
	if err != nil {
		return "", err
	}
	shortURL := s.config.BaseURL

	urlModel := &model.URL{
		UUID:        uuid,
		ShortURL:    shortener.GenerateShortURL(shortURL, uuid),
		OriginalURL: url,
	}

	if userID, ok := ctx.Value("UserID").(string); ok && userID != "" {
		urlModel.UserID = userID
	} else {
		fmt.Errorf("No UserID in context")
	}

	ctx, cancel := context.WithTimeout(ctx, longOperationTimeout)
	defer cancel()

	urlModel, err = s.repo.Save(ctx, urlModel)
	return urlModel.ShortURL, err
}

func (s *URLService) GetID(ctx context.Context, shortCode string) (string, error) {

	if shortCode == "" {
		return "", errors.New("empty short code")
	}

	url, err := s.repo.Get(ctx, shortCode)
	if err != nil {
		return "", err
	}

	return url.OriginalURL, nil
}

func (s *URLService) ShortenJSON(ctx context.Context, url string) (*model.ResponseShortener, error) {

	if url == "" {
		return nil, errors.New("empty URL")
	}
	uuid, err := shortener.GenerateUUID()
	if err != nil {
		return nil, err
	}
	shortURL := s.config.BaseURL
	urlModel := &model.URL{
		UUID:        uuid,
		ShortURL:    shortener.GenerateShortURL(shortURL, uuid),
		OriginalURL: url,
	}

	urlModel, err = s.repo.Save(ctx, urlModel)

	return &model.ResponseShortener{
		Result: urlModel.ShortURL,
	}, err

}

func (s *URLService) Ping(ctx context.Context) error {
	return s.repo.Ping(ctx)
}

func (s *URLService) ShortenJSONBatch(ctx context.Context, urls model.RequestShortenerBatchArray) (model.ResponseShortenerBatchArray, error) {

	var shortenerBatch model.ResponseShortenerBatchArray

	for _, url := range urls {
		uuid, err := shortener.GenerateUUID()
		if err != nil {
			return nil, err
		}
		shortURL := s.config.BaseURL
		urlModel := &model.URL{
			UUID:        uuid,
			ShortURL:    shortener.GenerateShortURL(shortURL, uuid),
			OriginalURL: url.OriginalURL,
		}

		if _, err := s.repo.Save(ctx, urlModel); err != nil {
			return nil, err
		}
		shortenerBatch = append(shortenerBatch, model.ResponseShortenerBatch{
			CorrelationID: url.CorrelationID,
			ShortURL:      urlModel.ShortURL,
		})
	}
	return shortenerBatch, nil
}

func (s *URLService) GetJSONBatch(ctx context.Context, userID string) (*model.URLUserBatch, error) {

	UURLUserBatch, err := s.repo.GetBatch(ctx, userID)
	if err != nil {
		return nil, err
	}
	return UURLUserBatch, nil
}
