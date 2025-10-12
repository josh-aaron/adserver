package main

import (
	"log"
	"net/http"
	"time"

	"github.com/josh-aaron/adserver/internal/model"
)

type application struct {
	config     config
	repository model.Repository
}

type config struct {
	addr string
	db   dbConfig
}

type dbConfig struct {
	addr         string
	maxOpenConns int
	maxIdleConns int
	maxIdleTime  string
}

func (app *application) mount() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", app.healthCheckHandler)
	mux.HandleFunc("GET /campaigns", app.getCampaignsHandler)
	mux.HandleFunc("GET /campaigns/{id}", app.getCampaignById)
	mux.HandleFunc("POST /campaigns", app.createCampaignHandler)
	mux.HandleFunc("DELETE /campaigns/{id}", app.deleteCampaignHandler)
	mux.HandleFunc("PUT /campaings/{id}", app.updateCampaignHandler)

	return mux
}

func (app *application) run(handler http.Handler) error {

	server := &http.Server{
		Addr:         app.config.addr,
		Handler:      handler,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Minute,
	}

	log.Printf("server has started at %s", app.config.addr)

	return server.ListenAndServe()
}
