package model

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"time"
)

type VastResponseRepo struct {
	db *sql.DB
}

func (r *VastResponseRepo) GetVast(ctx context.Context, campaign *Campaign, totalDuration int, transactionId int64) (*VAST, int, error) {
	log.Println("vastResponse.GetVast()")

	isCampaignActive := checkIsCampaignActive(campaign.StartDate, campaign.EndDate)

	// Check if the campaign is active or inactive. If inactive, let's return an emtpy VAST response
	if !isCampaignActive {
		campaign = &Campaign{}
	}
	// TODO: implement logic that will append the transactionId to the impression url and trackingEvent urls
	vast, err := constructVast(campaign, transactionId)
	if err != nil {
		log.Println(err)
		return nil, 0, err
	}

	vastDuration := calculateTotalDuration(vast)

	return vast, vastDuration, nil
}

// TODO: Currently, every VAST that is returned contains one Creative that is 15 seconds long.
// We would need to update this logic to find the sum of the duration within every Linear node in the VAST.
func calculateTotalDuration(vast *VAST) int {
	if len(vast.Ads) == 0 {
		return 0
	}

	return 15
}

func constructCallbackUrl(callbackName string, transactionId int64) string {
	log.Printf("appendTransactionIdToUri, appending %v to %v%v", transactionId, callbackUrlHost, callbackName)
	transactionIdStr := strconv.FormatInt(transactionId, 10)
	return callbackUrlHost + "beacons?cn=" + callbackName + "&t=" + transactionIdStr
}

// TODO: Based on the rate limiting requirements, we would need to use the current ad duration served as part of the
// ad selection logic to ensure that we do not exceed the limit,
func constructVast(campaign *Campaign, transactionId int64) (*VAST, error) {
	log.Print("vastResponse.constructVast")
	var vast = &VAST{}
	if campaign.Id == 0 {
		log.Print("vastResponse.constructVast campaign is inactive, return empty VAST")
		vast = &VAST{
			Version:      vastVersion,
			XsiNamespace: vastXsiNamespace,
			Ads:          []Ad{},
		}
		return vast, nil
	}

	// TODO: Refactor with methods to instantiate each struct
	vast = &VAST{
		Version:      vastVersion,
		XsiNamespace: vastXsiNamespace,
		Ads: []Ad{
			{
				ID: campaign.AdId,
				InLine: &InLine{
					AdSystem: &AdSystem{
						Version: adSystemVersion,
						Name:    adSystemName,
					},
					AdTitle: CDATAString{campaign.AdName},
					Pricing: &Pricing{
						Model:    pricingModel,
						Currency: pricingCurrency,
						Value:    pricingValue,
					},
					Errors: []CDATAString{
						{constructCallbackUrl("error", transactionId)},
					},
					Impressions: []Impression{
						// In a real life scenario, the ImpressionId and URI would be dynamically generated or retrieved from a DB
						{
							ID:  impressionId,
							URI: constructCallbackUrl("defaultImpression", transactionId),
						},
					},
					Creatives: []Creative{
						{
							ID:       campaign.AdCreativeId,
							Sequence: sequence,
							Linear: &Linear{
								// TODO: UPDATE COLUMN TO BE EXPRESSED AS STRING IN THE FORMAT BELOW
								Duration: linearDuration,
								TrackingEvents: []Tracking{
									// TODO: UPDATE TRACKING WITH REST OF QUARTILES
									{
										Event: trackingEventStart,
										URI:   constructCallbackUrl("start", transactionId),
									},
									{
										Event: trackingEventComplete,
										URI:   constructCallbackUrl("complete", transactionId),
									},
								},
								VideoClicks: &VideoClicks{
									ClickThroughs: []VideoClick{
										{
											ID:  clickThroughId,
											URI: clickThroughURI,
										},
									},
								},
								MediaFiles: []MediaFile{
									{
										ID:                  mediaFileId,
										Delivery:            mediaFileDelivery,
										Type:                mediaFileType,
										Codec:               mediaFileCodec,
										Bitrate:             mediaFileBitrate,
										Width:               mediaFileWidth,
										Height:              mediaFileHeight,
										MinBitrate:          mediaFileMinBitrate,
										MaxBitrate:          mediaFileMaxBitrate,
										Scalable:            mediaFileScalable,
										MaintainAspectRatio: mediaFileMaintainAspectRation,
										URI:                 campaign.AdCreativeUrl,
									},
								},
							},
						},
					},
					Extensions: &[]Extension{
						{
							Type: extensionType,
							Data: []byte(extensionData),
						},
					},
				},
			},
		},
	}

	return vast, nil
}

// TODO: Think about if this would make more sense to live elsewhere
func checkIsCampaignActive(startDate string, endDate string) bool {
	log.Printf("vastResponse.checkIsCampaignActive()")

	// Convert startDate and endDate to timestamps
	layout := "2006-01-02"
	parsedStartDate, err := time.Parse(layout, startDate)
	if err != nil {
		fmt.Println("Error parsing date:", err)
		return false
	}
	parsedEndDate, err := time.Parse(layout, endDate)
	if err != nil {
		fmt.Println("Error parsing date:", err)
		return false
	}

	currentDate := time.Now()

	if currentDate.Compare(parsedStartDate) == -1 {
		log.Println("Campaign is not yet active")
		return false
	}

	if currentDate.Compare(parsedEndDate) == 0 || currentDate.Compare(parsedEndDate) == 1 {
		log.Println("Campaign has expired")
		return false
	}

	return true
}
