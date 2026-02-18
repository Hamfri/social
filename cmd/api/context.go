package main

import (
	"context"
	"net/http"
	"social/internal/repository"
)

type contextKey string

const userContextKey = contextKey("user")

func (app *application) setAuthUserContext(r *http.Request, user *repository.User) *http.Request {
	ctx := context.WithValue(r.Context(), userContextKey, user)
	return r.WithContext(ctx)
}

func (app *application) getAuthUserContext(r *http.Request) *repository.User {
	user, ok := r.Context().Value(userContextKey).(*repository.User)
	if !ok {
		panic("user data missing in context")
	}

	return user
}
