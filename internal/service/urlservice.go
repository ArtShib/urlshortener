package service

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/ArtShib/urlshortener/internal/lib/loghelper"
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
	Stats(ctx context.Context) (*model.Stats, error)
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
	log := loghelper.New(s.logger, "URLService.Shorten")
	if url == "" {
		return "", log.LogAndReturnError(ctx, "empty URL", fmt.Errorf("empty URL"))
	}

	uuid, err := s.shortener.GenerateUUID()
	if err != nil {
		//log.Error(op, "error", err)
		//return "", fmt.Errorf("%s: %w", op, err)
		return "", log.LogAndReturnError(ctx, "error server.GenerateUUID", err)
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
		log.LogError(ctx, "error func repo.Save", err, slog.String("OriginalURL", urlModel.OriginalURL))

	}
	return urlModel.ShortURL, err
}

// GetID метод сервисного слоя, получения оригинального url
func (s *URLService) GetID(ctx context.Context, shortCode string) (*model.URL, error) {

	log := loghelper.New(s.logger, "URLService.GetID")

	if shortCode == "" {
		//log.Error(op, "error", fmt.Errorf("empty short code"))
		//return nil, fmt.Errorf("%s: %w", op, fmt.Errorf("empty short code"))
		return nil, log.LogAndReturnError(ctx, "empty short code", fmt.Errorf("empty short code"))
	}

	url, err := s.repo.Get(ctx, shortCode)
	if err != nil {
		//log.Error(op, "error", err)
		//return nil, fmt.Errorf("%s: %w", op, err)
		return nil, log.LogAndReturnError(ctx, "error repo.Get", err)
	}

	return url, nil
}

// ShortenJSON метод сервисного слоя, сокращение url. На вход подается json
func (s *URLService) ShortenJSON(ctx context.Context, url string) (*model.ResponseShortener, error) {
	log := loghelper.New(s.logger, "URLService.ShortenJSON")

	if url == "" {
		return nil, log.LogAndReturnError(ctx, "empty URL", fmt.Errorf("empty URL"))
	}

	uuid, err := s.shortener.GenerateUUID()
	if err != nil {
		return nil, log.LogAndReturnError(ctx, "error server.GenerateUUID", err)
	}

	shortURL := s.config.BaseURL
	urlModel := &model.URL{
		UUID:        uuid,
		ShortURL:    s.shortener.GenerateShortURL(shortURL, uuid),
		OriginalURL: url,
	}

	urlModel, err = s.repo.Save(ctx, urlModel)
	if err != nil {
		return nil, log.LogAndReturnError(ctx, "error repo.Save", err)
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
	log := loghelper.New(s.logger, "URLService.ShortenJSONBatch")

	var shortenerBatch model.ResponseShortenerBatchArray

	for _, url := range urls {
		uuid, err := s.shortener.GenerateUUID()
		if err != nil {
			return nil, log.LogAndReturnError(ctx, "error server.GenerateUUID", err)
		}
		shortURL := s.config.BaseURL
		urlModel := &model.URL{
			UUID:        uuid,
			ShortURL:    s.shortener.GenerateShortURL(shortURL, uuid),
			OriginalURL: url.OriginalURL,
		}

		if _, err := s.repo.Save(ctx, urlModel); err != nil {
			return nil, log.LogAndReturnError(ctx, "error repo.Save", err)
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
	log := loghelper.New(s.logger, "URLService.GetJSONBatch")

	UURLUserBatch, err := s.repo.GetBatch(ctx, userID)
	if err != nil {
		return nil, log.LogAndReturnError(ctx, "error repo.GetBatch", err)
	}
	return UURLUserBatch, nil
}

// DeleteBatch метод сервисного слоя, удаления записи (соотношения uuid ( - оригинального url) из репозитория
func (s *URLService) DeleteBatch(ctx context.Context, batch model.URLUserRequestArray) error {
	log := loghelper.New(s.logger, "URLService.DeleteBatch")

	if err := s.repo.DeleteBatch(ctx, batch); err != nil {
		return log.LogAndReturnError(ctx, "error repo.DeleteBatch", err)
	}
	return nil
}

// Stats метод для получения статистики
func (s *URLService) Stats(ctx context.Context) (*model.Stats, error) {

	log := loghelper.New(s.logger, "URLService.Stats")

	stats, err := s.repo.Stats(ctx)
	if err != nil {
		//log.Error(op, "error", err)
		//return nil, fmt.Errorf("%s: %w", op, err)
		return nil, log.LogAndReturnError(ctx, "error repo.Stats", err)
	}

	return stats, nil
}
