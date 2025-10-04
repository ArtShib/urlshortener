package middleware

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/ArtShib/urlshortener/internal/lib/auth"
	"github.com/ArtShib/urlshortener/internal/model"
)

func Auth(auth *auth.AuthService, logger *slog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var userID string

			cookie, err := r.Cookie("User")
			if err != nil {
				userID, err = auth.GenerateUserID()
				if err != nil {
					logger.Error("failed to generate user id",
						"Error", err.Error())
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				http.SetCookie(w, addCookie(userID, auth))

			} else if !auth.ValidateToken(cookie.Value) {
				userID, err = auth.GenerateUserID()
				if err != nil {
					logger.Error("failed to generate user id",
						"Error", err.Error())
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				http.SetCookie(w, addCookie(userID, auth))
			} else {
				userID = auth.GetUserID(cookie.Value)
				if userID == "" {
					logger.Error("failed to find user id",
						"UserID", "empty")
					http.Error(w, "Unauthorized", http.StatusUnauthorized)
					return
				}
				logger.Debug("user id is %s",
					"UserID", userID)
			}
			ctx := context.WithValue(r.Context(), model.UserIDKey, userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func addCookie(userID string, auth *auth.AuthService) *http.Cookie {
	token := auth.CreateToken(userID)
	return &http.Cookie{
		Name:  "User",
		Value: token,
	}
}
