package main

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

// extracts param from `url/id`
func (app *application) readIntParam(r *http.Request, key string) (int64, error) {
	id := chi.URLParam(r, key)
	val, err := strconv.ParseInt(id, 10, 64)

	if err != nil || val < 1 {
		return 0, errors.New("invalid int parameter")
	}
	return val, nil
}
