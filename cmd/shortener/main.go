package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ArtShib/urlshortener/internal/app"
	"github.com/ArtShib/urlshortener/internal/config"
	myLogger "github.com/ArtShib/urlshortener/internal/lib/logger"
	"github.com/ArtShib/urlshortener/internal/lib/loghelper"
	"github.com/ArtShib/urlshortener/internal/repository"
)

var (
	buildVersion string = "N/A"
	buildDate    string = "N/A"
	buildCommit  string = "N/A"
)

func main() {

	fmt.Printf("Build version: %s\nBuild date: %s\nBuild commit: %s\n", buildVersion, buildDate, buildCommit)

	//const op = "main"
	var err error
	logger := myLogger.NewLogger()
	logHelper := loghelper.New(logger, "main")

	ctx := context.Background()

	cfg, err := config.MustLoadConfig()
	if err != nil {
		logHelper.LogError(ctx, "run MustLoadConfig", err)
	}

	initCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	var urlRepo repository.URLRepository

	if cfg.RepoConfig.DatabaseDSN != "" {
		urlRepo, err = repository.NewURLRepository(initCtx, "db", cfg.RepoConfig.DatabaseDSN, logger)
	} else {
		urlRepo, err = repository.NewURLRepository(initCtx, "file", cfg.RepoConfig.FileStoragePath, logger)
	}
	if err != nil && !os.IsNotExist(err) {
		logHelper.LogErrorAndExit(initCtx, "init repository", err)
	}

	eventRepo, err := repository.NewEventRepository(cfg.AuditConfig.AuditFile, cfg.AuditConfig.AuditURL, logger)
	if err != nil {
		logHelper.LogError(ctx, "init event repository", err)
	}

	application := app.NewApp(ctx, cfg, &urlRepo, &eventRepo, logger)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGQUIT)
	errCh := application.Run()

	select {
	case err := <-errCh:
		logHelper.LogErrorAndExit(initCtx, "run application", err)
	case <-quit:
		shutdownCtx, shutdownCancel := context.WithTimeout(ctx, 10*time.Second)
		defer shutdownCancel()

		if err := application.Stop(shutdownCtx); err != nil {
			logHelper.LogError(shutdownCtx, "application shutdown error", err)
		}
	}
}
