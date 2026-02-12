package main

import (
	"net/http"
)

func (app *application) internalServerErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Errorw("internal error", "error", err.Error(), "method", r.Method, "path", r.URL.Path)

	message := "server encountered an error while processing your request"
	writeJSONError(w, http.StatusInternalServerError, message)
}

func (app *application) badRequestErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Warnf("bad request error", "error", err.Error(), "method", r.Method, "path", r.URL.Path)

	writeJSONError(w, http.StatusBadRequest, err.Error())
}

func (app *application) notFoundErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Warnf("not found error", "error", err.Error(), "method", r.Method, "path", r.URL.Path)

	message := "record not found"
	writeJSONError(w, http.StatusNotFound, message)
}

func (app *application) conflictResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Errorw("edit conflict error", "error", err.Error(), "method", r.Method, "path", r.URL.Path)

	writeJSONError(w, http.StatusConflict, err.Error())
}
