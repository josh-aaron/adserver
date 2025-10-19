const videoElement = document.getElementById("videoElement");
const contentVideoSrc = videoElement.getAttribute("src");
const playButton = document.getElementById("playButton");

playButton.addEventListener("click", () => {
	start();
	console.log("after start()")
});

function start() {
	submitAdRequest()
	.then(vastXml => getMediaFile(vastXml))
	.then(adSrc => updateVideoElementSrc(adSrc))
	.finally(playAdVideo)
}

function playAdVideo(){
	console.log("playAdVideo")
	videoElement.play();
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
		const response = await fetch("/ads?dma=501");

		if (!response.ok) {
			throw new Error(`HTTP error! status: ${response.status}`);
		}
		const data = await response.text()
		xml = parseXmlString(data)
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
