package router

import (
	"net/http"

	"github.com/ArtShib/urlshortener/internal/handler"
	"github.com/go-chi/chi"
)

type UrlService interface {
	Shorten(url string) (string, error)
	GetID(shortCode string) (string, error)
}

func NewRouter(svc UrlService) http.Handler {
	//mux := http.NewServeMux()
	mux := chi.NewRouter()
	handler := handler.NewUrlHandler(svc)

	// mux.HandleFunc(`/`, handler.Shorten)
	// mux.HandleFunc(`/{shortCode}`, handler.GetID)
	mux.Post("/", handler.Shorten)
	mux.Get("/{shortCode}", handler.GetID)
	
	return mux
}
