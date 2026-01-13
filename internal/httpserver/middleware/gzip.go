package middleware

import (
	"compress/gzip"
	"net/http"
	"strings"
)

// GzipMiddleware конструктор middleware gzip
func GzipMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		contentType := r.Header.Get("Content-Type")

		if r.Header.Get("Content-Encoding") == "gzip" {
			gzipReader, err := gzip.NewReader(r.Body)
			if err != nil {
				http.Error(w, "Invalid gzip body", http.StatusBadRequest)
				return
			}
			defer gzipReader.Close()
			r.Body = gzipReader
		}

		if strings.HasPrefix(contentType, "application/json") || strings.HasPrefix(contentType, "text/html") {

			if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
				next.ServeHTTP(w, r)
				return
			}

			w.Header().Set("Content-Encoding", "gzip")
			gzipWriter := gzip.NewWriter(w)
			defer gzipWriter.Close()

			gzw := gzipResponseWriter{ResponseWriter: w, Writer: gzipWriter}

			next.ServeHTTP(gzw, r)
		} else {
			next.ServeHTTP(w, r)
			return
		}
	})
}

type gzipResponseWriter struct {
	http.ResponseWriter
	Writer *gzip.Writer
}

// Write переопределенный метод записи gzip.Writer
func (w gzipResponseWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}
