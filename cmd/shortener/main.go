package main

import (
	"net/http"
	"strings"

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

	err := http.ListenAndServe(strings.Replace(cfg.HttpServer.Port, "::", ":", -1), router)	
	if err != nil{
		panic(err)
	}
}
