# adserver

Welcome to the README for the adserver project! Please also refer to the WORKLOG file, which maps each git commit to a high level unit of work/implementation step.

## Tech Stack, Dependencies
- Go 1.25.2
- Postgres v1.10.9
- github.com/joho/godotenv 
- github.com/rs/vast
- Migrate CLI Tool

## Project Setup & Verification Steps
1. Clone adserver repo
2. Download/install Postgres (if necessary)
    - If installing for the first time, ensure that you take note of the password and install directory
3. Create a .env file in the main directory (can copy .env.SAMPLE) and update with your Postgres username, password, and preferred port for the HTTP server.
4. Create the adserver DB
    First connect to the psql shell
    ```
    psql -U {username}
    ```
    And run this script:

    ```
    \i scripts/init_db.sql
    ```
5. Create the Campaign table by running the following commands in the psql shell. Or you can run the migrations in ``migrate/migrations`` (I used golang-migrate https://github.com/golang-migrate/migrate).
    
    ```
    \c adserver
    \i scripts/create_campaign_table.sql
    ```

6. Build and run the executable:
    - On Mac,  run ```make run``` in the terminal
    - On Windows,  run ```make runWindows``` in the terminal*

7. Test the API endpoints! Use the curl commands (or snag the campaign JSON objects) in curl-commands.txt to create some campaign data, or leverage an API testing tool like Postman or Thunderclient
    - See the "API Documentation" section below for details
    - You can find market names mapped to DMA codes here: https://www.spstechnical.com/DMACodes.htm
8. Test the other available Campaign API endpoints
9. Use the adrequest endpoint to retrieve a VAST response
10. Enter the VAST XML into an online VAST validator (e.g., https://tools.springserve.com/tagtest) and watch the ad play!
11. Run ```make test``` to run the unit tests

## API Documentaiton

### Campaign API

#### Get Campaign by Campaign ID

```
GET /campaigns/{id}
```
Url must contain the ID of a campaign, as an int.

Example:
```
GET /campaigns/{1}
```

Success:
- Returns a 200 HTTP status code and a JSON object of the Campaign if the campaign exists in the response body

Error:
- Returns a 404 HTTP error code if the campaign does not exist
- Retunrs a 400 HTTP error code if a non-integer is passed to the endpoint

#### Get All Campaigns

```
GET /campaigns/
```

Succes:
- Returns a 200 HTTP status code and a list of JSON objects of all Campaigns in the response body

Error
- Returns a 500 HTTP error code If there was any error retrieving all campaigns

#### Create Campaign

```
POST /campaigns/

Response Body:
{
    "name": "{name of the ad}", [string]
    "startDate": "{start date of the campaign}", [string]
    "endDate": "{end date of the campaign}", [string]
    "targetDmaId": {DMA code of the client}, [int]
    "adId": {adId}, [int]
    "adName": "{name of the ad creative}", [string]
    "adDuration": {duration of the ad creative, in seconds}, [int]
    "adCreativeId": {id of the ad creative}, [int]
    "AdCreativeUrl": "{url of the ad creative}" [string]
}
```
Request body must contain a JSON object of the campaign to be created.

Example:
```
POST /campaigns/{1}
{
    "name": "fender",
    "startDate": "2024-01-01",
    "endDate": "2025-01-01",
    "targetDmaId": 807,
    "adId": 4,
    "adName": "ForBiggerEscapes",
    "adDuration": 15,
    "adCreativeId": 104,
    "AdCreativeUrl": "http://commondatastorage.googleapis.com/gtv-videos-bucket/sample/ForBiggerEscapes.mp4"
}
```

Success:
- Returns a 201 HTTP status code and a JSON object of the Campaign that was created in the response body

Error:
- Returns a 400 HTTP error code if there is an issue with the JSON in the request body
- Returns a 500 HTTP error code for all other errors

#### Update Campaign

```
PUT /campaigns/{id}

Response Body:
{
    "name": "{name of the ad}", [string]
    "startDate": "{start date of the campaign}", [string]
    "endDate": "{end date of the campaign}", [string]
    "targetDmaId": {DMA code of the client}, [int]
    "adId": {adId}, [int]
    "adName": "{name of the ad creative}", [string]
    "adDuration": {duration of the ad creative, in seconds}, [int]
    "adCreativeId": {id of the ad creative}, [int]
    "AdCreativeUrl": "{url of the ad creative}" [string]
}
```
- Url must contain the ID of a campaign, as an int
- Request body must contain a JSON object of the campaign to be updated

Example:
```
PUT /campaigns/{1}

Response Body:
{
    "name": "gibson",
    "startDate": "2024-01-01",
    "endDate": "2025-01-01",
    "targetDmaId": 807,
    "adId": 4,
    "adName": "ForBiggerEscapes",
    "adDuration": 15,
    "adCreativeId": 104,
    "AdCreativeUrl": "http://commondatastorage.googleapis.com/gtv-videos-bucket/sample/ForBiggerEscapes.mp4"
}
```

Success:
- If the campaign was found and updated correctly, returns a 200 HTTP status code, and a JSON object of the Campaign that was updated in the response body

Error:
- Returns a 400 HTTP error code if there is an issue with the JSON in the request body
- Returns a 404 HTTP error code if the campaign does not exist 
- Returns a 500 HTTP error code for all other errors

#### Delete Campaign by Campaign ID

```
DELETE /campaigns/{id}
```
Url must contain the ID of a campaign, as an int.

Example:
```
DELETE /campaigns/{1}
```

Success:
- If the campaign exists, returns an 200 HTTP status code

Error:
- Returns a 500 HTTP error code for all other errors

### VAST Ad Response API

#### Get VAST Ad Response

```
GET /ads?dma={dmaCode}
```
Url must contain the DMA code of the client submitting the ad request, as an int.

Example:
```
GET /ads?dma=501
```

Success:
- Returns a 200 HTTP status code and a XML VAST ad response in the response body
- If a campaign is found with a target DMA code that matches the ad request, but the campaign is not active, then a 200 HTTP status code will be returned with an empty VAST ad response in the response body 

Error:
- Returns a 404 HTTP error code if no campaign is found with a target DMA code that matches the ad request's DMA code
- Returns a 500 HTTP error code for all other issues


## Assumptions and Limitations

- Only one publisher network (e.g., BobsAwesomeCatVideos.com is our only ad space publisher). Therefore, when any ad request is received, we only filter by targetDma and by if campaign is active
- DMA code/id will be hardcoded into ad request
- Ad and Creative data will live in the Campaign table. This is for ease of testing and updating the exposed Campaign API endpoints - i.e., testers can add/modify campaigns and see those changes reflected in the VAST response.
    - All other VAST data will be hardcoded.
    - The VAST response returned will only contain nodes outlined in the VAST 3.0 Inline Linear Example (https://github.com/InteractiveAdvertisingBureau/VAST_Samples/blob/master/VAST%203.0%20Samples/Inline_Linear_Tag-test.xml)
    - Excludes certain optional nodes (e.g., Wrapper, Non-Linear, Companion).
- If DMA in request does not match any active campaign target DMA, we will return an empty VAST
- A campaign will only have one DMA
- Only one ad/creative will be returned in the VAST response and all creatives will have a duration of 15 seconds
- IP address from HTTP request will be used to identify "unique users" for rate limitng

## Guiding Principles:
1. Use clean, layered architecture / Separation of controls (Transport, Service, Repoistory Layer)
2. Leverage dependency inversion principle, promote loose coupling
3. Strike a balance between familiar tech and learning opportunities

## Implementation Decisions, Considerations, Tradeoffs
WIP
- Tech stack
    - Goal was to use as few external packages/dependencies as possile
    - Postgres vs. MySQL
- Architecture/project set up
    - Combined service layer into repository layer
    - However, based on the complexity of the vast response service logic, it may have made more sense to separate out the business logic from the model
- Campaign API
- Ad Request API
- Rate Limiting


## Future Enhancements and Considerations

### Adtech Related
- Expand ad request query parameters (e.g., slot parameters)
- Implement external API for geolocation look up using RemoteAddr from adrequest HTTP request header
- Add table to register publisher networks
- Add additional tables mapped to each VAST node
- Wrap VAST file inside of an ad response
    - Allow for mutiple ads to be included in a single response, reducing number of requests received
    - Allow for temporal slots (i.e., preroll, midroll, postroll)
    - Allow for slotImpressions
        - slotImpressions can be used to calculate avails vs. filled
    - Or could use the Extension node in VAST for slot beacons/info
- Invalid traffic (IVT) filtering
    - Block traffic from known bots and automated tools
    - Suspicious non-human patterns (e.g., high number of clicks or requests)

### Engineering/Dev Related
- Implement rate limiting using Redis cache instead of in-memory
    - Distrubted cache allows for rate limiting for users across multiple server instances
- Authentication & Authorization
    - Especially required for managing campaigns
- Database caching for frequently used creatives
    - Use Redis to cache creatives frequently selected for insertion in ad responses
- Explore optimistic concurrency control/locks to avoid race conditions
    - Example: use timestamp/version column
- Log ad impressions, quartile beacons, etc. using event stereaming
    - Example: Kafka, RabbitMQ 
    - Can alo use ELK stack (Elasticsearch, Logstash, Kibana)
- Campaign Active vs. Inactive 
    - Run chronjob daily to designate campaign as active/inactive based on date that chronjob is run
    - May optimize ad response process by reducing the number of steps/calculations required

## New Experiences, Lessons Learned, and Takeaways
WIP

## Implementation Decision Making Hieracrhy
1. Leverage past experience from similar projects
2. Check official documentation
3. Use **reputable** StackOverflow posts, guides, and tutorials

Leverage 2 & 3 to find more optimal approaches from 1