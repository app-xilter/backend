package server

import (
	"backend/internal/routes"
	"fmt"
	"log"
	"net/http"
)

func SetupRoutes(mux *http.ServeMux) {
	routes.Post(mux)
	routes.Get(mux)
	routes.Health(mux)
	routes.PrivacyPolicy(mux)
}

func SetupMiddleware(handler http.Handler) http.Handler {
	corsHandler := CorsMiddleware(handler)
	loggingHandler := LoggerMiddleware(corsHandler)
	return loggingHandler
}

func StartServer(handler http.Handler, port string) {
	fmt.Printf("Server is starting on port %s...\n", port)
	log.Fatal(http.ListenAndServe(":"+port, handler))
}
