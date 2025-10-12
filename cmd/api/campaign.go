package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/josh-aaron/adserver/internal/model"
)

func (app *application) deleteCampaignHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("deleteCampaignHandler()")

	campaignIdParam := r.PathValue("id")
	campaignIdInt, err := strconv.ParseInt(campaignIdParam, 10, 64)
	if err != nil {
		log.Print(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	ctx := r.Context()

	err = app.repository.Campaign.Delete(ctx, campaignIdInt)
	if err != nil {
		log.Print(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (app *application) createCampaignHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("createCampaignHandler()")
	var newCampaign model.Campaign
	err := json.NewDecoder(r.Body).Decode(&newCampaign)
	if err != nil {
		log.Print(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	campaign := &model.Campaign{
		Id:            newCampaign.Id,
		Name:          newCampaign.Name,
		StartDate:     newCampaign.StartDate,
		EndDate:       newCampaign.EndDate,
		TargetDmaId:   newCampaign.TargetDmaId,
		AdId:          newCampaign.AdId,
		AdName:        newCampaign.AdName,
		AdDuration:    newCampaign.AdDuration,
		AdCreativeId:  newCampaign.AdCreativeId,
		AdCreativeUrl: newCampaign.AdCreativeUrl,
	}

	ctx := r.Context()
	err = app.repository.Campaign.Create(ctx, campaign)
	if err != nil {
		log.Print(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	js, err := json.Marshal(newCampaign)

	if err != nil {
		log.Print(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusCreated)
	w.Write(js)
}
