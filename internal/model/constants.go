package model

import (
	"errors"

	"github.com/josh-aaron/adserver/internal/env"
)

var (
	ErrNotFound = errors.New("resource not found")
)

var callbackUrlHost string

func init() {
	env.LoadEnv()
	callbackUrlHost = env.GetCallbackUrlHost()

}

// The following consts comprise the hardcoded data for the VAST response.
// The const names camelcase since they are not exported outside the model package.
const vastXsiNamespace = "http://www.w3.org/2001/XMLSchema"
const vastVersion = "3.0"
const adSystemVersion = "4.0"
const adSystemName = "Rockbot"
const pricingModel = "cpm"
const pricingCurrency = "USD"
const pricingValue = "25.00"
const impressionId = "Impression-ID-01"
const sequence = 1
const linearDuration = "00:00:15"
const trackingEventStart = "start"
const trackingEventFirstQuartile = "firstQuartile"
const trackingEventMidpoint = "midpoint"
const trackingEventThirdQuartile = "thirdQuartile"
const trackingEventComplete = "complete"
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
