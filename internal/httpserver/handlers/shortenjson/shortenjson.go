package shortenjson

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/ArtShib/urlshortener/internal/model"
	"github.com/go-chi/chi/v5/middleware"
)

// URLService интерфейс сервиса для создания сокращенного url. ответ json
type URLService interface {
	ShortenJSON(ctx context.Context, url string) (*model.ResponseShortener, error)
}

// New конструктор HandlerFunc для создания сокращенного url. ответ json
func New(log *slog.Logger, svc URLService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "ShortenJSON.Post"

		log := log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		log.Info("received request")

		var req model.RequestShortener

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

		responseShortener, err := svc.ShortenJSON(r.Context(), req.URL)
		if err != nil && !errors.Is(err, model.ErrURLConflict) {
			log.Error("service shortenJSON", "error", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if errors.Is(err, model.ErrURLConflict) {
			w.WriteHeader(http.StatusConflict)
		} else {
			w.WriteHeader(http.StatusCreated)
		}
		w.Header().Set("OriginalURL", req.URL)
		encoder := json.NewEncoder(w)
		if err := encoder.Encode(responseShortener); err != nil {
			log.Error("Encode response", "error", err)
			return
		}
	}
}
