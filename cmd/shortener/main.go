package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ArtShib/urlshortener/internal/app"
	"github.com/ArtShib/urlshortener/internal/config"
	myLogger "github.com/ArtShib/urlshortener/internal/lib/logger"
	"github.com/ArtShib/urlshortener/internal/repository"
)

func main() {
	const op = "main"
	var err error
	logger := myLogger.NewLogger()
	cfg, err := config.MustLoadConfig()
	if err != nil {
		logger.Error(op, "error", err)
	}
	initCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var urlRepo repository.URLRepository

	if cfg.RepoConfig.DatabaseDSN != "" {
		urlRepo, err = repository.NewURLRepository(initCtx, "db", cfg.RepoConfig.DatabaseDSN, logger)
	} else {
		urlRepo, err = repository.NewURLRepository(initCtx, "file", cfg.RepoConfig.FileStoragePath, logger)
	}
	if err != nil && !os.IsNotExist(err) {
		logger.Error(op, "error", err)
		os.Exit(1)
	}

	eventRepo, err := repository.NewEventRepository(cfg.AuditConfig.AuditFile, cfg.AuditConfig.AuditUrl, logger)
	if err != nil {
		logger.Error(op, "error", err)
	}
	application := app.NewApp(context.Background(), cfg, &urlRepo, &eventRepo, logger)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGQUIT)
	errCh := application.Run()

	select {
	case err := <-errCh:
		logger.Error(op, "error", err)
		os.Exit(1)
	case <-quit:
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer shutdownCancel()

		if err := application.Stop(shutdownCtx); err != nil {
			logger.Error(op, "shutdown error", err)
		}
	}
}
