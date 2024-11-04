package api

import (
	"encoding/json"
	"log"
	"main/lib/api/response"
	"net/http"
)

// TODO: add logger&response middleware for RequestID, log URL, timestamp, ip, etc.

func RespondWithOK(w http.ResponseWriter, requestID string, response response.Response, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	_ = json.NewEncoder(w).Encode(response)

	log.Printf("[INFO] [%s] Completed %d OK", requestID, status)
}

func RespondWithError(w http.ResponseWriter, requestID string, response response.Response, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(response)

	log.Printf("[ERROR]: [%s] %s", requestID, response.Error)
}
