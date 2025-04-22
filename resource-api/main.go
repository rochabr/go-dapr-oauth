package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type dataResponse struct {
	Message string `json:"message"`
	User    string `json:"user"`
}

func main() {
	http.HandleFunc("/protected", func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token == "" {
			http.Error(w, "missing token", http.StatusUnauthorized)
			return
		}

		response := dataResponse{
			Message: "Access to protected resource granted",
			User:    "token claims would be decoded here",
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	})

	log.Println("[resource-api] Listening on :6002")
	log.Fatal(http.ListenAndServe(":6002", nil))
}
