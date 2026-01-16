package shortenjsonbatch

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/ArtShib/urlshortener/internal/model"
	"github.com/go-chi/chi/v5/middleware"
)

// URLService интерфейс сервиса для создания пачки сокращенных url.
type URLService interface {
	ShortenJSONBatch(ctx context.Context, urls model.RequestShortenerBatchArray) (model.ResponseShortenerBatchArray, error)
}

// New конструктор HandlerFunc для создания пачки сокращенных url.
func New(log *slog.Logger, svc URLService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "ShortenJSONBatch.Post"

		log := log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		log.Info("received request")

		var req model.RequestShortenerBatchArray

		decoder := json.NewDecoder(r.Body)
		defer func() {
			if err := r.Body.Close(); err != nil {
				log.Error("close body", "error", err)
			}
		}()

		if err := decoder.Decode(&req); err != nil {
			log.Error("Body decode", "error", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		responseShortener, err := svc.ShortenJSONBatch(r.Context(), req)
		if err != nil {
			log.Error("service ShortenJSONBatch", "error", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		encoder := json.NewEncoder(w)
		if err := encoder.Encode(responseShortener); err != nil {
			log.Error("Encode response", "error", err)
			return
		}
	}
}
