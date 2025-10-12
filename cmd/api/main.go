package main

import (
	"log"

	"github.com/josh-aaron/adserver/internal/env"
)

func main() {

	conf := config{
		addr: env.GetString("ADDR", ":8080"),
	}

	app := &application{
		config: conf,
	}

	mux := app.mount()
	log.Fatal(app.run(mux))
}
