package routes

import (
	"backend/internal/model"
	"encoding/json"
	"log"
	"net/http"
)

func Get(mux *http.ServeMux) {
	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		tags := []model.SystemTag{
			{Id: 1, Name: "tag1"},
			{Id: 2, Name: "tag2"},
		}

		res, err := json.Marshal(model.SystemTagResponse{Tags: tags})
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
