package config

import (
	"flag"
	
	"github.com/ArtShib/urlshortener/internal/model"
)

type Config struct {
	HTTPServer *model.HTTPServerConfig
	ShortService *model.ShortServiceConfig
} 

func (c *Config) LoadConfig() {
	flag.StringVar(&c.HTTPServer.Port, "a", ":8080", "HTTP server startup address")
	flag.StringVar(&c.ShortService.ShortURL, "b", "http://localhost:8080", "Address of the resulting shortened URL")
	flag.Parse()
}

func MustLoadConfig() *Config {
	cfg := Config{
		HTTPServer: &model.HTTPServerConfig{},
		ShortService: &model.ShortServiceConfig{},
	}
	cfg.LoadConfig()
	return &cfg
}
