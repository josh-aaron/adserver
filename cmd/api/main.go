package main

import (
	"log"

	"github.com/josh-aaron/adserver/internal/db"
	"github.com/josh-aaron/adserver/internal/env"
	"github.com/josh-aaron/adserver/internal/model"
)

func main() {

	config := config{
		addr: env.GetString("ADDR", ":8080"),
		db: dbConfig{
			addr:         env.GetString("DB_ADDR", "postgres://admin:adminpassword@localhost/adserver?sslmode=disable"),
			maxOpenConns: 30,
			maxIdleConns: 30,
			maxIdleTime:  "15m",
		},
	}

	db, err := db.New(config.db.addr, config.db.maxOpenConns, config.db.maxIdleConns, config.db.maxIdleTime)
	if err != nil {
		log.Print(err)
		log.Panic()
	}

	defer db.Close()
	log.Println("database connection pool established")

	repository := model.NewRepository(db)

	app := &application{
		config:     config,
		repository: repository,
	}

	mux := app.mount()
	log.Fatal(app.run(mux))
}
