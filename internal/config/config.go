package config

import (
	"flag"

	"github.com/ArtShib/urlshortener/internal/model"
)

type Config struct {
	HttpServer *model.HttpServerConfig
	ShortService *model.ShortServiceConfig
} 

func (c *Config) ParseConfig() {
	flag.StringVar(&c.HttpServer.Port, "a", ":8080", "HTTP server startup address")
	flag.StringVar(&c.HttpServer.Host, "b", "localhost", "Address of the resulting shortened URL")
	flag.Parse()
}

func LoadConfig() *Config {
	cfg := Config{
		HttpServer: &model.HttpServerConfig{},
		ShortService: &model.ShortServiceConfig{},
	}
	cfg.ParseConfig()
	return &cfg
}
