package ping

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
)

// URLService интерфейс сервиса для проверки доступности хранилища.
type URLService interface {
	Ping(ctx context.Context) error
}

// New конструктор HandlerFunc для проверки доступности хранилища.
func New(log *slog.Logger, svc URLService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "Ping.Get"

		log := log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		log.Info("received request")

		if err := svc.Ping(r.Context()); err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}
