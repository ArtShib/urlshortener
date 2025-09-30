package middleware

import (
	"context"
	"fmt"
	"net/http"

	"github.com/ArtShib/urlshortener/internal/lib/auth"
	"github.com/go-chi/chi/middleware"
)

func Auth(auth *auth.AuthService) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var userID string
			type contextKey string
			const userIDKey contextKey = "userID"
			c, err := r.Cookie("User")
			if err != nil {
				userID, err = auth.GenerateUserID()
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
			} else if !auth.ValidateToken(c.Value) {
				userID, err = auth.GenerateUserID()
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
			} else {
				userID = auth.GetUserID(c.Value)
				if userID == "" {
					http.Error(w, "Unauthorized", http.StatusUnauthorized)
					return
				}
			}
			token := auth.CreateToken(userID)
			fmt.Printf(userID)
			http.SetCookie(w, &http.Cookie{
				Name:  "User",
				Value: token,
			})
			ctx := context.WithValue(r.Context(), userIDKey, userID)
			r = r.WithContext(ctx)
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			next.ServeHTTP(ww, r)
		})
	}
}
