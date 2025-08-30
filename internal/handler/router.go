package handler

import (
	"net/http"

	"github.com/go-chi/chi"
)


func NewRouter(svc URLService) http.Handler {
	//mux := http.NewServeMux()
	mux := chi.NewRouter()
	handler := NewURLHandler(svc)

	// mux.HandleFunc(`/`, handler.Shorten)
	// mux.HandleFunc(`/{shortCode}`, handler.GetID)
	mux.Post("/", handler.Shorten)
	mux.Get("/{shortCode}", handler.GetID)
	
	return mux
}
