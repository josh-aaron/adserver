package main

import (
	"log"
	"net"
	"net/http"
	"time"

	"github.com/josh-aaron/adserver/internal/model"
	"github.com/josh-aaron/adserver/internal/ratelimiter"
)

type application struct {
	config      config
	repository  model.Repository
	rateLimiter ratelimiter.Limiter
}

type config struct {
	addr        string
	db          dbConfig
	rateLimiter ratelimiter.Config
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
	mux.HandleFunc("GET /campaigns/{id}", app.getCampaignByIdHandler)
	mux.HandleFunc("POST /campaigns", app.createCampaignHandler)
	mux.HandleFunc("DELETE /campaigns/{id}", app.deleteCampaignHandler)
	mux.HandleFunc("PUT /campaigns/{id}", app.updateCampaignHandler)

	// /ads?dma={dmaId}&
	mux.Handle("GET /ads", app.rateLimiterMiddleware(http.HandlerFunc(app.getVastHandler)))

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

func (app *application) rateLimiterMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		ip := app.getIpHost(r.RemoteAddr)

		if allow, retryAfter := app.rateLimiter.Allow(ip); !allow {
			w.Header().Set("Retry-After", retryAfter.String())
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (app *application) getIpHost(remoteAddr string) string {
	ip, _, err := net.SplitHostPort(remoteAddr)
	if err != nil {
		log.Println(err)
	}
	return ip
}
