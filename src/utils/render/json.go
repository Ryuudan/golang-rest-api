package render

import (
	"encoding/json"
	"log"
	"net/http"
)

func JSON(w http.ResponseWriter, status int, data interface{}) {
	response, err := json.Marshal(data)

	if err != nil {
		log.Printf("Error marshaling JSON: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)

	if _, err := w.Write(response); err != nil {
		log.Printf("Error writing response: %v", err)
	}
}
