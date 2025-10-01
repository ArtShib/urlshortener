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

			_, err := r.Cookie("User")
			if err != nil {
				userID, err = auth.GenerateUserID()
				token = auth.CreateToken(userID)
				http.SetCookie(w, &http.Cookie{
					Name:  "User",
					Value: token,
				})
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
			}
			//else if !auth.ValidateToken(cookie.Value) {
			//	userID, err = auth.GenerateUserID()
			//	token = auth.CreateToken(userID)
			//	if err != nil {
			//		http.Error(w, err.Error(), http.StatusInternalServerError)
			//		return
			//	}
			//} else {
			//	userID = auth.GetUserID(cookie.Value)
			//	token = cookie.Value
			//	if userID == "" {
			//		http.Error(w, "Unauthorized", http.StatusUnauthorized)
			//		return
			//	}
			//}
			//log.Fatal("EEEEEEEEEEEEEEEEEEEE", cookie.Value)
			//userID = auth.GetUserID(cookie.Value)
			//if userID == "" {
			//	http.Error(w, "Unauthorized", http.StatusUnauthorized)
			//	return
			//}
			//token = cookie.Value
			//http.SetCookie(w, &http.Cookie{
			//	Name:  "User",
			//	Value: token,
			//})
			ctx := context.WithValue(r.Context(), model.UserIDKey, userID)
			//r = r.WithContext(ctx)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
