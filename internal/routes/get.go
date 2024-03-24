package routes

import (
	"backend/internal/durable"
	"backend/internal/model"
	"encoding/json"
	"log"
	"net/http"
)

func Get(mux *http.ServeMux) {
	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		tagsList := model.SystemTagResponse{}
		if err := durable.Connection().Table("tags").Find(&tagsList.Tags).Order("id ASC").Error; err != nil {
			http.Error(w, "Database error: getting tags", http.StatusInternalServerError)
			log.Println(err)
			return
		}

		res, err := json.Marshal(tagsList)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if _, err = w.Write(res); err != nil {
			log.Printf("Error writing response: %v", err)
		}
	})
}
