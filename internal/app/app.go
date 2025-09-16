package app

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/ArtShib/urlshortener/internal/config"
	"github.com/ArtShib/urlshortener/internal/handler"
	"github.com/ArtShib/urlshortener/internal/lib/logger"
	"github.com/ArtShib/urlshortener/internal/repository"
	"github.com/ArtShib/urlshortener/internal/service"
)


type App struct {
	Logger *slog.Logger
	Repository repository.URLRepository
	Server	*http.Server
	Config *config.Config
}
func NewApp(cfg *config.Config, repo *repository.URLRepository) *App {
	app := &App{
		Config: cfg,
		Repository: *repo,
	}
	app.Logger = logger.NewLogger()
	svc := service.NewURLService(app.Repository, cfg.ShortService)
	app.Server =  &http.Server{
				Addr: app.Config.HTTPServer.ServerAddress,
				Handler: handler.NewRouter(svc, app.Logger),
			} 
	return app
}

func (a *App) Run() {
	go func() {
		if err := a.Server.ListenAndServe(); err != nil {
			a.Logger.Error(err.Error())
		}
		
	}()
}

func (a *App) Stop(ctx context.Context){
	a.Repository.Close()
	a.Server.Shutdown(ctx)
}
