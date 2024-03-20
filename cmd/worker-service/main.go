package main

import (
	"backend/internal/durable"
	"backend/internal/model"
	"encoding/json"
	"flag"
	"fmt"
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
	http.HandleFunc("POST /", func(w http.ResponseWriter, r *http.Request) {
		enableCors(&w)

		defer func() {
			if err := recover(); err != nil {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				log.Printf("Panic: %v", err)
			}
		}()

		var t model.Tweets
		err := json.NewDecoder(r.Body).Decode(&t)
		if err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		err = durable.ValidateStruct(t)
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		w.Write([]byte(`{"message": "ok"}`))

	})

	fmt.Println("Server is starting on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}
