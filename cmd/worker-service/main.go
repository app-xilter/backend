package main

import (
	"backend/internal/durable"
	"backend/internal/server"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
)

func init() {
	// setup logger
	durable.SetupLogger()

	// load .env file
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("Error loading .env file")
	}

	// connect to database
	if err := durable.ConnectDB(os.Getenv("DB_DSN")); err != nil {
		log.Fatal("Error connecting to database")
	}
}

func main() {
	mux := http.NewServeMux()
	server.SetupRoutes(mux)

	middlewareMux := server.SetupMiddleware(mux)
	server.StartServer(middlewareMux, "8080")
}
