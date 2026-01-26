package middleware

import (
	"net/http"

	"github.com/ArtShib/urlshortener/internal/lib/trustedsubnet"
)

// NewTrusted конструктор middleware проверки trustedSubnet
func NewTrusted(trustedSubnet *trustedsubnet.TrustedSubnet) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip := r.Header.Get("X-Real-IP")
			trusted := trustedSubnet.IsTrustedSubnet(r.Context(), ip)
			if !trusted {
				http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
