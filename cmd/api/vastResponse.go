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

	vast, err := app.repository.VastResponse.GetVast(ctx, campaign)
	if err != nil {
		log.Print(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/xml")
	vastXml, err := xml.MarshalIndent(vast, "", "  ")
	if err != nil {
		log.Print(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Write(vastXml)
}
