package service

import (
	"context"
	"errors"
	"time"

	"github.com/ArtShib/urlshortener/internal/lib/shortener"
	"github.com/ArtShib/urlshortener/internal/model"
)

type URLRepository interface {
	Save(ctx context.Context, url *model.URL)  (*model.URL, error)
	Get(ctx context.Context, shortCode string) (*model.URL, error)
	Ping(ctx context.Context) error
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
	ctx, cansel := context.WithTimeout(context.Background(), 10 * time.Second)
	defer cansel()
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

	// if urlModel, err = s.repo.Save(ctx, urlModel); err != nil {
	// 	return urlModel.ShortURL, err
	// }
	
	urlModel, err = s.repo.Save(ctx, urlModel)
	return urlModel.ShortURL, err
}

func (s *URLService) GetID(shortCode string) (string, error) {
	ctx, cansel := context.WithTimeout(context.Background(), 10 * time.Second)
	defer cansel()

	if shortCode == "" {
		return "", errors.New("empty short code")
	}

	url, err := s.repo.Get(ctx, shortCode)
	if err != nil {
		return "", err
	}

	return url.OriginalURL, nil
}

func (s *URLService) ShortenJSON(url string) (*model.ResponseShortener, error) {
	ctx, cansel := context.WithTimeout(context.Background(), 10 * time.Second)
	defer cansel()

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

	// if urlModel, err = s.repo.Save(ctx, urlModel); err != nil {
	// 	return &model.ResponseShortener{
	// 		Result: urlModel.ShortURL,
	// 	}, err
	// }

	urlModel, err = s.repo.Save(ctx, urlModel)

	return &model.ResponseShortener{
		Result: urlModel.ShortURL,
	}, err

}

func (s *URLService) Ping(ctx context.Context) error {
	return s.repo.Ping(ctx)
}

func (s *URLService) ShortenJSONBatch(urls model.RequestShortenerBatchArray) (model.ResponseShortenerBatchArray, error) {
	ctx, cansel := context.WithTimeout(context.Background(), 10 * time.Second)
	defer cansel()

	var shortenerBatch model.ResponseShortenerBatchArray

	for _, url := range urls {
		uuid, err := shortener.GenerateUUID()
		if err != nil {
			return nil, err
		}
		shortURL := s.config.BaseURL 
		urlModel := &model.URL{
			UUID: uuid,
			ShortURL: shortener.GenerateShortURL(shortURL, uuid),
			OriginalURL: url.OriginalURL,
		}
		
		if _, err := s.repo.Save(ctx, urlModel); err != nil {
			return nil, err
		}
		shortenerBatch = append(shortenerBatch, model.ResponseShortenerBatch{
			CorrelationID: url.CorrelationID,
			ShortURL: urlModel.ShortURL,
		})	
	}
	return shortenerBatch, nil
}
