package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/ArtShib/urlshortener/internal/app"
	"github.com/ArtShib/urlshortener/internal/config"
)

func main() {
	cfg := config.MustLoadConfig()
	app, _ := app.NewApp(cfg)
	
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGQUIT)
	go app.Run()

	<-quit
	
	app.Stop()
}
