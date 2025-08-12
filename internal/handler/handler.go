package handler

import (
	"io"
	"net/http"
	"strconv"
)


type UrlService interface {
	Shorten(url string) (string, error)
	GetID(shortCode string) (string, error)
}

type UrlHandler struct{
	service UrlService
}

func NewUrlHandler(svc UrlService) *UrlHandler {
	return &UrlHandler{
		service: svc,
	}
}

func (h *UrlHandler) Shorten(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Only POST requests are allowed!", http.StatusMethodNotAllowed)
		return
	}
	
	body, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	shortUrl, err := h.service.Shorten(string(body))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(201)
	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("Content-Length", strconv.Itoa(len(shortUrl)))

	w.Write([]byte(shortUrl))

}

func (h *UrlHandler) GetID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET requests are allowed!", http.StatusMethodNotAllowed)
		return
	}

	shortCode := r.URL.Path[1:]
	if shortCode == "" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	longUrl, err := h.service.GetID(shortCode) 
	
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Location", longUrl)
	w.WriteHeader(307)

}
