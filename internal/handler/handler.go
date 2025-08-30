package handler

import (
	"io"
	"net/http"
	"strconv"
)


type URLHandler struct{
	service URLService
}

func NewURLHandler(svc URLService) *URLHandler {
	return &URLHandler{
		service: svc,
	}
}

func (h *URLHandler) Shorten(w http.ResponseWriter, r *http.Request) {
	
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
	
	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("Content-Length", strconv.Itoa(len(shortURL)))
	w.WriteHeader(201)

	w.Write([]byte(shortURL))
}

func (h *URLHandler) GetID(w http.ResponseWriter, r *http.Request) {

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
