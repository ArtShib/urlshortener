package handler

import (
	"log/slog"
	"net/http"

	customMiddleware "github.com/ArtShib/urlshortener/internal/handler/middleware"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)


func NewRouter(svc URLService, log *slog.Logger) http.Handler {
	//mux := http.NewServeMux()
	mux := chi.NewRouter()

	mux.Use(middleware.RequestID)
	mux.Use(middleware.Logger)
	mux.Use(middleware.Recoverer)
	mux.Use(customMiddleware.New(log))
	mux.Use(customMiddleware.GzipMiddleware)

	handler := NewURLHandler(svc)
	mux.Post("/", handler.Shorten)
	mux.Post("/api/shorten", handler.ShortenJson)
	mux.Get("/{shortCode}",  handler.GetID)
	
	return mux
}
