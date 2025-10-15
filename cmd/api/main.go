package main

import (
	"log"
	"time"

	"github.com/josh-aaron/adserver/internal/db"
	"github.com/josh-aaron/adserver/internal/env"
	"github.com/josh-aaron/adserver/internal/model"
	"github.com/josh-aaron/adserver/internal/ratelimiter"
)

func main() {

	// Load .env variables
	env.LoadEnv()

	// Initialize config struct, with required configurations for the server, db, and rateLimiter
	config := config{
		addr: env.GetString("PORT", ":8080"),
		db: dbConfig{
			addr:         env.GetString("DB_ADDR", "postgres://admin:adminpassword@localhost/adserver?sslmode=disable"),
			maxOpenConns: 30,
			maxIdleConns: 30,
			maxIdleTime:  "15m",
		},
		rateLimiter: ratelimiter.Config{
			AdDurationLimit: 300,
			TimeFrame:       time.Minute * 60,
		},
	}

	// Initialize the db,
	db, err := db.New(config.db.addr, config.db.maxOpenConns, config.db.maxIdleConns, config.db.maxIdleTime)
	if err != nil {
		log.Print(err)
		log.Panic()
	}

	defer db.Close()
	log.Println("database connection pool established")

	repository := model.NewRepository(db)

	rateLimiter := ratelimiter.NewFixedWindowLimiter(
		config.rateLimiter.AdDurationLimit,
		config.rateLimiter.TimeFrame,
	)

	app := &application{
		config:      config,
		repository:  repository,
		rateLimiter: rateLimiter,
	}

	// Initialize the multiplexer, aka API router
	mux := app.mount()

	// If there's an error running the app, let's shut it down
	log.Fatal(app.run(mux))
}
