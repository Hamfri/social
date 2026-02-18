package main

import (
	"os"
	"social/internal/auth"
	"social/internal/db"
	"social/internal/env"
	"social/internal/mailer"
	"social/internal/repository"
	"sync"
	"time"

	"go.uber.org/zap"
)

const version = "0.0.1"

//	@title			Social
//	@description	simple social network implementation
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

// @securityDefinitions.apiKey	ApiKeyAuth
// @in							header
// @name						Authorization
// @description
func main() {
	cfg := config{
		port:   env.GetString("APP_PORT", ":8080"),
		apiURL: env.GetString("APP_URL", "localhost:8080"),
		env:    env.GetString("APP_ENV", "development"),
		db: dbConfig{
			dsn:          env.GetString("DB_DSN", ""),
			maxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 30),
			maxIdleConns: env.GetInt("DB_MAX_IDLE_CONNS", 30),
			maxIdleTime:  env.GetString("DB_MAX_IDLE_TIME", "15m"),
		},
		smtp: smtp{
			host:     env.GetString("SMTP_HOST", ""),
			port:     env.GetInt("SMTP_PORT", 25),
			username: env.GetString("SMTP_USERNAME", ""),
			password: env.GetString("SMTP_PASSWORD", ""),
			sender:   env.GetString("SMTP_SENDER", ""),
		},
		auth: authConfig{
			basic: basicConfig{
				username: env.GetString("BASIC_AUTH_USERNAME", "admin"),
				password: env.GetString("BASIC_AUTH_PASSWORD", "admin"),
			},
			token: tokenConfig{
				secret: env.GetString("JWT_SECRET", ""),
				aud:    env.GetString("JWT_AUD", ""),
				iss:    env.GetString("JWT_ISS", ""),
				exp:    time.Duration(env.GetInt("JWT_EXP", 24)),
			},
		},
	}

	logger := zap.Must(zap.NewProduction()).Sugar()
	defer logger.Sync()

	db, err := db.New(cfg.db.dsn, cfg.db.maxOpenConns, cfg.db.maxIdleConns, cfg.db.maxIdleTime)
	if err != nil {
		logger.Fatal(err)
		os.Exit(1)
	}

	defer db.Close()

	logger.Info("database connection pool established")

	repository := repository.New(db)
	mailer, err := mailer.NewMailtrap(cfg.smtp.host, cfg.smtp.port, cfg.smtp.username, cfg.smtp.password, cfg.smtp.sender)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	JWTauthenticator := auth.NewJWTAuthenticator(cfg.auth.token.secret, cfg.auth.token.aud, cfg.auth.token.iss)

	app := &application{
		config:        cfg,
		repository:    repository,
		logger:        logger,
		mailer:        mailer,
		wg:            &sync.WaitGroup{},
		authenticator: JWTauthenticator,
	}

	mux := app.mount()

	logger.Fatal(app.run(mux))
}
