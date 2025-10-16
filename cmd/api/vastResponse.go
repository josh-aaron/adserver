package main

import (
	"encoding/xml"
	"log"
	"net/http"
	"strconv"
)

func (app *application) getVastHandler(w http.ResponseWriter, r *http.Request) {
	println("getVastHandler()")
	w.Header().Set("Content-Type", "application/xml")

	queryParams := r.URL.Query()
	dmaIdParam := queryParams.Get("dma")
	dmaIdInt, err := strconv.ParseInt(dmaIdParam, 10, 64)
	if err != nil {
		log.Print(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	ctx := r.Context()

	// Before we can request/construct the VAST, we need to get the Campaign associated with the targetDMA
	// provided in query params of the ad request
	campaign, err := app.repository.Campaign.GetByDma(ctx, dmaIdInt)
	if err != nil {
		log.Print(err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// In the future, if there's a business requirement to limit ad duraiton served, then a real adserver would probably want to incorporate the duration already served
	// into it's ad selection process, to ensure it's not breaching the limit. So, let's pass the currentDurationServed to the VAST response service.
	ip := app.getIpHost(r.RemoteAddr)
	currentDurationServed := app.rateLimiter.GetCurrentAdDurationServed(ip)
	vast, vastDuration, err := app.repository.VastResponse.GetVast(ctx, campaign, currentDurationServed)
	if err != nil {
		log.Print(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	// Update the current ad duration served with the total duration in the latest vast response
	app.rateLimiter.UpdateCurrentAdDurationServed(ip, vastDuration)

	// Marshal the VAST struct into xml to be written to the response
	vastXml, err := xml.MarshalIndent(vast, "", "  ")
	if err != nil {
		log.Print(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Write(vastXml)
}
