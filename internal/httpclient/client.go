package httpclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"

	"github.com/ArtShib/urlshortener/internal/model"
)

// Client структура http клиена
type Client struct {
	log        *slog.Logger
	httpClient *http.Client
	auditUrl   string
}

// New конструктор Client
func New(log *slog.Logger, auditUrl string) *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: time.Second * 10,
			Transport: &http.Transport{
				MaxIdleConnsPerHost: 100,
				MaxIdleConns:        100,
				//MaxConnsPerHost:     10,
				IdleConnTimeout: time.Second * 90,
			},
		},
		log:      log,
		auditUrl: auditUrl,
	}
}

// SendAuditRecord отправка записи аудита на удаленный http сервер
func (c *Client) SendAuditRecord(ctx context.Context, record *model.Event) error {
	const op = "Client.SendEventRecord"
	log := c.log.With(
		slog.String("op", op),
	)

	log.Info("request urlConnect")

	body, err := json.Marshal(record)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	request, err := http.NewRequestWithContext(ctx, http.MethodPost, c.auditUrl, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	request.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(request)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	defer func() {
		_, _ = io.Copy(io.Discard, resp.Body)
		if err := resp.Body.Close(); err != nil {
			c.log.Error(op, "error", err)
		}
	}()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("%s: invalid status code: %d", op, resp.StatusCode)
	}

	return nil
}

// Close закрытие соединения
func (c *Client) Close() error {
	c.httpClient.CloseIdleConnections()
	return nil
}
