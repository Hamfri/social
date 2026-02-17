package main

import (
	"fmt"
	"net/http"
	"social/docs"
	"social/internal/mailer"
	"social/internal/repository"
	"sync"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
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
type config struct {
	port   string
	env    string
	db     dbConfig
	apiURL string
	smtp   smtp
}
type application struct {
	config     config
	repository repository.Repository
	logger     *zap.SugaredLogger
	mailer     *mailer.SMTPMailer
	wg         *sync.WaitGroup
}

func (app *application) mount() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Use(middleware.Timeout(60 * time.Second))

	r.Route("/v1", func(r chi.Router) {
		r.Get("/health", app.healthCheckHandler)

		docsURL := fmt.Sprintf("%s/swagger/doc.json", app.config.port)

		r.Get("/swagger/*", httpSwagger.Handler(httpSwagger.URL(docsURL)))

		r.Route("/posts", func(r chi.Router) {
			r.Post("/", app.createPostHandler)

			r.Route("/{id}", func(r chi.Router) {
				r.Use(app.postsContextMiddleware)
				r.Get("/", app.getPostHandler)
				r.Patch("/", app.updatePostHandler)
				r.Delete("/", app.deletePostHandler)
			})
		})

		r.Route("/users", func(r chi.Router) {
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

	app.logger.Infow("server started on", "port", app.config.port, "env", app.config.env)

	return srv.ListenAndServe()
}
