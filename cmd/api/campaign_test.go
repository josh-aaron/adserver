package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/josh-aaron/adserver/internal/model"
)

// Use Go's built in testing library along with the mock structs, methods, and data
// TODO: write tests for the rest of the Campaign API, VAST API, and rate limiter

func TestGetCampaign(t *testing.T) {

	app := newTestApplication(t, config{})
	mux := app.mount()

	t.Run("should return valid Campaign", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/campaigns/1", nil)
		if err != nil {
			t.Fatal(err)
		}
		requestRecorder := httptest.NewRecorder()
		mux.ServeHTTP(requestRecorder, req)

		expectedCampaign := model.Campaign{
			Name:          "ford",
			StartDate:     "2025-10-12",
			EndDate:       "2026-01-01",
			TargetDmaId:   501,
			AdId:          2,
			AdName:        "ForBiggerEscapes",
			AdDuration:    1,
			AdCreativeId:  102,
			AdCreativeUrl: "http://commondatastorage.googleapis.com/gtv-videos-bucket/sample/ForBiggerEscapes.mp4",
		}

		var actualCampaign model.Campaign
		err = json.NewDecoder(requestRecorder.Body).Decode(&actualCampaign)
		if err != nil {
			t.Fatal(err)
		}

		if expectedCampaign != actualCampaign {
			t.Error("expected campaign does not match actual campaign")
		}
	})
}
