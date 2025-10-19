package main

import (
	"log"
	"net"
	"net/http"
	"time"

	"github.com/josh-aaron/adserver/internal/model"
	"github.com/josh-aaron/adserver/internal/ratelimiter"
)

// Leverage depenency injection so that our structs rely on behaviors, not specific implementations
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

// TODO: update routes with API verison, e.g., /v1/campaigns, /v1/ads
func (app *application) mount() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", app.healthCheckHandler)

	mux.HandleFunc("GET /campaigns", app.getCampaignsHandler)
	mux.HandleFunc("GET /campaigns/{id}", app.getCampaignByIdHandler)
	mux.HandleFunc("POST /campaigns", app.createCampaignHandler)
	mux.HandleFunc("DELETE /campaigns/{id}", app.deleteCampaignHandler)
	mux.HandleFunc("PUT /campaigns/{id}", app.updateCampaignHandler)

	// Use query parameters for the ad request, i.e., /ads?dma={dmaId}&
	// Wrap the getVastHandler in the rateLimiter middleware
	mux.Handle("GET /ads", app.rateLimiterMiddleware(http.HandlerFunc(app.getVastHandler)))

	// Endpoint for client video players to send impression, error, and quartile beacons, using a transactionId.
	// The transactionId will be unique to each ad request, and be dynmically appended to the beacon URIs in the VAST
	mux.HandleFunc("GET /beacons", app.logBeaconsHandler)

	mux.HandleFunc("GET /adTransactions", app.getAdTransactionsHandler)

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

		// If we try to send another ad request after already hitting the limit, return a 429 HTTP error
		//  with the amount of time until we can be served ads again in the 'Retry-after' header
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
