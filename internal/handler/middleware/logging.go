package middleware

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/ArtShib/urlshortener/internal/model"
	"github.com/go-chi/chi/middleware"
)

func New(log *slog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		log.Info("logger middleware enabled")
		fn := func(w http.ResponseWriter, r *http.Request) {
			entry := log.With(
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
				slog.String("remote_address", r.RemoteAddr),
				slog.String("user_agent", r.UserAgent()),
				slog.String("content_type", r.Header.Get("Content-Type")),
				slog.String("content_encoding", r.Header.Get("Content-Encoding")),
				slog.String("request_id", middleware.GetReqID(r.Context())),
			)

			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			t1 := time.Now()

			userID, _ := r.Context().Value(model.UserIDKey).(string)
			token, _ := r.Cookie("User")
			defer func() {
				entry.Info("request completed",
					slog.String("userID", userID),
					slog.String("token", token.Value),
					slog.Int("status", ww.Status()),
					slog.Int("bytes", ww.BytesWritten()),
					slog.String("duration", time.Since(t1).String()),
				)
			}()
			next.ServeHTTP(ww, r)
		}
		return http.HandlerFunc(fn)
	}
}
