package config

import (
	"flag"

	"github.com/ArtShib/urlshortener/internal/model"
	"github.com/caarlos0/env"
	"github.com/joho/godotenv"
)

type Config struct {
	HTTPServer *model.HTTPServerConfig 
	ShortService *model.ShortServiceConfig 
	RepoConfig *model.RepositoryConfig
} 

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
	return nil
}
func (c *Config) LoadConfigFlag() {
	if c.HTTPServer.ServerAddress == "" {
		flag.StringVar(&c.HTTPServer.ServerAddress, "a", ":8080", "HTTP server startup address")
	}
	if c.ShortService.BaseURL == "" {
		flag.StringVar(&c.ShortService.BaseURL, "b", "http://localhost:8080", "Address of the resulting shortened URL")
	}
	if c.RepoConfig.FileStoragePath == "" {
		flag.StringVar(&c.RepoConfig.FileStoragePath, "f", "/Users/shibakin-av/IdeaProjects/go/urlshortener/storage/storage.json", "File storage path")
	}
	flag.Parse()
}

func MustLoadConfig() *Config {
	cfg := Config{
		HTTPServer: &model.HTTPServerConfig{},
		ShortService: &model.ShortServiceConfig{},
		RepoConfig: &model.RepositoryConfig{},
	}
	cfg.LoadConfigEnv()
	cfg.LoadConfigFlag()
	return &cfg
}
