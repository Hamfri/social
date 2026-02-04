package main

import (
	"log"
	"social/internal/db"
	"social/internal/env"
	"social/internal/repository"
)

func main() {
	appEnv := env.GetString("APP_ENV", "")
	if appEnv != "development" {
		log.Printf("this command can only be run on development instances only")
		return
	}

	dsn := env.GetString("DB_DSN", "")

	conn, err := db.New(dsn, 3, 3, "15m")
	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()

	repository := repository.New(conn)

	db.Seed(repository)
}
