package main

import (
	"log"
	"net/http"
)

func generalErrorResponse(w http.ResponseWriter, r *http.Request, err error, status int, message string) {
	log.Printf("Error on %s %s, error: %s", r.Method, r.URL.Path, err)
	writeJSONError(w, status, message)
}

func (app *application) internalServerError(w http.ResponseWriter, r *http.Request, err error) {
	generalErrorResponse(w, r, err, http.StatusInternalServerError, "Internal Server Error")
}

func (app *application) badRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	generalErrorResponse(w, r, err, http.StatusBadRequest, err.Error())
}

func (app *application) notFoundResponse(w http.ResponseWriter, r *http.Request, err error) {
	generalErrorResponse(w, r, err, http.StatusNotFound, err.Error())
}
