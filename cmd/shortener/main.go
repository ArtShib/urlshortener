package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

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
	defer repo.SavingRepository(cfg.RepoConfig.FileStoragePath)

	svc := service.NewURLService(repo, cfg.ShortService)
	router := handler.NewRouter(svc, logg)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Fatal(http.ListenAndServe(cfg.HTTPServer.ServerAddress, router))
	}()
	
	<-sigChan
	
}
