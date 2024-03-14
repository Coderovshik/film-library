package middleware

import (
	"log"
	"net/http"
)

func NewLogMiddleware() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Printf("INFO: incoming reqeust method=%s path=%s remoteAddr=%s\n", r.Method, r.URL.Path, r.RemoteAddr)

			next.ServeHTTP(w, r)

			log.Printf("INFO: processed reqeust successfully method=%s path=%s remoteAddr=%s\n", r.Method, r.URL.Path, r.RemoteAddr)
		})
	}
}
