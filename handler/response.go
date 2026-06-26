package handler

import (
	"encoding/json"
	"net/http"
)

// JSONResponse defines the structure of the API response envelope.
type JSONResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

func respondWithJSON(w http.ResponseWriter, status int, payload JSONResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func respondWithError(w http.ResponseWriter, status int, errMsg string) {
	respondWithJSON(w, status, JSONResponse{Success: false, Error: errMsg})
}

func respondWithSuccess(w http.ResponseWriter, status int, data interface{}) {
	respondWithJSON(w, status, JSONResponse{Success: true, Data: data})
}
