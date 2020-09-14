package handlers

import (
	"encoding/json"
	"net/http"
	"time"
)

type Health struct {
	Start time.Time
}

type healthResponse struct {
	Uptime string `json:"uptime"`
}

func (h *Health) ServeHTTP(w http.ResponseWriter, _ *http.Request) {
	_ = json.NewEncoder(w).Encode(healthResponse{
		time.Since(h.Start).String(),
	})
}
