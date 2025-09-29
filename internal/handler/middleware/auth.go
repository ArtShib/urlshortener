package middleware

import (
	"context"
	"net/http"

	"github.com/ArtShib/urlshortener/internal/lib/auth"
	"github.com/go-chi/chi/middleware"
)

func Auth(auth *auth.AuthService) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			c, err := r.Cookie("User")
			if err != nil {
				userId, err := auth.GenerateUserID()
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				token := auth.CreateToken(userId)
				http.SetCookie(w, &http.Cookie{
					Name:  "User",
					Value: token,
				})
				ctx := context.WithValue(r.Context(), "UserIDKey", userId)
				r = r.WithContext(ctx)

				ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
				next.ServeHTTP(ww, r)
				return
			}

			if !auth.ValidateToken(c.Value) {
				userId, err := auth.GenerateUserID()
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				token := auth.CreateToken(userId)
				http.SetCookie(w, &http.Cookie{
					Name:  "User",
					Value: token,
				})
				ctx := context.WithValue(r.Context(), "UserIDKey", userId)
				r = r.WithContext(ctx)
				//w.WriteHeader(http.StatusUnauthorized)
				//return
			}

			//userId := auth.GetUserID(c.Value)
			//ctx := context.WithValue(r.Context(), "UserIDKey", userId)
			//r = r.WithContext(ctx)

			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
			next.ServeHTTP(ww, r)
		}
		return http.HandlerFunc(fn)
	}
}
