package handler

import (
	"io"
	"net/http"
	"strconv"
)


type URLService interface {
	Shorten(url string) (string, error)
	GetID(shortCode string) (string, error)
}

type URLHandler struct{
	service URLService
}

func NewURLHandler(svc URLService) *URLHandler {
	return &URLHandler{
		service: svc,
	}
}

func (h *URLHandler) Shorten(w http.ResponseWriter, r *http.Request) {

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

	shortURL, err := h.service.Shorten(string(body))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(201)
	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("Content-Length", strconv.Itoa(len(shortURL)))

	w.Write([]byte(shortURL))

}

func (h *URLHandler) GetID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET requests are allowed!", http.StatusMethodNotAllowed)
		return
	}

	shortCode := r.URL.Path[1:]
	if shortCode == "" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	longURL, err := h.service.GetID(shortCode) 
	
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Location", longURL)
	w.WriteHeader(307)

}
