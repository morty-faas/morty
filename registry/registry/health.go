package registry

import (
	"encoding/json"
	"net/http"
)

type healthcheck struct {
	Status string `json:"status"`
}

const (
	up           = "UP"
	outOfService = "OUT_OF_SERVICE"
)

func (s *Server) HealthcheckHandler(w http.ResponseWriter, _ *http.Request) {

	if err := s.storage.Healthcheck(); err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		json.NewEncoder(w).Encode(&healthcheck{Status: outOfService})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&healthcheck{Status: up})
}
