package service

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/ArtShib/urlshortener/internal/model"
)

const (
	longOperationTimeout = 10 * time.Second
)

// URLRepository описывает интерфейс для работы с репозиторием данных urlshort
type URLRepository interface {
	Save(ctx context.Context, url *model.URL) (*model.URL, error)
	Get(ctx context.Context, shortCode string) (*model.URL, error)
	Ping(ctx context.Context) error
	GetBatch(ctx context.Context, userID string) (model.URLUserBatch, error)
	DeleteBatch(ctx context.Context, deleteRequest model.URLUserRequestArray) error
}

// Shortener описывает интерфейс для генерации uuid и ShortURL
type Shortener interface {
	GenerateUUID() (string, error)
	GenerateShortURL(url string, uuid string) string
}

// URLService структура URLService
type URLService struct {
	repo      URLRepository
	config    *model.ShortServiceConfig
	shortener Shortener
	logger    *slog.Logger
}

// NewURLService конструктор для URLService
func NewURLService(repo URLRepository, cfg *model.ShortServiceConfig, shortener Shortener, logger *slog.Logger) *URLService {
	urlService := &URLService{
		repo:      repo,
		config:    cfg,
		shortener: shortener,
		logger:    logger,
	}
	return urlService
}

// Shorten метод сервисного слоя, сокращения url
func (s *URLService) Shorten(ctx context.Context, url string) (string, error) {
	const op = "URLService.Shorten"
	log := s.logger.With(
		slog.String("op", op),
	)
	if url == "" {
		log.Error(op, "error", fmt.Errorf("empty URL"))
		return "", fmt.Errorf("%s: %w", op, fmt.Errorf("empty URL"))
	}

	uuid, err := s.shortener.GenerateUUID()
	if err != nil {
		log.Error(op, "error", err)
		return "", fmt.Errorf("%s: %w", op, err)
	}
	shortURL := s.config.BaseURL

	urlModel := &model.URL{
		UUID:        uuid,
		ShortURL:    s.shortener.GenerateShortURL(shortURL, uuid),
		OriginalURL: url,
	}

	userID, ok := ctx.Value(model.UserIDKey).(string)
	if ok && userID != "" {
		urlModel.UserID = userID
	}

	ctx, cancel := context.WithTimeout(ctx, longOperationTimeout)
	defer cancel()

	urlModel, err = s.repo.Save(ctx, urlModel)
	if err != nil {
		log.Error(op, "error", err)
		return "", fmt.Errorf("%s: %w", op, err)
	}
	return urlModel.ShortURL, err
}

// GetID метод сервисного слоя, получения оригинального url
func (s *URLService) GetID(ctx context.Context, shortCode string) (*model.URL, error) {
	const op = "URLService.GetID"
	log := s.logger.With(
		slog.String("op", op),
	)

	if shortCode == "" {
		log.Error(op, "error", fmt.Errorf("empty short code"))
		return nil, fmt.Errorf("%s: %w", op, fmt.Errorf("empty short code"))
	}

	url, err := s.repo.Get(ctx, shortCode)
	if err != nil {
		log.Error(op, "error", err)
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return url, nil
}

// ShortenJSON метод сервисного слоя, сокращение url. На вход подается json
func (s *URLService) ShortenJSON(ctx context.Context, url string) (*model.ResponseShortener, error) {
	const op = "URLService.ShortenJSON"
	log := s.logger.With(
		slog.String("op", op),
	)

	if url == "" {
		log.Error(op, "error", fmt.Errorf("empty URL"))
		return nil, fmt.Errorf("%s: %w", op, fmt.Errorf("empty URL"))
	}

	uuid, err := s.shortener.GenerateUUID()
	if err != nil {
		log.Error(op, "error", err)
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	shortURL := s.config.BaseURL
	urlModel := &model.URL{
		UUID:        uuid,
		ShortURL:    s.shortener.GenerateShortURL(shortURL, uuid),
		OriginalURL: url,
	}

	urlModel, err = s.repo.Save(ctx, urlModel)
	if err != nil {
		log.Error(op, "error", err)
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &model.ResponseShortener{
		Result: urlModel.ShortURL,
	}, err
}

// Ping метод сервисного слоя, проверка доступности репозитория
func (s *URLService) Ping(ctx context.Context) error {
	return s.repo.Ping(ctx)
}

// ShortenJSONBatch метод сервисного слоя сокращение url пачками
func (s *URLService) ShortenJSONBatch(ctx context.Context, urls model.RequestShortenerBatchArray) (model.ResponseShortenerBatchArray, error) {
	const op = "URLService.ShortenJSONBatch"
	log := s.logger.With(
		slog.String("op", op),
	)

	var shortenerBatch model.ResponseShortenerBatchArray

	for _, url := range urls {
		uuid, err := s.shortener.GenerateUUID()
		if err != nil {
			log.Error(op, "error", err)
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		shortURL := s.config.BaseURL
		urlModel := &model.URL{
			UUID:        uuid,
			ShortURL:    s.shortener.GenerateShortURL(shortURL, uuid),
			OriginalURL: url.OriginalURL,
		}

		if _, err := s.repo.Save(ctx, urlModel); err != nil {
			log.Error(op, "error", err)
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		shortenerBatch = append(shortenerBatch, model.ResponseShortenerBatch{
			CorrelationID: url.CorrelationID,
			ShortURL:      urlModel.ShortURL,
		})
	}
	return shortenerBatch, nil
}

// GetJSONBatch метод сервисного слоя, получения оригинального url по id пользователя
func (s *URLService) GetJSONBatch(ctx context.Context, userID string) (model.URLUserBatch, error) {
	const op = "URLService.GetJSONBatch"
	log := s.logger.With(
		slog.String("op", op),
	)
	UURLUserBatch, err := s.repo.GetBatch(ctx, userID)
	if err != nil {
		log.Error(op, "error", err)
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return UURLUserBatch, nil
}

// DeleteBatch метод сервисного слоя, удаления записи (соотношения uuid ( - оригинального url) из репозитория
func (s *URLService) DeleteBatch(ctx context.Context, batch model.URLUserRequestArray) error {
	const op = "URLService.DeleteBatch"
	log := s.logger.With(
		slog.String("op", op),
	)
	if err := s.repo.DeleteBatch(ctx, batch); err != nil {
		log.Error(op, "error", err)
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}
