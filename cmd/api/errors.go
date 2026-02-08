package main

import (
	"log"
	"net/http"
)

func (app *application) internalServerErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("internal error: %s method: %s path: %s", err.Error(), r.Method, r.URL.Path)

	message := "server encountered an error while processing your request"
	writeJSONError(w, http.StatusInternalServerError, message)
}

func (app *application) badRequestErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("bad request error: %s method: %s path: %s", err.Error(), r.Method, r.URL.Path)

	writeJSONError(w, http.StatusBadRequest, err.Error())
}

func (app *application) notFoundErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("not found error: %s method: %s path: %s", err.Error(), r.Method, r.URL.Path)

	message := "record not found"
	writeJSONError(w, http.StatusNotFound, message)
}

func (app *application) editConflictResponse(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("edit conflict error: %s method: %s path: %s", err.Error(), r.Method, r.URL.Path)

	message := "edit conflict"
	writeJSONError(w, http.StatusConflict, message)
}
