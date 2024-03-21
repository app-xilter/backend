package main

import (
	"backend/internal/durable"
	"backend/internal/model"
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

	if err := durable.ConnectDB(os.Getenv("DB_DSN")); err != nil {
		log.Fatal("Error connecting to database")
	}

	// migrate database
	if err := durable.Connection().AutoMigrate(
		&model.Tags{},
		&model.Tweets{}); err != nil {
		log.Fatal(err)
	}

	// init tags
	var count int64
	durable.Connection().Model(&model.Tags{}).Count(&count)
	if count == 0 {
		tags := []model.Tags{
			{Name: "football"},
			{Name: "basketball"},
			{Name: "sport"},
			{Name: "politics"},
			{Name: "news"},
			{Name: "economy"},
			{Name: "crypto"},
			{Name: "technology"},
			{Name: "music"},
		}

		if err := durable.Connection().Create(&tags).Error; err != nil {
			log.Fatal(err)
		}
	}

}

func main() {
	mux := http.NewServeMux()
	server.SetupRoutes(mux)

	middlewareMux := server.SetupMiddleware(mux)
	server.StartServer(middlewareMux, os.Getenv("SERVER_PORT"))
}
