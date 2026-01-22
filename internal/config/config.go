package config

import (
	"context"
	"flag"
	"log/slog"
	"os"

	"github.com/ArtShib/urlshortener/internal/lib/loghelper"
	"github.com/ArtShib/urlshortener/internal/model"
	"github.com/caarlos0/env"
	"github.com/joho/godotenv"
)

// Config структура конфига
type Config struct {
	HTTPServer   *model.HTTPServerConfig
	ShortService *model.ShortServiceConfig
	RepoConfig   *model.RepositoryConfig
	Concurrency  *model.Concurrency
	AuditConfig  *model.AuditConfig
	TLSConfig    *model.TLSConfig
	ConfigFile   *model.ConfigFile
	logger       *slog.Logger
}

// LoadConfigEnv загрузка данных в конфиг из env
func (c *Config) LoadConfigEnv() error {
	if err := godotenv.Load(); err != nil {
		return err
	}
	if err := env.Parse(c.HTTPServer); err != nil {
		return err
	}
	if err := env.Parse(c.ShortService); err != nil {
		return err
	}
	if err := env.Parse(c.RepoConfig); err != nil {
		return err
	}
	if err := env.Parse(c.AuditConfig); err != nil {
		return err
	}
	if err := env.Parse(c.TLSConfig); err != nil {
		return err
	}
	if err := env.Parse(c.ConfigFile); err != nil {
		return err
	}
	return nil
}

// LoadConfigFlag загрузка данных в конфиг из cmd
func (c *Config) LoadConfigFlag() {
	if c.HTTPServer.ServerAddress == "" {
		flag.StringVar(&c.HTTPServer.ServerAddress, "a", ":8080", "HTTP server startup address")
	}
	if c.ShortService.BaseURL == "" {
		flag.StringVar(&c.ShortService.BaseURL, "b", "http://localhost:8080", "Address of the resulting shortened URL")
	}
	if c.RepoConfig.FileStoragePath == "" {
		flag.StringVar(&c.RepoConfig.FileStoragePath, "f", "", "File storage path")
	}
	if c.RepoConfig.DatabaseDSN == "" {
		flag.StringVar(&c.RepoConfig.DatabaseDSN, "d", "", "DataBase connection string")
	}
	if c.AuditConfig.AuditFile == "" {
		flag.StringVar(&c.AuditConfig.AuditFile, "AUDIT_FILE", "/home/artem/GolandProjects/urlshortener/storage/audit.json", "Audit file path")
	}
	if c.AuditConfig.AuditURL == "" {
		flag.StringVar(&c.AuditConfig.AuditURL, "AUDIT_URL", "", "URL to audit")
	}
	if c.TLSConfig.Enabled == false {
		flag.BoolVar(&c.TLSConfig.Enabled, "s", false, "Enable TLS")
	}
	if c.ConfigFile.Path == "" {
		flag.StringVar(&c.ConfigFile.Path, "c", "", "Configuration file path")
	}
	if c.ConfigFile.Path == "" {
		flag.StringVar(&c.ConfigFile.Path, "config", "", "Configuration file path")
	}

	flag.Parse()
}

// MustLoadConfig конструктор Config
func MustLoadConfig(ctx context.Context, logger *slog.Logger) (*Config, error) {
	logHelper := loghelper.New(logger, "config.MustLoadConfig")
	var err error
	cfg := Config{
		HTTPServer:   &model.HTTPServerConfig{},
		ShortService: &model.ShortServiceConfig{},
		RepoConfig: &model.RepositoryConfig{
			FileStoragePath: os.Getenv("FILE_STORAGE_PATH"),
			DatabaseDSN:     os.Getenv("DATABASE_DSN"),
		},
		AuditConfig: &model.AuditConfig{},
		Concurrency: &model.Concurrency{
			WorkerPoolDelete: &model.WorkerPoolDelete{
				CountWorkers:   3,
				InputChainSize: 20,
				BufferSize:     10,
				BatchSize:      10,
			},
			WorkerPoolEvent: &model.WorkerPoolEvent{
				CountWorkers:   3,
				EventChainSize: 100,
			},
		},
		TLSConfig: &model.TLSConfig{
			Cert: "cert/cert.pem",
			Key:  "cert/key.pem",
		},
		ConfigFile: &model.ConfigFile{},
	}

	err = cfg.LoadConfigEnv()
	if err != nil {
		logHelper.LogError(ctx, "LoadConfigEnv", err)
	}
	cfg.LoadConfigFlag()

	configFile, err := LoadConfigFile(cfg.ConfigFile.Path)
	if err != nil {
		logHelper.LogError(ctx, "Error loading config file", err)
	} else {
		configFile.LoadConfig(&cfg)
	}

	return &cfg, err
}
