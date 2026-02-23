package main

import (
	"net/http"
	"net/http/httptest"
	"social/internal/auth"
	"social/internal/cache"
	"social/internal/ratelimiter"
	"social/internal/repository"
	"testing"

	"go.uber.org/zap"
)

func newTestApplication(t *testing.T, cfg config) *application {
	t.Helper()

	logger := zap.Must(zap.NewProduction()).Sugar()
	// logger := zap.NewNop().Sugar()
	mockRepostory := repository.MockNewRepository()

	mockCache := cache.MockNewRedisStorage()

	mockAuthenticator := auth.MockNewJWTAuthenticator("test", "mock-aud", "mock-iss")

	ratelimiter := ratelimiter.NewFixedWindowLimiter(
		cfg.ratelimiter.RequestsPerTimeFrame,
		cfg.ratelimiter.TimeFrame,
	)

	return &application{
		logger:        logger,
		repository:    mockRepostory,
		redisCache:    mockCache,
		authenticator: mockAuthenticator,
		config:        cfg,
		rateLimiter:   ratelimiter,
	}
}

func runRequest(req *http.Request, mux http.Handler) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)
	return rr
}

func assertResponse(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("expected: %d got: %d", expected, actual)
	}
}
