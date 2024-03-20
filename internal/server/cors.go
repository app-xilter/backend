package server

import (
	"backend/internal/durable"
	"net/http"
)

func CorsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		durable.EnableCors(&w)
		next.ServeHTTP(w, r)
	})
}
