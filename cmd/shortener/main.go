package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/ArtShib/urlshortener/internal/config"
	"github.com/ArtShib/urlshortener/internal/handler"
	"github.com/ArtShib/urlshortener/internal/lib/logger"
	"github.com/ArtShib/urlshortener/internal/repository"
	"github.com/ArtShib/urlshortener/internal/service"
)

func main() {
	cfg := config.MustLoadConfig()
	repo, _ := repository.NewRepository(cfg.RepoConfig.FileStoragePath)
	defer func(){
		if err := repo.SavingRepository("/Users/shibakin-av/IdeaProjects/go/urlshortener/storage/storage3.json"); err != nil {
			fmt.Println(err)
		}
	}()
	
	svc := service.NewURLService(repo, cfg.ShortService)
	router := handler.NewRouter(svc, logger.NewLogger())

	go func() {
		log.Fatal(http.ListenAndServe(cfg.HTTPServer.ServerAddress, router))	
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT)
	<-sigChan
}
