package app

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

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
func NewApp(cfg *config.Config) (*App, error) {
	app := &App{
		Config: cfg,
	}
	app.Logger = logger.NewLogger()
	var err error
	app.Repository, err = repository.NewRepository(cfg.RepoConfig.FileStoragePath)
	if err != nil && !os.IsNotExist(err) {
		return nil, err 
	} 
	svc := service.NewURLService(app.Repository, cfg.ShortService)
	app.Server =  &http.Server{
				Addr: app.Config.HTTPServer.ServerAddress,
				Handler: handler.NewRouter(svc, app.Logger),
			} 
	return app, nil
}

func (a *App) Run() {
	go func() {
		if err := a.Server.ListenAndServe(); err != nil {
			fmt.Println(err)
			os.Exit(2)
		}
		
	}()
}

func (a *App) Stop(){
	ctx, cansel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cansel()
	a.Repository.SavingRepository(a.Config.RepoConfig.FileStoragePath)
	
	a.Server.Shutdown(ctx)
}
