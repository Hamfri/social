package main

import (
	"context"
	"errors"
	"expvar"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"social/docs"
	"social/internal/auth"
	"social/internal/cache"
	"social/internal/env"
	"social/internal/mailer"
	"social/internal/ratelimiter"
	"social/internal/repository"
	"sync"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	httpSwagger "github.com/swaggo/http-swagger/v2"
	"go.uber.org/zap"
)

type dbConfig struct {
	dsn          string
	maxOpenConns int
	maxIdleConns int
	maxIdleTime  string
}

type smtp struct {
	host     string
	port     int
	username string
	password string
	sender   string
}

// Terrible idea
// Don't use in any production system
type basicConfig struct {
	username string
	password string
}

type tokenConfig struct {
	secret string
	exp    time.Duration
	iss    string
	aud    string
}
type authConfig struct {
	basic basicConfig
	token tokenConfig
}

type redisConfig struct {
	addr    string
	pw      string
	db      int
	enabled bool
}

type config struct {
	port        string
	env         string
	db          dbConfig
	apiURL      string
	smtp        smtp
	auth        authConfig
	redis       redisConfig
	ratelimiter ratelimiter.Config
}
type application struct {
	config        config
	repository    repository.Repository
	redisCache    cache.Storage
	logger        *zap.SugaredLogger
	mailer        *mailer.SMTPMailer
	wg            *sync.WaitGroup
	authenticator *auth.JWTAuthenticator
	rateLimiter   *ratelimiter.FixedWindowRateLimiter
}

func (app *application) mount() http.Handler {
	r := chi.NewRouter()

	// Basic CORS
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{env.GetString("CORS_ALLOWED_ORIGIN", "localhost:4200")}, // avoid using arterisk
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Use(app.RateLimiterMiddleware)

	r.Use(middleware.Timeout(60 * time.Second))

	r.Route("/v1", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			// Terrible idea
			// in prod block such routes using nginx or caddy and only allow access via localhost
			// we can use an ssh tunnel locally to access the route
			r.Use(app.BasicAuthMiddleware)
			r.Get("/health", app.healthCheckHandler)
			docsURL := fmt.Sprintf("%s/swagger/doc.json", app.config.port)

			r.Get("/swagger/*", httpSwagger.Handler(httpSwagger.URL(docsURL)))
			r.Get("/debug/vars", expvar.Handler().ServeHTTP)
		})

		r.Route("/posts", func(r chi.Router) {
			r.Use(app.TokenAuthMiddleware)
			r.Post("/", app.createPostHandler)

			r.Route("/{id}", func(r chi.Router) {
				r.Use(app.postsContextMiddleware)
				r.Get("/", app.getPostHandler)
				r.Patch("/", app.CheckPostOwnership("moderator", app.updatePostHandler))
				r.Delete("/", app.CheckPostOwnership("admin", app.deletePostHandler))
			})
		})

		r.Route("/users", func(r chi.Router) {
			r.Use(app.TokenAuthMiddleware)
			r.Route("/{id}", func(r chi.Router) {
				r.Use(app.usersContextMiddleware)
				r.Get("/", app.getUserHandler)
				r.Put("/follow", app.followUserHandler)
				r.Delete("/unfollow", app.unFollowUserHandler)
			})

			r.Group(func(r chi.Router) {
				r.Get("/feed", app.getUserFeedHandler)
			})
		})

		r.Route("/accounts", func(r chi.Router) {
			r.Post("/register", app.registerAccountHandler)
			r.Put("/activate", app.activateAccountHandler)
			r.Post("/login", app.loginAccountHandler)
			r.Delete("/logout", app.loginAccountHandler)
		})
	})

	return r
}

// mux *chi.Mux == mux http.Handler
func (app *application) run(mux http.Handler) error {
	// swagger
	docs.SwaggerInfo.Version = version
	docs.SwaggerInfo.Host = app.config.apiURL
	docs.SwaggerInfo.BasePath = "/v1"
	docs.SwaggerInfo.Schemes = []string{"http", "https"}

	srv := &http.Server{
		Addr:         app.config.port,
		Handler:      mux,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Minute,
	}

	shutdownError := make(chan error)

	go func() {
		quit := make(chan os.Signal, 1)

		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

		s := <-quit

		app.logger.Infow("shutdown server", "signal", s.String())

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		shutdownError <- srv.Shutdown(ctx)
	}()

	app.logger.Infow("server started on", "port", app.config.port, "env", app.config.env)

	err := srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	if err = <-shutdownError; err != nil {
		return err
	}

	app.logger.Infow("server stopped", "port", app.config.port, "env", app.config.env)

	return nil
}
