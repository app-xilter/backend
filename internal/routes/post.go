package routes

import (
	"backend/internal/durable"
	"backend/internal/model"
	"encoding/json"
	"log"
	"net/http"
)

func Post(mux *http.ServeMux) {
	mux.HandleFunc("POST /", func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				log.Printf("Panic: %v", err)
			}
		}()

		var req model.Request
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		err = durable.ValidateStruct(req)
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}

		var responseModel model.Response
		for _, tweet := range req.Tweets {
			responseModel.Results = append(responseModel.Results, model.Result{
				Link:     tweet.Link,
				Category: 1,
			})
		}

		res, err := json.Marshal(responseModel)
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
