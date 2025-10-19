const videoElement = document.getElementById("videoElement");
const contentVideoSrc = videoElement.currentSrc
const playButton = document.getElementById("playButton");
const dmaTextArea = document.getElementById("dmaTextArea")

let defaultImpression = ""

playButton.addEventListener("click", () => {
	start();
});

videoElement.addEventListener("ended", (event) => {
	console.log("ended")
	if (event.target.currentSrc != contentVideoSrc) {
		videoElement.setAttribute("src", contentVideoSrc)
		videoElement.play()
	}
})

// Per IAB MRC guidelines, default ad impression should fire as soon as ad renders
videoElement.addEventListener("loadeddata", (event) => {
	console.log("loaded")
	if (event.target.currentSrc != contentVideoSrc) {
		fireEventCallbackUrl(defaultImpression)
	}
})

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
}

function getDefaultImpression(vastXml) {
	console.log("getDefaultImpression")
	const impressionNodes = vastXml.getElementsByTagName("Impression");
	const impressionUrl = impressionNodes[0].textContent
	defaultImpression = impressionUrl
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
