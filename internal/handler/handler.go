package handler

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/ArtShib/urlshortener/internal/model"
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
	if err != nil && err != model.ErrURLConflict {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("Content-Length", strconv.Itoa(len(shortURL)))
	
	if err == model.ErrURLConflict {
		w.WriteHeader(http.StatusConflict)
	}else{
		w.WriteHeader(201)
	}

	w.Write([]byte(shortURL))
}

func (h *URLHandler) GetID(w http.ResponseWriter, r *http.Request) {

	shortCode := r.URL.Path[1:]
	if shortCode == "" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	originalURL, err := h.service.GetID(shortCode) 
	
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Location", originalURL)
	w.WriteHeader(307)

}

func (h *URLHandler) ShortenJSON(w http.ResponseWriter, r *http.Request) {
	
	var req *model.RequestShortener

	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()

	if err := decoder.Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	responseShortener, err := h.service.ShortenJSON(req.URL)
	if err != nil && err != model.ErrURLConflict {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	if err == model.ErrURLConflict {
		w.WriteHeader(http.StatusConflict)
	}else{
		w.WriteHeader(201)
	}
	encoder := json.NewEncoder(w)
	if err := encoder.Encode(responseShortener); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *URLHandler) Ping(w http.ResponseWriter, r *http.Request) {
	cxt, cansel := context.WithTimeout(r.Context(), 10 * time.Second)
	defer cansel()
	if err := h.service.Ping(cxt); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)	
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *URLHandler) ShortenJSONBatch(w http.ResponseWriter, r *http.Request) {
	
	var req model.RequestShortenerBatchArray

	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()

	if err := decoder.Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	responseShortener, err := h.service.ShortenJSONBatch(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	encoder := json.NewEncoder(w)
	if err := encoder.Encode(responseShortener); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
