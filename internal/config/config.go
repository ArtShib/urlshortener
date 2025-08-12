package config

import (
	"flag"
	"strings"
	
	"github.com/ArtShib/urlshortener/internal/model"
)

type Config struct {
	HTTPServer *model.HTTPServerConfig
	ShortService *model.ShortServiceConfig
} 

func (c *Config) ParseConfig() {
	flag.StringVar(&c.HTTPServer.Port, "a", ":8080", "HTTP server startup address")
	flag.StringVar(&c.HTTPServer.Host, "b", "localhost", "Address of the resulting shortened URL")
	flag.Parse()
	c.HTTPServer.Port = strings.Replace(c.HTTPServer.Port, "::", ":", -1)
}

func LoadConfig() *Config {
	cfg := Config{
		HTTPServer: &model.HTTPServerConfig{},
		ShortService: &model.ShortServiceConfig{},
	}
	cfg.ParseConfig()
	return &cfg
}
