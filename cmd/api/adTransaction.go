package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

func (app *application) logBeaconsHandler(w http.ResponseWriter, r *http.Request) {
	println("logBeaconsHandler()")

	transactionIdStr := app.ExtractQueryParam("t", r.URL.Query())
	transactionIdInt, err := strconv.ParseInt(transactionIdStr, 10, 64)
	if err != nil {
		log.Print(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	beaconName := app.ExtractQueryParam("cn", r.URL.Query())
	beaconUrl := r.URL.String()

	ctx := r.Context()

	err = app.repository.AdTransaction.LogBeacons(ctx, transactionIdInt, beaconUrl, beaconName)
	if err != nil {
		log.Print(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (app *application) getAdTransactionsHandler(w http.ResponseWriter, r *http.Request) {
	println("getAdTransactions")
	w.Header().Set("Content-Type", "application/json")

	ctx := r.Context()
	adTransactions, err := app.repository.AdTransaction.GetAllAdTransactions(ctx)
	if err != nil {
		log.Print(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(adTransactions)
}
