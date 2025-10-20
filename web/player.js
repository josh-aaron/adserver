// TODO: refactor player.js into Class syntax

const videoElement = document.getElementById("videoElement");
const contentVideoSrc = videoElement.currentSrc;
const playButton = document.getElementById("playButton");
const dmaTextArea = document.getElementById("dmaTextArea");

const trackingEvents = {}
let defaultImpression = ""

let firstQuartileFired = false
let midpointFired = false
let thirdQuartileFired = false

playButton.addEventListener("click", () => {
	start();
});

videoElement.addEventListener("ended", (event) => {
	if (event.target.currentSrc != contentVideoSrc) {
		console.log("ad video ended")
		const trackingEventComplete = trackingEvents["complete"]
		if (trackingEventComplete) {
			fireEventCallbackUrl(trackingEventComplete)
		} else {
			console.log("no start tracking event in VAST")
		}
		videoElement.setAttribute("src", contentVideoSrc)
		videoElement.play()
		return
	}
	console.log("content video ended")
})

// Per IAB MRC guidelines, default ad impression should fire as soon as ad renders
videoElement.addEventListener("loadeddata", (event) => {
	if (event.target.currentSrc != contentVideoSrc) {
		console.log("ad video loaded")
		fireEventCallbackUrl(defaultImpression)
		return
	}
	console.log("content video loaded")
})

videoElement.addEventListener("playing", (event) => {
	if (event.target.currentSrc != contentVideoSrc) {
		console.log("ad video playing")
		const trackingEventStart = trackingEvents["start"]
		if (trackingEventStart) {
			fireEventCallbackUrl(trackingEventStart)
		} else {
			console.log("No start tracking event in VAST")
		}
		return
	}
	console.log("content video playing")

})

videoElement.addEventListener("timeupdate", (event) => {
	if (event.target.currentSrc != contentVideoSrc) {
		fireQuartileBeacons()
	}
})

function fireQuartileBeacons() {
	const currentTime = videoElement.currentTime
	const duration = videoElement.duration

	const trackingEventFirstQuartile = trackingEvents["firstQuartile"]
	const trackingEventMidpoint = trackingEvents["midpoint"]
	const trackingEventThirdQuartile = trackingEvents["thirdQuartile"]

	if (!firstQuartileFired && currentTime >= duration * 0.25) {
		console.log("fireQuartileBeacons firstQuartile")
		firstQuartileFired = true
		if (trackingEventFirstQuartile) {
			fireEventCallbackUrl(trackingEventFirstQuartile)
		} else {
			console.log("fireQuartileBeacons no firstQuartile tracking event in VAST")
		}
	} 
	else if (!midpointFired && currentTime >= duration * 0.5) {
		console.log("fireQuartileBeacons midpoint")
		midpointFired = true
		if (trackingEventMidpoint) {
			fireEventCallbackUrl(trackingEventMidpoint)
		} else {
			console.log("fireQuartileBeacons no midpoint tracking event in VAST")
		}
	} 
	else if (!thirdQuartileFired && currentTime >= duration * 0.75) {
		thirdQuartileFired = true
		if (trackingEventThirdQuartile) {
			fireEventCallbackUrl(trackingEventThirdQuartile)
		} else {
			console.log("fireQuartileBeacons no thirdQuartile tracking event in VAST")
		}
	}
}

async function fireEventCallbackUrl(urlString) {
	console.log("fireEventCallbackUrl firing " + urlString)
	const urlObject = new URL(urlString)
	const urlSearch = "/beacons" + urlObject.search
	console.log(urlSearch)
	const response = await fetch(urlSearch)
	if (!response.ok) {
		throw new Error(`HTTP error! status: ${response.status}`);
	} 
} 

function start() {
	submitAdRequest()
	.then(vastXml => getMediaFile(vastXml))
	.then(adSrc => updateVideoElementSrc(adSrc))
	.catch(error => {console.error("An error occurred:", error)})
	.finally(playAdVideo)
}

function playAdVideo(){
	console.log("playAdVideo")
	videoElement.play();
}

function getEventCallbacks(vastXml) {
	console.log("getEventCallbacks")
	getDefaultImpression(vastXml)
	getTrackingEvents(vastXml)
}

function getTrackingEvents(vastXml) {
	console.log("getTrackingEvents")
	const trackingEventNodes = vastXml.getElementsByTagName("Tracking");
	for (node of trackingEventNodes) {
		const eventName = node.getAttribute("event")
		const eventUrl = node.textContent
		trackingEvents[eventName] = eventUrl
	}
}

function getDefaultImpression(vastXml) {
	console.log("getDefaultImpression")
	const impressionNodes = vastXml.getElementsByTagName("Impression");
	if (impressionNodes) {
		const impressionUrl = impressionNodes[0].textContent
		defaultImpression = impressionUrl
	} else{
		console.log("impression node not found in VAST")
	}
	
}

function updateVideoElementSrc(adSrc) {
	console.log("updateVideoElementSrc")
	videoElement.setAttribute("src", adSrc);
	videoElement.width="640"
	videoElement.height="360"
}

function getMediaFile(vastXml) {
	console.log("getMediaFile")
	const mediaFileNodes = vastXml.getElementsByTagName("MediaFile");
	const adSrc = mediaFileNodes[0].textContent
	return adSrc
}

async function submitAdRequest() {
	console.log("submitAdRequest")
	try {
		const dma = dmaTextArea.value
		const adRequestPath = "/ads?dma=" + dma
		const response = await fetch(adRequestPath);
		if (!response.ok) {
			throw new Error(`HTTP error! status: ${response.status}`);
		}
		const data = await response.text()
		xml = parseXmlString(data)
		getEventCallbacks(xml)
		return xml
	} catch (error) {
		console.error("Error fetching data:", error);
		throw error;
	}
}

function parseXmlString(xmlString) {
	const parser = new DOMParser();
	const xmlDoc = parser.parseFromString(xmlString, "text/xml");
	if (xmlDoc.getElementsByTagName("parsererror").length > 0) {
		console.error("Error parsing XML:", xmlDoc.getElementsByTagName("parsererror")[0].textContent);
		return null;
	}
	return xmlDoc;
}
