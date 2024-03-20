package server

import (
	"fmt"
	"net/http"
)

type statusResponseWriter struct {
	http.ResponseWriter
	status int
}

func (srw *statusResponseWriter) WriteHeader(code int) {
	srw.status = code
	srw.ResponseWriter.WriteHeader(code)
}

func LoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		srw := &statusResponseWriter{ResponseWriter: w}
		next.ServeHTTP(srw, r)
		fmt.Printf("%d %s %s (%s)\n", srw.status, r.Method, r.URL, r.RemoteAddr)
	})
}
