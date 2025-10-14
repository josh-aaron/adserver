package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/josh-aaron/adserver/internal/model"
)

// TODO: Implement helper methods to reduce repetitive code (e.g., header setting, error handling)
// TODO: Implement validation for data sent by client.

func (app *application) getCampaignsHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("getCampaignsHandler()")
	w.Header().Set("Content-Type", "application/json")

	ctx := r.Context()
	campaigns, err := app.repository.Campaign.GetAll(ctx)
	if err != nil {
		log.Print(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(campaigns)
}

func (app *application) getCampaignByIdHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("getCampaignById()")
	w.Header().Set("Content-Type", "application/json")

	campaignIdParam := r.PathValue("id")
	campaignIdInt, err := strconv.ParseInt(campaignIdParam, 10, 64)
	if err != nil {
		log.Print(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	campaign, err := app.repository.Campaign.GetById(ctx, campaignIdInt)
	if err != nil {
		log.Print(err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(campaign)
}

func (app *application) deleteCampaignHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("deleteCampaignHandler()")
	w.Header().Set("Content-Type", "application/json")

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

	w.WriteHeader(http.StatusOK)
}

func (app *application) createCampaignHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("createCampaignHandler()")
	w.Header().Set("Content-Type", "application/json")

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

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(campaign)
}

func (app *application) updateCampaignHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("updateCampaignHandler()")
	w.Header().Set("Content-Type", "application/json")

	campaignIdParam := r.PathValue("id")
	campaignIdInt, err := strconv.ParseInt(campaignIdParam, 10, 64)
	if err != nil {
		log.Print(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var updatedCampaign model.Campaign
	err = json.NewDecoder(r.Body).Decode(&updatedCampaign)
	if err != nil {
		log.Print(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	err = app.repository.Campaign.Update(ctx, campaignIdInt, &updatedCampaign)
	if err != nil {
		switch err {
		case model.ErrNotFound:
			http.Error(w, err.Error(), http.StatusNotFound)
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(updatedCampaign)
}
