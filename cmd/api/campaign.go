package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/josh-aaron/adserver/internal/model"
)

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
