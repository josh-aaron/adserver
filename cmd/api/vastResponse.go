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

	// As soon as an ad request is allowed through by the rate limiter, create a transactionId
	transactionId := app.repository.AdTransaction.CreateTransactionId()

	dmaIdParam := app.ExtractQueryParam("dma", r.URL.Query())
	dmaIdInt, err := strconv.ParseInt(dmaIdParam, 10, 64)
	if err != nil {
		log.Print(err)
		http.Error(w, "error: dmaId format invalid", http.StatusInternalServerError)
		return
	}

	ctx := r.Context()

	// Before we can request/construct the VAST, we need to get the Campaign associated with the targetDMA
	// provided in query params of the ad request
	campaign, err := app.repository.Campaign.GetByDma(ctx, dmaIdInt)
	if err != nil {
		log.Print(err)
		http.Error(w, "error: error retrieving campaign for dma", http.StatusNotFound)
		return
	}

	// In the future, if there's a business requirement to limit ad duraiton served, then a real adserver would probably want to incorporate the duration already served
	// into it's ad selection process, to ensure it's not breaching the limit. So, let's pass the currentDurationServed to the VAST response service.
	ip := app.getIpHost(r.RemoteAddr)
	currentDurationServed := app.rateLimiter.GetCurrentAdDurationServed(ip)
	vast, vastDuration, err := app.repository.VastResponse.GetVast(ctx, campaign, currentDurationServed, transactionId)
	if err != nil {
		log.Print(err)
		http.Error(w, "error: error retrieving vast response", http.StatusInternalServerError)
	}

	// Update the current ad duration served in our in-memory cache with the total duration in the latest vast response
	app.rateLimiter.UpdateCurrentAdDurationServed(ip, vastDuration)

	// Marshal the VAST struct into xml to be written to the response
	vastXml, err := xml.MarshalIndent(vast, "", "  ")
	if err != nil {
		log.Print(err)
		http.Error(w, "error: error construct vast xml", http.StatusInternalServerError)
	}

	w.Write(vastXml)

	// I thought it might be more performant to pass individual variables to LogAdTransaction, as opposed to
	// instantiating a new AdTransaction struct. However, I should research if this goes against best practices
	app.repository.AdTransaction.LogAdTransaction(ctx, transactionId, r.URL.String(), vastXml, dmaIdInt, campaign.Id)
}
