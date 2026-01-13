package getid

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/ArtShib/urlshortener/internal/model"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// URLService интерфейс сервиса для получения оригинального url.
type URLService interface {
	GetID(ctx context.Context, shortCode string) (*model.URL, error)
}

// New конструктор HandlerFunc для получения оригинального url.
func New(log *slog.Logger, svc URLService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "GetID.Get"

		log := log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		log.Info("received request")

		shortCode := chi.URLParam(r, "shortCode")
		if shortCode == "" {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}

		url, err := svc.GetID(r.Context(), shortCode)
		if err != nil {
			log.Error("service GetID", "error", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		w.Header().Set("OriginalURL", url.OriginalURL)

		if url.DeletedFlag {

			w.WriteHeader(http.StatusGone)
			return
		}

		w.Header().Set("Location", url.OriginalURL)
		w.WriteHeader(http.StatusTemporaryRedirect)
	}
}
