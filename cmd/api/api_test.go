package main

import (
	"net/http"
	"net/http/httptest"
	"social/internal/ratelimiter"
	"testing"
	"time"
)

func TestRateLimiterMiddleware(t *testing.T) {
	cfg := config{
		ratelimiter: ratelimiter.Config{
			RequestsPerTimeFrame: 20,
			TimeFrame:            time.Second * 5,
			Enabled:              true,
		},
		port: "9000",
		auth: authConfig{
			basic: basicConfig{
				username: "admin",
				password: "admin",
			},
		},
	}

	app := newTestApplication(t, cfg)
	ts := httptest.NewServer(app.mount())
	defer ts.Close()

	client := &http.Client{}
	mockIP := "192.168.1.1"
	marginOfError := 2

	for i := 0; i < cfg.ratelimiter.RequestsPerTimeFrame+marginOfError; i++ {
		req, err := http.NewRequest("GET", ts.URL+"/v1/health", nil)
		if err != nil {
			t.Fatalf("could not initialise request: %v", err)
		}

		req.SetBasicAuth("admin", "admin")
		req.Header.Set("X-Forwarded-For", mockIP)
		resp, err := client.Do(req)

		if err != nil {
			t.Fatalf("could not send request: %v", err)
		}

		resp.Body.Close()

		if i < cfg.ratelimiter.RequestsPerTimeFrame {
			if resp.StatusCode != http.StatusOK {
				t.Errorf("expected status OK; got %v", resp.Status)
			}
		} else {
			if resp.StatusCode != http.StatusTooManyRequests {
				t.Errorf("expected status too many requests; got %v", resp.Status)
			}
		}
	}
}
