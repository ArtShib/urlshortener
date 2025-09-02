package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ArtShib/urlshortener/internal/config"
	"github.com/ArtShib/urlshortener/internal/handler"
	"github.com/ArtShib/urlshortener/internal/lib/logger"
	"github.com/ArtShib/urlshortener/internal/repository"
	"github.com/ArtShib/urlshortener/internal/service"
)

func main() {
	cfg := config.MustLoadConfig()
	
	logg := logger.NewLogger()
	repo, _ := repository.NewRepository(cfg.RepoConfig.FileStoragePath)

	svc := service.NewURLService(repo, cfg.ShortService)
	router := handler.NewRouter(svc, logg)
	
	server := &http.Server{
		Addr: cfg.HTTPServer.ServerAddress,
		Handler: router,
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGQUIT)

	
	go func() {
		log.Fatal(server.ListenAndServe())
	}()
	
	<-quit
	
	ctx, cansel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cansel()
	
	repo.SavingRepository(cfg.RepoConfig.FileStoragePath)
	
	server.Shutdown(ctx)
}
