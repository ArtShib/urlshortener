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
			var userID, token string

			c, err := r.Cookie("User")
			if err != nil {
				userID, err = auth.GenerateUserID()
				token = auth.CreateToken(userID)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
			} else if !auth.ValidateToken(c.Value) {
				userID, err = auth.GenerateUserID()
				token = auth.CreateToken(userID)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
			} else {
				userID = auth.GetUserID(c.Value)
				token = c.Value
				if userID == "" {
					http.Error(w, "Unauthorized", http.StatusUnauthorized)
					return
				}
			}

			http.SetCookie(w, &http.Cookie{
				Name:  "User",
				Value: token,
			})
			ctx := context.WithValue(r.Context(), model.UserIDKey, userID)
			//r = r.WithContext(ctx)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
