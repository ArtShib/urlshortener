package getjsonbatch

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/ArtShib/urlshortener/internal/model"
	"github.com/go-chi/chi/v5/middleware"
)

// URLService интерфейс сервиса для получния списка url созданных пользователем.
type URLService interface {
	GetJSONBatch(ctx context.Context, userID string) (model.URLUserBatch, error)
}

// New конструктор HandlerFunc для получния списка url созданных пользователем.
func New(log *slog.Logger, svc URLService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "GetJSONBatch.Get"

		log := log.With(
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

		urlsBatch, err := svc.GetJSONBatch(r.Context(), userID)

		if err != nil {
			log.Error("GetJSONBatch", "error", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		if len(urlsBatch) == 0 {
			log.Error("StatusNoContent", "error", http.StatusText(http.StatusNoContent))
			http.Error(w, http.StatusText(http.StatusNoContent), http.StatusNoContent)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		encoder := json.NewEncoder(w)
		if err := encoder.Encode(urlsBatch); err != nil {
			log.Error("Encode response", "error", err)
			return
		}
	}
}
