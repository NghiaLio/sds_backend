package handler

import "net/http"

// PingHandler handles GET /ping requests.
type PingHandler struct{}

// NewPingHandler creates a new PingHandler.
func NewPingHandler() *PingHandler {
	return &PingHandler{}
}

// RegisterRoutes registers the ping route.
func (h *PingHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /ping", h.Ping)
}

// Ping returns a simple success pong message.
func (h *PingHandler) Ping(w http.ResponseWriter, r *http.Request) {
	respondWithSuccess(w, http.StatusOK, map[string]string{"message": "pong"})
}
