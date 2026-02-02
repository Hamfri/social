package main

import (
	"log"
	"os"
	"social/internal/db"
	"social/internal/env"
	"social/internal/repository"
)

func main() {
	cfg := config{
		port: env.GetString("APP_PORT", ":8080"),
		env:  env.GetString("APP_ENV", "development"),
		db: dbConfig{
			dsn:          env.GetString("DB_DSN", ""),
			maxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 30),
			maxIdleConns: env.GetInt("DB_MAX_IDLE_CONNS", 30),
			maxIdleTime:  env.GetString("DB_MAX_IDLE_TIME", "15m"),
		},
	}

	db, err := db.New(cfg.db.dsn, cfg.db.maxOpenConns, cfg.db.maxIdleConns, cfg.db.maxIdleTime)
	if err != nil {
		log.Panic(err)
		os.Exit(1)
	}

	defer db.Close()

	log.Printf("database connection pool established")

	repository := repository.New(db)

	app := &application{
		config:     cfg,
		repository: repository,
	}

	mux := app.mount()

	log.Fatal(app.run(mux))
}
