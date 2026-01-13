// Package deleteurls предоставляет функциональность для асинхронного удаления URL-адресов.
package deleteurls

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/ArtShib/urlshortener/internal/model"
	"github.com/go-chi/chi/v5/middleware"
)

// WorkerPoolDelete интерфейс воркера для обработки запросов на удаление.
type WorkerPoolDelete interface {
	AddRequest(req model.DeleteRequest)
}

// New конструктор HandlerFunc для обработки запросов на удаление URL.
func New(log *slog.Logger, svc WorkerPoolDelete) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "Delete.Delete"

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

		var uuids []string
		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&uuids); err != nil {
			log.Error("JsonDecode", "error", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		deleteRequest := model.DeleteRequest{
			UserID: userID,
			UUIDs:  uuids,
		}

		svc.AddRequest(deleteRequest)

		w.WriteHeader(http.StatusAccepted)
	}
}
