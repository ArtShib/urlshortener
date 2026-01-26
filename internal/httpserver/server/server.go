package server

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/ArtShib/urlshortener/internal/lib/loghelper"
	"github.com/ArtShib/urlshortener/internal/model"
)

// ServerHTTP структура для http сервера
type ServerHTTP struct {
	*http.Server
}

// New конструктор ServerHTTP
func New(serverAddress string, mux http.Handler) *ServerHTTP {
	return &ServerHTTP{
		&http.Server{
			Addr:    serverAddress,
			Handler: mux,
		},
	}
}

// Start запуск http сервера
func (s *ServerHTTP) Start(ctx context.Context, config *model.TLSConfig, log *slog.Logger) error {
	logHelper := loghelper.New(log, "ServerHTTP.Start")
	logHelper.LogDebug(ctx, "Server start")
	var err error
	if config.Enabled {
		logHelper.LogDebug(ctx, "TLS Enabled")
		err = s.ListenAndServeTLS(config.Cert, config.Key)
	} else {
		err = s.ListenAndServe()
	}
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		return logHelper.LogAndReturnError(ctx, "Error starting server", err)
	}
	logHelper.LogDebug(ctx, "Server stopped")
	return nil
}

// ShutDown остановка http сервера для реализации graceful shutdown
func (s *ServerHTTP) ShutDown(ctx context.Context, log *slog.Logger) error {
	logHelper := loghelper.New(log, "ServerHTTP.Stop")
	if err := s.Shutdown(ctx); err != nil {
		return logHelper.LogAndReturnError(ctx, "shutting down server failed", err)
	}
	logHelper.LogDebug(ctx, "Server shut down")
	return nil
}
