package main

import "net/http"

// This is a leftover from my original implementation of the webserver (ADS-001). Leaving here for nostalgia.

func (app *application) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("ok"))
}
