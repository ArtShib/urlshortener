package middleware

import (
	"context"
	"net/http"

	"github.com/ArtShib/urlshortener/internal/lib/auth"
	"github.com/ArtShib/urlshortener/internal/model"
)

func Auth(auth *auth.AuthService) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var userID string

			cookie, err := r.Cookie("User")
			if err != nil {
				userID, err = auth.GenerateUserID()
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				http.SetCookie(w, addCookie(userID, auth))

			} else if !auth.ValidateToken(cookie.Value) {
				userID, err = auth.GenerateUserID()
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				http.SetCookie(w, addCookie(userID, auth))
			} else {
				userID = auth.GetUserID(cookie.Value)
				if userID == "" {
					http.Error(w, "Unauthorized", http.StatusUnauthorized)
					return
				}
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
