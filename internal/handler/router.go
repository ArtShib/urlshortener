package handler

import (
	"log/slog"
	"net/http"

	customMiddleware "github.com/ArtShib/urlshortener/internal/handler/middleware"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)


func NewRouter(svc URLService, log *slog.Logger) http.Handler {

	mux := chi.NewRouter()

	mux.Use(middleware.RequestID)
	mux.Use(middleware.Logger)
	mux.Use(middleware.Recoverer)
	mux.Use(customMiddleware.New(log))
	mux.Use(customMiddleware.GzipMiddleware)

	handler := NewURLHandler(svc)
	mux.Post("/", handler.Shorten)
	mux.Post("/api/shorten", handler.ShortenJSON)
	mux.Get("/{shortCode}",  handler.GetID)
	mux.Get("/ping", handler.Ping)
	mux.Post("/api/shorten/batch", handler.ShortenJSONBatch)
	
	return mux
}

