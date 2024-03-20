package server

import (
	"backend/internal/routes"
	"fmt"
	"log"
	"net/http"
)

func SetupRoutes(mux *http.ServeMux) {
	routes.Post(mux)
	routes.GetTags(mux)
}

func StartServer(mux *http.ServeMux, port string) {
	fmt.Printf("Server is starting on port %s...\n", port)
	log.Fatal(http.ListenAndServe(":"+port, mux))
}
