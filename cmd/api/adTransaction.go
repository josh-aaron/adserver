package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

func (app *application) logBeaconHandler(w http.ResponseWriter, r *http.Request) {
	println("logBeaconHandler()")

	transactionIdStr := app.ExtractQueryParam("t", r.URL.Query())
	transactionIdInt, err := strconv.ParseInt(transactionIdStr, 10, 64)
	if err != nil {
		log.Print(err)
		http.Error(w, "error: invalid transactionId format", http.StatusBadRequest)
		return
	}

	beaconName := app.ExtractQueryParam("cn", r.URL.Query())
	beaconUrl := r.URL.String()

	ctx := r.Context()

	err = app.repository.AdTransaction.LogBeacon(ctx, transactionIdInt, beaconUrl, beaconName)
	if err != nil {
		log.Print(err)
		http.Error(w, "error: error logging beacon", http.StatusInternalServerError)
		return
	}
}

func (app *application) getAdTransactionsHandler(w http.ResponseWriter, r *http.Request) {
	println("getAdTransactionsHandler()")
	w.Header().Set("Content-Type", "application/json")

	ctx := r.Context()
	adTransactions, err := app.repository.AdTransaction.GetAllAdTransactions(ctx)
	if err != nil {
		log.Print(err)
		http.Error(w, "error: error retrieving campaigns", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(adTransactions)
}

func (app *application) getBeaconsByTransactionIdHandler(w http.ResponseWriter, r *http.Request) {
	println("getBeaconsLogHandler()")
	w.Header().Set("Content-Type", "application/json")

	transactionIdParam := r.PathValue("transactionId")
	transactionIdInt, err := strconv.ParseInt(transactionIdParam, 10, 64)
	if err != nil {
		log.Print(err)
		http.Error(w, "error: invalid transactionId format", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	beacons, err := app.repository.AdTransaction.GetBeaconsByTransactionId(ctx, transactionIdInt)
	if err != nil {
		log.Print(err)
		http.Error(w, "error: error retrieving beacons", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(beacons)
}
