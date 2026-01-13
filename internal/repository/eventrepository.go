package repository

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/ArtShib/urlshortener/internal/httpclient"
	"github.com/ArtShib/urlshortener/internal/model"
	"github.com/ArtShib/urlshortener/internal/repository/eventfile"
)

// EventRepository описывает интерфейс сохранения аудита
type EventRepository interface {
	SendAuditRecord(ctx context.Context, record *model.Event) error
	Close() error
}

// NewEventRepository конструктор создания репозитория под аудит
func NewEventRepository(auditFilePath string, auditUrl string, log *slog.Logger) (EventRepository, error) {
	const op = "repository.NewEventRepository"
	logger := log.With(
		slog.String("op", op),
	)
	if auditFilePath != "" {
		eventRepository, err := eventfile.New(auditFilePath, log)
		if err != nil {
			logger.Error(op, "error", err)
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		return eventRepository, nil
	}
	if auditUrl != "" {
		eventRepository := httpclient.New(log, auditUrl)
		return eventRepository, nil
	}
	return nil, fmt.Errorf("%s: %w", op, fmt.Errorf("audit file and url is empty"))
}
