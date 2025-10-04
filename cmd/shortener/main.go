package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ArtShib/urlshortener/internal/app"
	"github.com/ArtShib/urlshortener/internal/config"
	"github.com/ArtShib/urlshortener/internal/repository"
)

func main() {
	cfg := config.MustLoadConfig()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var repo repository.URLRepository
	var err error
	if cfg.RepoConfig.DatabaseDSN != "" {
		repo, err = repository.NewRepository(ctx, "db", cfg.RepoConfig.DatabaseDSN)
	} else {
		repo, err = repository.NewRepository(ctx, "file", cfg.RepoConfig.FileStoragePath)
	}
	if err != nil && !os.IsNotExist(err) {
		log.Fatal(err.Error())
	}

	app := app.NewApp(cfg, &repo)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGQUIT)
	go app.Run()

	<-quit

	app.Stop(ctx)
}
