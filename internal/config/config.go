package config

import (
	"flag"
	"os"

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
		flag.StringVar(&c.AuditConfig.AuditFile, "AUDIT_FILE", "", "Audit file path") ///home/artem/GolandProjects/urlshortener/storage/audit.json
	}
	if c.AuditConfig.AuditURL == "" {
		flag.StringVar(&c.AuditConfig.AuditURL, "AUDIT_URL", "", "URL to audit")
	}

	flag.Parse()
}

// MustLoadConfig конструктор Config
func MustLoadConfig() (*Config, error) {
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
	}
	err = cfg.LoadConfigEnv()
	cfg.LoadConfigFlag()
	return &cfg, err
}
