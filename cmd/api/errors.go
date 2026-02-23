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
	app.logger.Warnw("bad request error", "error", err.Error(), "method", r.Method, "path", r.URL.Path)

	writeJSONError(w, http.StatusBadRequest, err.Error())
}

func (app *application) notFoundErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Warnw("not found error", "error", err.Error(), "method", r.Method, "path", r.URL.Path)

	message := "record not found"
	writeJSONError(w, http.StatusNotFound, message)
}

func (app *application) conflictResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Errorw("edit conflict error", "error", err.Error(), "method", r.Method, "path", r.URL.Path)

	writeJSONError(w, http.StatusConflict, err.Error())
}

// authentication
func (app *application) unauthorizedBasicResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Warnw("unauthorized error", "error", err.Error(), "method", r.Method, "path", r.URL.Path)

	w.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)

	writeJSONError(w, http.StatusUnauthorized, err.Error())
}

// authentication
func (app *application) unauthorizedResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Warnw("unauthorized error", "error", err.Error(), "method", r.Method, "path", r.URL.Path)
	writeJSONError(w, http.StatusUnauthorized, err.Error())
}

// authorization
func (app *application) notPermittedResponse(w http.ResponseWriter, r *http.Request) {
	message := "your account doesn't have the necessary permissions to access this resource"

	app.logger.Warnw(message, "method", r.Method, "path", r.URL.Path)
	writeJSONError(w, http.StatusForbidden, message)
}

func (app *application) rateLimiterExceededResponse(w http.ResponseWriter, r *http.Request, retryAfter string) {
	message := "rate limit exceeded"

	app.logger.Warnw(message, "method", r.Method, "path", r.URL.Path)

	w.Header().Set("Retry-After", retryAfter)

	writeJSONError(w, http.StatusTooManyRequests, message+": "+retryAfter)
}
