package service

import (
	"context"
	"fmt"
	"log/slog"

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
	const op = "EventService.NewEventService"
	log := logger.With(
		slog.String("op", op),
	)
	if eventRepository == nil {
		log.Error(op, "error", fmt.Errorf("audit file and url is empty"))
		return &EventService{logger: log}, fmt.Errorf("%s: %w", op, fmt.Errorf("audit file and url is empty"))
	}
	return &EventService{
		eventRepository: eventRepository,
		logger:          logger,
	}, nil
}

// SendAuditRecord сохранение сообщения аудита
func (s *EventService) SendAuditRecord(ctx context.Context, record *model.Event) error {
	const op = "EventService.SendAuditRecord"
	log := s.logger.With(
		slog.String("op", op),
	)
	log.Debug("start EventService.SendAuditRecord")
	return s.eventRepository.SendAuditRecord(ctx, record)
}

// Close закрытие репозитория куда сохраняются сообщения аудита
func (s *EventService) Close() error {
	const op = "EventService.Close"
	log := s.logger.With(
		slog.String("op", op),
	)
	log.Debug("start EventService.Close")
	if s.eventRepository != nil {
		if err := s.eventRepository.Close(); err != nil {
			s.logger.Error(op, "error", err)
			return fmt.Errorf("%s: %w", op, err)
		}
	}
	return nil
}
