package server

import (
	"fmt"
	"net/http"
	"time"
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
		fmt.Printf("%s	%s	%s	%d	%s\n", time.Now().Format("2006/01/02 15:04:05"), r.Method, r.URL, srw.status, r.RemoteAddr)
	})
}
