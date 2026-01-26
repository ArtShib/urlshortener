package stats

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/ArtShib/urlshortener/internal/model"
	"github.com/go-chi/chi/v5/middleware"
)

// URLService интерфейс сервиса для получния статистики
type URLService interface {
	Stats(ctx context.Context) (*model.Stats, error)
}

// New конструктор HandlerFunc для получния статистики
func New(logger *slog.Logger, svc URLService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "Stats.Get"

		log := logger.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		log.Info("received request")

		userID, ok := r.Context().Value(model.UserIDKey).(string)
		if !ok || userID == "" {
			log.Error("Unauthorized", "error", http.StatusText(http.StatusUnauthorized))
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		stats, err := svc.Stats(r.Context())

		if err != nil {
			log.Error("Stats", "error", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		encoder := json.NewEncoder(w)
		if err := encoder.Encode(stats); err != nil {
			log.Error("Encode response", "error", err)
			return
		}
	}
}
