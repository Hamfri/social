package main

import (
	"encoding/base64"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrEmptyAuthorizationHeader     = errors.New("authorization header is empty")
	ErrMalformedAuthorizationHeader = errors.New("malformed authorization header")
)

// Terrible idea
// Don't use in any production system
func (app *application) BasicAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authorizationHeader := r.Header.Get("Authorization")

		if authorizationHeader == "" {
			app.unauthorizedBasicResponse(w, r, ErrEmptyAuthorizationHeader)
			return
		}

		parts := strings.Split(authorizationHeader, " ")
		if len(parts) != 2 || parts[0] != "Basic" {
			app.unauthorizedBasicResponse(w, r, ErrMalformedAuthorizationHeader)
			return
		}

		decoded, err := base64.StdEncoding.DecodeString(parts[1])
		if err != nil {
			app.unauthorizedBasicResponse(w, r, err)
			return
		}

		creds := strings.SplitN(string(decoded), ":", 2)
		if creds[0] != app.config.auth.basic.username || creds[1] != app.config.auth.basic.password {
			app.unauthorizedBasicResponse(w, r, errors.New("invalid password or username"))
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (app *application) TokenAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			app.unauthorizedResponse(w, r, ErrEmptyAuthorizationHeader)
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			app.unauthorizedResponse(w, r, ErrMalformedAuthorizationHeader)
			return
		}

		jwtToken, err := app.authenticator.ValidateToken(parts[1])
		if err != nil {
			app.unauthorizedResponse(w, r, err)
			return
		}

		claims := jwtToken.Claims.(jwt.MapClaims)

		sub, err := claims.GetSubject()
		if err != nil {
			app.unauthorizedResponse(w, r, err)
			return
		}

		userId, err := strconv.ParseInt(sub, 10, 64)
		if err != nil {
			app.unauthorizedResponse(w, r, err)
			return
		}

		ctx := r.Context()

		user, err := app.repository.Users.GetByID(ctx, userId)
		if err != nil {
			app.unauthorizedResponse(w, r, err)
			return
		}

		r = app.setAuthUserContext(r, user)

		next.ServeHTTP(w, r)
	})
}
