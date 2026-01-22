package app

import (
	"context"
	"errors"
	"log/slog"

	"github.com/ArtShib/urlshortener/internal/config"
	"github.com/ArtShib/urlshortener/internal/httpserver"
	"github.com/ArtShib/urlshortener/internal/httpserver/server"
	"github.com/ArtShib/urlshortener/internal/lib/auth"
	"github.com/ArtShib/urlshortener/internal/lib/loghelper"
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
	Server       *server.ServerHTTP
	Config       *config.Config
	Auth         *auth.Service
	URLService   *service.URLService
	EventService *service.EventService
	WPoolDelete  *requestdeletion.DeletePool
	WPoolEvent   *audit.WorkerPoolEvent
}

// NewApp конструктор App
func NewApp(ctx context.Context, cfg *config.Config, repo *repository.URLRepository, eventRepo *repository.EventRepository, log *slog.Logger) *App {

	logHelper := loghelper.New(log, "app.NewApp")
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
		logHelper.LogError(ctx, "run service.NewEventService", err)
	}
	app.WPoolEvent = audit.New(app.EventService, app.Logger, cfg.Concurrency.WorkerPoolEvent)
	if err == nil {
		app.WPoolEvent.Start(ctx)
	}
	app.Server = server.New(app.Config.HTTPServer.ServerAddress, httpserver.NewRouter(app.URLService, app.Logger, app.Auth, app.WPoolDelete, app.WPoolEvent))
	return app
}

// Run закпуск http сервера
func (a *App) Run() <-chan error {

	errCh := make(chan error, 1)

	go func() {
		ctx := context.Background()
		err := a.Server.Start(ctx, a.Config.TLSConfig, a.Logger)
		if err != nil {
			errCh <- err
		}
		close(errCh)
	}()

	return errCh
}

// Stop остановка сервисов для реализации graceful shutdown
func (a *App) Stop(ctx context.Context) error {
	logHelper := loghelper.New(a.Logger, "app.Stop")
	errServer := a.Server.Shutdown(ctx)
	a.WPoolDelete.Stop()
	if a.EventRepo != nil {
		a.WPoolEvent.Stop()
	}
	errRepo := a.URLRepo.Close()
	if err := errors.Join(errRepo, errServer); err != nil {

		return logHelper.LogAndReturnError(ctx, "failed to stop app gracefully", err)
	}

	return nil
}
