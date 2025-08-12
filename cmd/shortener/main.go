package main

import (
	"net/http"

	"github.com/ArtShib/urlshortener/internal/config"
	"github.com/ArtShib/urlshortener/internal/repository"
	"github.com/ArtShib/urlshortener/internal/router"
	"github.com/ArtShib/urlshortener/internal/service"
)

func main() {
	cfg := config.LoadConfig()
	repo := repository.NewRepository()
	svc := service.NewUrlService(repo, cfg.HttpServer)
	router := router.NewRouter(svc)

	err := http.ListenAndServe(cfg.HttpServer.Port, router)	
	if err != nil{
		panic(err)
	}
}
