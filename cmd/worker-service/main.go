package main

import (
	"backend/internal/durable"
	"backend/internal/server"
	"flag"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
)

func init() {
	flag.Parse()

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

	loggingMux := server.LoggerMiddleware(mux)

	server.StartServer(loggingMux, "8080")
}
