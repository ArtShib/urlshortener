package app

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/ArtShib/urlshortener/internal/config"
	"github.com/ArtShib/urlshortener/internal/httpserver"
	"github.com/ArtShib/urlshortener/internal/lib/auth"
	"github.com/ArtShib/urlshortener/internal/lib/shortener"
	"github.com/ArtShib/urlshortener/internal/repository"
	"github.com/ArtShib/urlshortener/internal/service"
	"github.com/ArtShib/urlshortener/internal/workerpool/audit"
	"github.com/ArtShib/urlshortener/internal/workerpool/requestdeletion"
)

// App структура слоя application
type App struct {
	Logger       *slog.Logger
	URLRepo      repository.URLRepository
	EventRepo    repository.EventRepository
	Server       *http.Server
	Config       *config.Config
	Auth         *auth.Service
	URLService   *service.URLService
	EventService *service.EventService
	WPoolDelete  *requestdeletion.DeletePool
	WPoolEvent   *audit.WorkerPoolEvent
}

// NewApp конструктор App
func NewApp(ctx context.Context, cfg *config.Config, repo *repository.URLRepository, eventRepo *repository.EventRepository, log *slog.Logger) *App {
	const op = "app.NewApp"
	app := &App{
		Config:    cfg,
		URLRepo:   *repo,
		EventRepo: *eventRepo,
		Logger:    log,
	}
	shortSvc := shortener.NewShortener()
	app.URLService = service.NewURLService(app.URLRepo, cfg.ShortService, shortSvc, app.Logger)
	app.WPoolDelete = requestdeletion.NewWorkerPool(app.URLService, app.Logger, cfg.Concurrency.WorkerPoolDelete)
	app.WPoolDelete.Start(ctx)
	app.Auth = auth.NewAuthService("048ff4ea240a9fdeac8f1422733e9f3b8b0291c969652225e25c5f0f9f8da654139c9e21")
	var err error
	app.EventService, err = service.NewEventService(app.EventRepo, app.Logger)
	if err != nil {
		log.Error(op, "error", fmt.Errorf("%s: %w", op, err))
	}
	app.WPoolEvent = audit.New(app.EventService, app.Logger, cfg.Concurrency.WorkerPoolEvent)
	if err == nil {
		app.WPoolEvent.Start(ctx)
	}
	app.Server = &http.Server{
		Addr:    app.Config.HTTPServer.ServerAddress,
		Handler: httpserver.NewRouter(app.URLService, app.Logger, app.Auth, app.WPoolDelete, app.WPoolEvent),
	}
	return app
}

// Run закпуск http сервера
func (a *App) Run() <-chan error {
	errCh := make(chan error, 1)

	go func() {
		if err := a.Server.ListenAndServe(); err != nil {
			a.Logger.Error(err.Error())
			errCh <- err
		}
		close(errCh)
	}()

	return errCh
}

// Stop остановка сервисов для реализации graceful shutdown
func (a *App) Stop(ctx context.Context) error {
	a.WPoolDelete.Stop()
	a.WPoolEvent.Stop()
	errRepo := a.URLRepo.Close()
	errServer := a.Server.Shutdown(ctx)

	if err := errors.Join(errRepo, errServer); err != nil {
		return err // fmt.Errorf("failed to stop app gracefully: %w", err)
	}

	return nil
}
