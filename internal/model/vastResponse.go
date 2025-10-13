package model

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"
)

// The structs below fulfill the VAST 3.0 Inline Linear Example (https://github.com/InteractiveAdvertisingBureau/VAST_Samples/blob/master/VAST%203.0%20Samples/Inline_Linear_Tag-test.xml
// Excludes certain optional nodes (e.g., Wrapper, Non-Linear, Companion).

type VAST struct {
	// The version of the VAST spec
	Version string `xml:"version,attr"`
	// One or more Ad elements. Advertisers and video content publishers may
	// associate an <Ad> element with a line item video ad defined in contract
	// documentation, usually an insertion order. These line item ads typically
	// specify the creative to display, price, delivery schedule, targeting,
	// and so on.
	Ads []Ad `xml:"Ad"`
}

type CDATAString struct {
	CDATA string `xml:",cdata"`
}

type Ad struct {
	// An ad server-defined identifier string for the ad
	ID     int     `xml:"id,attr,omitempty"`
	InLine *InLine `xml:",omitempty"`
}

type InLine struct {
	// The name of the ad server that returned the ad
	AdSystem *AdSystem
	// The common name of the ad
	AdTitle CDATAString
	// One or more URIs that directs the video player to a tracking resource file that the
	// video player should request when the first frame of the ad is displayed
	Impressions []Impression `xml:"Impression"`
	// The container for one or more <Creative> elements
	Creatives []Creative `xml:"Creatives>Creative"`
	// A URI representing an error-tracking pixel; this element can occur multiple
	// times.
	Errors []CDATAString `xml:"Error,omitempty"`
	// Provides a value that represents a price that can be used by real-time bidding
	// (RTB) systems. VAST is not designed to handle RTB since other methods exist,
	// but this element is offered for custom solutions if needed.
	Pricing *Pricing `xml:",omitempty"`
	// XML node for custom extensions, as defined by the ad server. When used, a
	// custom element should be nested under <Extensions> to help separate custom
	// XML elements from VAST elements. The following example includes a custom
	// xml element within the Extensions element.
	Extensions *[]Extension `xml:"Extensions>Extension,omitempty"`
}

type AdSystem struct {
	Version string `xml:"version,attr,omitempty"`
	Name    string `xml:",cdata"`
}

type Impression struct {
	ID  string `xml:"id,attr,omitempty"`
	URI string `xml:",cdata"`
}

type Pricing struct {
	// Identifies the pricing model as one of "cpm", "cpc", "cpe" or "cpv".
	Model string `xml:"model,attr"`
	// The 3 letter ISO-4217 currency symbol that identifies the currency of
	// the value provided
	Currency string `xml:"currency,attr"`
	// If the value provided is to be obfuscated/encoded, publishers and advertisers
	// must negotiate the appropriate mechanism to do so. When included as part of
	// a VAST Wrapper in a chain of Wrappers, only the value offered in the first
	// Wrapper need be considered.
	Value string `xml:",cdata"`
}

type Creative struct {
	// An ad server-defined identifier for the creative
	ID int `xml:"id,attr,omitempty"`
	// The preferred order in which multiple Creatives should be displayed
	Sequence int `xml:"sequence,attr,omitempty"`
	// Identifies the ad with which the creative is served
	AdID string `xml:"AdID,attr,omitempty"`
	// If present, defines a linear creative
	Linear *Linear `xml:",omitempty"`
}

type Linear struct {
	// To specify that a Linear creative can be skipped, the ad server must
	// include the skipoffset attribute in the <Linear> element. The value
	// for skipoffset is a time value in the format HH:MM:SS or HH:MM:SS.mmm
	// or a percentage in the format n%. The .mmm value in the time offset
	// represents milliseconds and is optional. This skipoffset value
	// indicates when the skip control should be provided after the creative
	// begins playing.
	SkipOffset *Offset `xml:"skipoffset,attr,omitempty"`
	// Duration in standard time format, hh:mm:ss
	Duration       string       `xml:"Duration"`
	TrackingEvents []Tracking   `xml:"TrackingEvents>Tracking,omitempty"`
	VideoClicks    *VideoClicks `xml:",omitempty"`
	MediaFiles     []MediaFile  `xml:"MediaFiles>MediaFile,omitempty"`
}

type Tracking struct {
	// The name of the event to track for the element. The creativeView should
	// always be requested when present.
	//
	// Possible values are creativeView, start, firstQuartile, midpoint, thirdQuartile,
	// complete, mute, unmute, pause, rewind, resume, fullscreen, exitFullscreen, expand,
	// collapse, acceptInvitation, close, skip, progress.
	Event string `xml:"event,attr"`
	// The time during the video at which this url should be pinged. Must be present for
	// progress event. Must match (\d{2}:[0-5]\d:[0-5]\d(\.\d\d\d)?|1?\d?\d(\.?\d)*%)
	Offset *Offset `xml:"offset,attr,omitempty"`
	URI    string  `xml:",cdata"`
}

type Offset struct {
	// If not nil, the Offset is duration based
	Duration *Duration
	// If Duration is nil, the Offset is percent based
	Percent float32
}

type Duration time.Duration

type VideoClicks struct {
	ClickThroughs  []VideoClick `xml:"ClickThrough,omitempty"`
	ClickTrackings []VideoClick `xml:"ClickTracking,omitempty"`
	CustomClicks   []VideoClick `xml:"CustomClick,omitempty"`
}

type VideoClick struct {
	ID  string `xml:"id,attr,omitempty"`
	URI string `xml:",cdata"`
}

type MediaFile struct {
	// Optional identifier
	ID string `xml:"id,attr,omitempty"`
	// Method of delivery of ad (either "streaming" or "progressive")
	Delivery string `xml:"delivery,attr"`
	// MIME type. Popular MIME types include, but are not limited to
	// “video/x-ms-wmv” for Windows Media, and “video/x-flv” for Flash
	// Video. Image ads or interactive ads can be included in the
	// MediaFiles section with appropriate Mime types
	Type string `xml:"type,attr"`
	// The codec used to produce the media file.
	Codec string `xml:"codec,attr,omitempty"`
	// Bitrate of encoded video in Kbps. If bitrate is supplied, MinBitrate
	// and MaxBitrate should not be supplied.
	Bitrate int `xml:"bitrate,attr,omitempty"`
	// Minimum bitrate of an adaptive stream in Kbps. If MinBitrate is supplied,
	// MaxBitrate must be supplied and Bitrate should not be supplied.
	MinBitrate int `xml:"minBitrate,attr,omitempty"`
	// Maximum bitrate of an adaptive stream in Kbps. If MaxBitrate is supplied,
	// MinBitrate must be supplied and Bitrate should not be supplied.
	MaxBitrate int `xml:"maxBitrate,attr,omitempty"`
	// Pixel dimensions of video.
	Width int `xml:"width,attr"`
	// Pixel dimensions of video.
	Height int `xml:"height,attr"`
	// Whether it is acceptable to scale the image.
	Scalable bool `xml:"scalable,attr,omitempty"`
	// Whether the ad must have its aspect ratio maintained when scales.
	MaintainAspectRatio bool `xml:"maintainAspectRatio,attr,omitempty"`
	// The APIFramework defines the method to use for communication if the MediaFile
	// is interactive. Suggested values for this element are “VPAID”, “FlashVars”
	// (for Flash/Flex), “initParams” (for Silverlight) and “GetVariables” (variables
	// placed in key/value pairs on the asset request).
	APIFramework string `xml:"apiFramework,attr,omitempty"`
	URI          string `xml:",cdata"`
}

type Extension struct {
	Type string `xml:"type,attr,omitempty"`
	Data []byte `xml:",innerxml"`
}

type VastResponseRepo struct {
	db *sql.DB
}

func (s *VastResponseRepo) GetVast(ctx context.Context, campaign *Campaign) (*VAST, error) {
	log.Println("vastResponse.GetByDma()")

	isCampaignActive, err := checkIsCampaignActive(campaign.StartDate, campaign.EndDate)
	if err != nil {
		//TODO: return custom error message
		return nil, err
	}

	var vast *VAST
	if !isCampaignActive {
		campaign = &Campaign{}
	}
	vast, err = constructVast(campaign)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return vast, nil
}

func constructVast(campaignPayload *Campaign) (*VAST, error) {
	log.Print("vastResponse.constructVast")
	var vast = &VAST{}
	if campaignPayload.Id == 0 {
		log.Print("vastResponse.constructVast campaign is inactive, return empty VAST")
		vast = &VAST{
			Version: "3.0",
			Ads:     []Ad{},
		}
		return vast, nil
	}

	//TODO: Refactor with methods to instantiate each struct
	vast = &VAST{
		Version: "3.0",
		Ads: []Ad{
			{
				ID: campaignPayload.AdId,
				InLine: &InLine{
					AdSystem: &AdSystem{
						Version: "4.0",
						Name:    "Rockbot",
					},
					AdTitle: CDATAString{campaignPayload.AdName},
					Pricing: &Pricing{
						Model:    "cpm",
						Currency: "USD",
						Value:    "25.00",
					},
					Errors: []CDATAString{
						//TODO: update beacons with unique transaction ID for tracking purposes
						{"http://example.com/error"},
					},
					Impressions: []Impression{
						{
							ID:  "Impression-ID-01",
							URI: "http://example.com/error",
						},
					},
					Creatives: []Creative{
						{
							ID:       campaignPayload.AdCreativeId,
							Sequence: 1,
							Linear: &Linear{
								//TODO: UPDATE COLUMN TO BE EXPRESSED AS STRING IN THE FORMAT BELOW
								Duration: "00:00:15",
								TrackingEvents: []Tracking{
									//TODO: UPDATE TRACKING WITH REST OF QUARTILES
									{
										Event: "start",
										URI:   "http://example.com/tracking/start",
									},
									{
										Event: "complete",
										URI:   "http://example.com/tracking/complete",
									},
								},
								VideoClicks: &VideoClicks{
									ClickThroughs: []VideoClick{
										{
											ID:  "ClickThrough-Impression-01",
											URI: "http://iabtechlab.com",
										},
									},
								},
								MediaFiles: []MediaFile{
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
										URI:                 campaignPayload.AdCreativeUrl,
									},
								},
							},
						},
					},
					Extensions: &[]Extension{
						{
							Type: "iab-Count",
							Data: []byte(`<total_available><![CDATA[ 2 ]]></total_available>`),
						},
					},
				},
			},
		},
	}

	return vast, nil
}

// TODO: Think about if this would make more sense to live elsewhere
func checkIsCampaignActive(startDate string, endDate string) (bool, error) {
	log.Println("vastResponse.checkIsCampaignActive()")
	// Convert startDate and endDate to timestamps
	layout := "2006-01-02"

	parsedStartDate, err := time.Parse(layout, startDate)
	if err != nil {
		fmt.Println("Error parsing date:", err)
		return false, err
	}
	parsedEndDate, err := time.Parse(layout, endDate)
	if err != nil {
		fmt.Println("Error parsing date:", err)
		return false, err
	}

	currentDate := time.Now()

	if currentDate.Compare(parsedStartDate) == -1 {
		log.Println("Campaign is not yet active")
		return false, nil
	}

	if currentDate.Compare(parsedEndDate) == 0 || currentDate.Compare(parsedEndDate) == 1 {
		log.Println("Campaign has expired")
		return false, nil
	}

	return true, nil
}
