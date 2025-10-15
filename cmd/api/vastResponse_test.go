package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/josh-aaron/adserver/internal/model"
)

func TestGetVastResponse(t *testing.T) {

	app := newTestApplication(t, config{})
	mux := app.mount()

	t.Run("should return valid vast response", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/ads?dma=501", nil)
		if err != nil {
			t.Fatal(err)
		}
		requestRecorder := httptest.NewRecorder()
		mux.ServeHTTP(requestRecorder, req)

		var actualVast *model.VAST
		err = json.NewDecoder(requestRecorder.Body).Decode(&actualVast)
		if err != nil {
			t.Fatal(err)
		}

		expectedVast := &model.VAST{
			Version:      "3.0",
			XsiNamespace: "http://www.w3.org/2001/XMLSchema",
			Ads: []model.Ad{
				{
					ID: 2,
					InLine: &model.InLine{
						AdSystem: &model.AdSystem{
							Version: "4.0",
							Name:    "Rockbot",
						},
						AdTitle: model.CDATAString{"ForBiggerEscapes"},
						Pricing: &model.Pricing{
							Model:    "cpm",
							Currency: "USD",
							Value:    "25.00",
						},
						Errors: []model.CDATAString{
							{"http://example.com/error"},
						},
						Impressions: []model.Impression{
							{
								ID:  "Impression-ID-01",
								URI: "http://example.com/impression",
							},
						},
						Creatives: []model.Creative{
							{
								ID:       102,
								Sequence: 1,
								Linear: &model.Linear{
									Duration: "00:00:15",
									TrackingEvents: []model.Tracking{
										{
											Event: "start",
											URI:   "http://example.com/tracking/start",
										},
										{
											Event: "complete",
											URI:   "http://example.com/tracking/complete",
										},
									},
									VideoClicks: &model.VideoClicks{
										ClickThroughs: []model.VideoClick{
											{
												ID:  "ClickThrough-Impression-01",
												URI: "http://iabtechlab.com",
											},
										},
									},
									MediaFiles: []model.MediaFile{
										{
											ID:                  "5241",
											Delivery:            "progressive",
											Type:                "video/mp4",
											Codec:               "",
											Bitrate:             500,
											Width:               400,
											Height:              300,
											MinBitrate:          360,
											MaxBitrate:          1080,
											Scalable:            true,
											MaintainAspectRatio: true,
											URI:                 "http://commondatastorage.googleapis.com/gtv-videos-bucket/sample/ForBiggerEscapes.mp4",
										},
									},
								},
							},
						},
						Extensions: &[]model.Extension{
							{
								Type: "iab-Count",
								Data: []byte(`<total_available><![CDATA[ 2 ]]></total_available>`),
							},
						},
					},
				},
			},
		}

		if expectedVast != actualVast {
			t.Error("expected campaign does not match actual campaign")
		}
	})
}
