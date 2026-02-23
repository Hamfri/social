package main

import (
	"fmt"
	"net/http"
	"strconv"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func TestGetUser(t *testing.T) {
	redisCfg := config{
		redis: redisConfig{
			enabled: false,
		},
	}
	app := newTestApplication(t, redisCfg)

	mux := app.mount()

	claims := jwt.MapClaims{
		"aud": "mock-aud",
		"iss": "mock-iss",
		"sub": strconv.Itoa(int(1)),
		"exp": time.Now().Add(time.Hour).Unix(),
	}

	jwtToken, err := app.authenticator.GenerateToken(claims)
	if err != nil {
		t.Fatal(err)
	}
	t.Run("should not allow unauthenticated requests", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/v1/users/1", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := runRequest(req, mux)

		assertResponse(t, http.StatusUnauthorized, rr.Code)
	})

	t.Run("should allow authenticated requests", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/v1/users/1", nil)
		if err != nil {
			t.Fatal(err)
		}

		fmt.Println("token", jwtToken)
		req.Header.Set("Authorization", "Bearer "+*jwtToken)
		rr := runRequest(req, mux)

		assertResponse(t, http.StatusOK, rr.Code)
	})
}
