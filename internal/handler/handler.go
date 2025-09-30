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

const (
	defaultRequestTimeout = 3 * time.Second
	longOperationTimeout  = 10 * time.Second
)

type URLHandler struct {
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

	shortURL, err := h.service.Shorten(r.Context(), string(body))
	if err != nil && err != model.ErrURLConflict {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("Content-Length", strconv.Itoa(len(shortURL)))

	if err == model.ErrURLConflict {
		w.WriteHeader(http.StatusConflict)
	} else {
		w.WriteHeader(201)
	}

	w.Write([]byte(shortURL))
}

func (h *URLHandler) GetID(w http.ResponseWriter, r *http.Request) {

	ctx, cancel := context.WithTimeout(r.Context(), longOperationTimeout)
	defer cancel()

	shortCode := r.URL.Path[1:]
	if shortCode == "" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	originalURL, err := h.service.GetID(ctx, shortCode)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Location", originalURL)
	w.WriteHeader(307)

}

func (h *URLHandler) ShortenJSON(w http.ResponseWriter, r *http.Request) {

	ctx, cancel := context.WithTimeout(r.Context(), longOperationTimeout)
	defer cancel()

	var req *model.RequestShortener

	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()

	if err := decoder.Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	responseShortener, err := h.service.ShortenJSON(ctx, req.URL)
	if err != nil && err != model.ErrURLConflict {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err == model.ErrURLConflict {
		w.WriteHeader(http.StatusConflict)
	} else {
		w.WriteHeader(201)
	}
	encoder := json.NewEncoder(w)
	if err := encoder.Encode(responseShortener); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *URLHandler) Ping(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), defaultRequestTimeout)
	defer cancel()
	if err := h.service.Ping(ctx); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *URLHandler) ShortenJSONBatch(w http.ResponseWriter, r *http.Request) {

	ctx, cancel := context.WithTimeout(r.Context(), longOperationTimeout)
	defer cancel()

	var req model.RequestShortenerBatchArray

	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()

	if err := decoder.Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	responseShortener, err := h.service.ShortenJSONBatch(ctx, req)
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

func (h *URLHandler) GetJSONBatch(w http.ResponseWriter, r *http.Request) {

	userID, _ := r.Context().Value("UserID").(string)
	if userID == "" {
		//http.Error(w, "Not found", http.StatusNotFound)
		//return
		//fmt.Printf("Unauthorized")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), longOperationTimeout)
	defer cancel()

	urlsBatch, err := h.service.GetJSONBatch(ctx, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	if len(urlsBatch) == 0 {
		http.Error(w, "Not content", http.StatusNoContent)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	encoder := json.NewEncoder(w)
	if err := encoder.Encode(urlsBatch); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
