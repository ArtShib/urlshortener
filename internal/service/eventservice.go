package service

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/ArtShib/urlshortener/internal/lib/loghelper"
	"github.com/ArtShib/urlshortener/internal/model"
)

// EventRepository описывает интерфейс сохранения аудита
type EventRepository interface {
	SendAuditRecord(ctx context.Context, record *model.Event) error
	Close() error
}

// EventService структура сервиса аудита
type EventService struct {
	eventRepository EventRepository
	logger          *slog.Logger
}

// NewEventService констструктор сервиса аудита
func NewEventService(eventRepository EventRepository, logger *slog.Logger) (*EventService, error) {
	log := loghelper.New(logger, "EventService.NewEventService")

	if eventRepository == nil {
		return nil, log.LogAndReturnError(context.Background(), "EventService.NewEventService", fmt.Errorf("audit file and url is empty"))
	}
	return &EventService{
		eventRepository: eventRepository,
		logger:          logger,
	}, nil
}

// SendAuditRecord сохранение сообщения аудита
func (s *EventService) SendAuditRecord(ctx context.Context, record *model.Event) error {
	log := loghelper.New(s.logger, "EventService.SendAuditRecord")

	log.LogDebug(ctx, "start EventService.SendAuditRecord")
	return s.eventRepository.SendAuditRecord(ctx, record)
}

// Close закрытие репозитория куда сохраняются сообщения аудита
func (s *EventService) Close() error {
	ctx := context.Background()
	log := loghelper.New(s.logger, "EventService.Close")
	log.LogDebug(ctx, "start EventService.Close")
	if s.eventRepository != nil {
		if err := s.eventRepository.Close(); err != nil {
			return log.LogAndReturnError(ctx, "EventService.Close", err)
		}
	}
	return nil
}
