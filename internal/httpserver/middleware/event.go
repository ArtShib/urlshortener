package middleware

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/ArtShib/urlshortener/internal/model"
)

// ServiceEvent описывает интерфейс сохранения аудита
type ServiceEvent interface {
	AddEventRecord(event *model.Event)
}

// NewEvent конструктор middleware записи аудита
func NewEvent(log *slog.Logger, svc ServiceEvent) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			const op = "middleware.Event"
			logger := log.With(
				slog.String("op", op),
			)

			next.ServeHTTP(w, r)

			userID, ok := r.Context().Value(model.UserIDKey).(string)
			if !ok || userID == "" {
				logger.Error(op, "error", http.StatusText(http.StatusUnauthorized))
			}

			action := "shorten"
			if r.Method == http.MethodGet {
				action = "follow"
			}
			event := &model.Event{
				TimeStamp:   time.Now().Unix(),
				Action:      action,
				UserID:      userID,
				OriginalURL: w.Header().Get("OriginalURL"),
			}

			svc.AddEventRecord(event)
		})
	}
}
