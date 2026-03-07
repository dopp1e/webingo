package main

import (
	"net/http"
)

func generalErrorResponse(w http.ResponseWriter, r *http.Request, err error, status int, message string) {
	writeJSONError(w, status, message)
}

func (app *application) internalServerError(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Errorw("Internal Server Error", "error", err)
	generalErrorResponse(w, r, err, http.StatusInternalServerError, "Internal Server Error")
}

func (app *application) badRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Warnw("Bad Request", "error", err)
	generalErrorResponse(w, r, err, http.StatusBadRequest, err.Error())
}

func (app *application) notFoundResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Warnw("Not Found", "error", err)
	generalErrorResponse(w, r, err, http.StatusNotFound, err.Error())
}

func (app *application) conflictResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Warnw("Conflict", "error", err)
	generalErrorResponse(w, r, err, http.StatusConflict, err.Error())
}

func (app *application) unauthorizedResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Warnw("Unauthorized", "error", err)
	generalErrorResponse(w, r, err, http.StatusUnauthorized, err.Error())
}

func (app *application) unauthorizedBasicResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Warnw("Unauthorized (Basic)", "error", err)

	w.Header().Set("WWW-Authenticate", `Basic realm="Restricted", charset="UTF-8"`)
	generalErrorResponse(w, r, err, http.StatusUnauthorized, err.Error())
}

func (app *application) forbiddenResponse(w http.ResponseWriter, r *http.Request) {
	app.logger.Warnw("forbidden", "method", r.Method, "path", r.URL.Path, "error")

	writeJSONError(w, http.StatusForbidden, "forbidden")
}
