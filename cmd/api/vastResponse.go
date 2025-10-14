package main

import (
	"encoding/xml"
	"log"
	"net/http"
	"strconv"
)

func (app *application) getVastHandler(w http.ResponseWriter, r *http.Request) {
	println("getVastResponseHandler()")
	queryParams := r.URL.Query()
	dmaIdParam := queryParams.Get("dma")
	dmaIdInt, err := strconv.ParseInt(dmaIdParam, 10, 64)
	if err != nil {
		log.Print(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	ctx := r.Context()

	campaign, err := app.repository.Campaign.GetByDma(ctx, dmaIdInt)
	if err != nil {
		log.Print(err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	ip := app.getIpHost(r.RemoteAddr)
	currentDurationServed := app.rateLimiter.GetCurrentAdDurationServed(ip)
	log.Printf("getVastResponseHandler currentDurationServed: %v", currentDurationServed)

	vast, vastDuration, err := app.repository.VastResponse.GetVast(ctx, campaign, currentDurationServed)
	log.Printf("getVastResponseHandler new currentDurationServed: %v", vastDuration)

	if err != nil {
		log.Print(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	app.rateLimiter.UpdateCurrentAdDurationServed(ip, vastDuration)

	w.Header().Set("Content-Type", "application/xml")
	vastXml, err := xml.MarshalIndent(vast, "", "  ")
	if err != nil {
		log.Print(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Write(vastXml)
}
