package main

import (
	"log"
	"net/http"

	"github.com/ArtShib/urlshortener/internal/config"
	"github.com/ArtShib/urlshortener/internal/handler"
	"github.com/ArtShib/urlshortener/internal/repository"
	"github.com/ArtShib/urlshortener/internal/service"
)

func main() {
	cfg := config.MustLoadConfig()
	repo := repository.NewRepository()
	svc := service.NewURLService(repo, cfg.ShortService)
	router := handler.NewRouter(svc)

	log.Fatal(http.ListenAndServe(cfg.HTTPServer.Port, router))	
}
