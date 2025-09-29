package main

import (
	"net/http"
)

// healthCheckHandler godoc
//
//	@Summary		Healthcheck
//	@Description	Checks the health of the API service
//	@Tags			ops
//	@Produce		json
//	@Param			userID	path		string	true	"User ID"	Format(uuid)
//	@Success		200		{object}	string	"ok"
//	@Router			/health [get]
func (app *application) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{
		"status":  "ok",
		"env":     app.config.env,
		"version": version,
	}

	if err := app.jsonResponse(w, http.StatusOK, data); err != nil {
		app.internalServerError(w, r, err)
	}
}
