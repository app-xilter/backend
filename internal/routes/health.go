package routes

import (
	"log"
	"net/http"
)

func Health(mux *http.ServeMux) {
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)

		if _, err := w.Write([]byte("ok")); err != nil {
			log.Printf("Error writing response: %v", err)
		}
	})
}
