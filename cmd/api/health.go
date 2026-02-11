package main

import (
	"net/http"
)

type HealthResponse struct {
	Status  string `json:"status"`
	Env     string `json:"env"`
	Version string `json:"version"`
}

// @Summary		App health check
// @Description	Check if application is up and running
// @Tags			Operations
// @Produce		json
// @Success		200	{object}	HealthResponse
// @Failure		500	{object}	error
// @Router			/health [get]
func (app *application) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{
		"status":  "ok",
		"env":     app.config.env,
		"version": version,
	}

	if err := app.jsonResponse(w, http.StatusOK, data); err != nil {
		app.internalServerErrorResponse(w, r, err)
	}
}
