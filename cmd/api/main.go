package main

import (
	"os"
	"social/internal/db"
	"social/internal/env"
	"social/internal/repository"

	"go.uber.org/zap"
)

const version = "0.0.3"

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

	app := &application{
		config:     cfg,
		repository: repository,
		logger:     logger,
	}

	mux := app.mount()

	logger.Fatal(app.run(mux))
}
