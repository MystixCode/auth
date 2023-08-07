package api

import (
	"auth/log"
	"auth/services/health"

	"net/http"
)

type HealthEndpoint struct {
	service *health.Service
	log     *log.Logger
}

func NewHealthEndpoint(log *log.Logger, service *health.Service) *HealthEndpoint {
	return &HealthEndpoint{
		service: service,
		log:     log,
	}
}

func (e *HealthEndpoint) GetHealth(w http.ResponseWriter, _ *http.Request) {
	var status int
	healthResponse := e.service.GetHealth()

	switch healthResponse.State {
	case health.StateStarting:
		status = http.StatusAccepted
	case health.StateRunning:
		status = http.StatusOK
	case health.StateStopping:
		status = http.StatusExpectationFailed
	default:
		status = http.StatusInternalServerError
	}

	respond(w, e.log, status, "", healthResponse)
}
