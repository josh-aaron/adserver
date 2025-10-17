package model

import "errors"

var (
	ErrNotFound = errors.New("resource not found")
)

// The following consts comprise the hardcoded data for the VAST response.
// The const names camelcase since they are not exported outside the model package.
const vastXsiNamespace = "http://www.w3.org/2001/XMLSchema"
const vastVersion = "3.0"
const adSystemVersion = "4.0"
const adSystemName = "Rockbot"
const pricingModel = "cpm"
const pricingCurrency = "USD"
const pricingValue = "25.00"
const errorURI = "localhost:8080/beacons?cn=error&t="
const impressionId = "Impression-ID-01"
const impressionURI = "http://example.com/impression"
const sequence = 1
const linearDuration = "00:00:15"
const trackingEventStart = "start"
const trackingEventStartURI = "localhost:8080/beacons?cn=start&t="
const trackingEventComplete = "complete"
const trackingEventCompleteURI = "localhost:8080/beacons?cn=complete&t="
const clickThroughId = "ClickThrough-Impression-01"
const clickThroughURI = "http://iabtechlab.com"
const mediaFileId = "5241"
const mediaFileDelivery = "progressive"
const mediaFileType = "video/mp4"
const mediaFileCodec = ""
const mediaFileBitrate = 500
const mediaFileWidth = 400
const mediaFileHeight = 300
const mediaFileMinBitrate = 360
const mediaFileMaxBitrate = 1080
const mediaFileScalable = true
const mediaFileMaintainAspectRation = true
const extensionType = "iab-Count"
const extensionData = `<total_available><![CDATA[ 2 ]]></total_available>`
