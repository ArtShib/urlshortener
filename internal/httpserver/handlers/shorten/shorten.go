package shorten

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/ArtShib/urlshortener/internal/model"
	"github.com/go-chi/chi/v5/middleware"
)

// URLService интерфейс сервиса для создания сокращенного url.
type URLService interface {
	Shorten(ctx context.Context, url string) (string, error)
}

// New конструктор HandlerFunc для создания сокращенного url.
func New(log *slog.Logger, svc URLService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "Shorten.Post"

		log := log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		log.Info("received request")

		r.Body = http.MaxBytesReader(w, r.Body, 1024*1024)

		body, err := io.ReadAll(r.Body)
		defer func() {
			if err := r.Body.Close(); err != nil {
				log.Error("close body", "error", err)
			}
		}()

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		shortURL, err := svc.Shorten(r.Context(), string(body))

		if err != nil && !errors.Is(err, model.ErrURLConflict) {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/plain")
		w.Header().Set("Content-Length", strconv.Itoa(len(shortURL)))

		if errors.Is(err, model.ErrURLConflict) {
			w.WriteHeader(http.StatusConflict)
		} else {
			w.WriteHeader(http.StatusCreated)
		}

		w.Header().Set("OriginalURL", string(body))
		_, err = w.Write([]byte(shortURL))
		if err != nil {
			log.Error("response write failed", "error", err)
		}
	}
}
