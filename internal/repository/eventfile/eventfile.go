package eventfile

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"sync"

	"github.com/ArtShib/urlshortener/internal/model"
)

// AuditEvent структура для работы с файлом для сохранения аудита
type AuditEvent struct {
	mu            sync.Mutex
	auditFilePath string
	logger        *slog.Logger
	encoder       *json.Encoder
	auditFile     *os.File
}

// New конструктор для AuditEvent
func New(auditFilePath string, log *slog.Logger) (*AuditEvent, error) {
	const op = "EventFile.New"
	logger := log.With(
		slog.String("op", op),
	)
	auditFile, err := os.OpenFile(auditFilePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		logger.Error(op, "error", err)
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &AuditEvent{
		auditFilePath: auditFilePath,
		auditFile:     auditFile,
		encoder:       json.NewEncoder(auditFile),
		logger:        log,
	}, nil
}

// SendAuditRecord запись в файл записи аудита
func (a *AuditEvent) SendAuditRecord(ctx context.Context, record *model.Event) error {
	const op = "EventFile.SendAuditRecord"
	log := a.logger.With(
		slog.String("op", op),
	)

	if err := ctx.Err(); err != nil {
		log.Error(op, "error", err)
		return fmt.Errorf("%s: %w", op, err)
	}
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.encoder.Encode(record)
}

// Close закрытие файла
func (a *AuditEvent) Close() error {
	const op = "EventFile.Close"
	log := a.logger.With(
		slog.String("op", op),
	)

	if err := a.auditFile.Sync(); err != nil {
		log.Error(op, "error", err)
		return fmt.Errorf("%s: %w", op, err)
	}
	return a.auditFile.Close()
}
