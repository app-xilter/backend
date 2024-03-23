package routes

import (
	"io"
	"net/http"
	"os"
)

func PrivacyPolicy(mux *http.ServeMux) {
	mux.HandleFunc("GET /privacy-policy", func(w http.ResponseWriter, r *http.Request) {
		file, err := os.Open("./assets/privacy-policy.html")
		if err != nil {
			http.Error(w, "Error reading file", http.StatusInternalServerError)
			return
		}
		defer func(file *os.File) {
			err := file.Close()
			if err != nil {
				return
			}
		}(file)

		if _, err = io.Copy(w, file); err != nil {
			http.Error(w, "Error writing response", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)
		return
	})
}
