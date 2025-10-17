# adserver


Welcome to the README for the adserver project! Please also refer to the WORKLOG file, which maps each git commit to a high level unit of work/implementation step.


## Tech Stack & Dependencies
- Go 1.25.2
- Postgres v1.10.9
- github.com/joho/godotenv
- github.com/rs/vast
- golang-migrate (CLI Tool)


## Project Setup & Verification Steps

### Core Requirements

1. Clone adserver repo

2. Download/install Postgres (if necessary)
   - If installing for the first time, ensure that you take note of the password and install path (you'll need the install path to set your psql command environment variable)

3. Create a .env file in the main adserver directory (can copy .env.SAMPLE) and update with your Postgres username, password, and preferred port for the HTTP server.

4. Create the adserver DB

First, connect to the psql shell:
   ```
   psql -U {username}
   ```
   Then run this script:

   ```
   \i scripts/init_db.sql
   ```
5. Create the tables by running the following commands in the psql shell:
  
   ```
   \c adserver
   \i scripts/create_all_tables.sql
   ```

Note: I originally used migrations for the Campaign table - see ``migrate/migrations`` (using golang-migrate https://github.com/golang-migrate/migrate). However, to create the subsequent tables I just SQL in the psql shell.

6. Build and run the executable:
   - On Mac,  run ```make run``` in the terminal
   - On Windows,  run ```make runWindows``` in the terminal

7. Create some Campaign data and test the API endpoints!
   - Use the curl commands in curl-commands.txt to create some campaign data, or leverage an API testing tool like Postman or Thunderclient and snag the campaign JSON objects in curl-commands.txt
   - See the "API Documentation" section below for details
   - You can find market names mapped to DMA codes here: https://www.spstechnical.com/DMACodes.htm

8. Test the other available Campaign API endpoints (GET by id, GET all, PUT, and DELETE)

9. Use the ad response endpoint to retrieve a VAST response using GET requests
   - Ensure that: a. the ```dma``` parameter matches a targetDmaId from a Campaign you created, and b. the campaign is active in order to receive a populated VAST response
   - If a campaign with a matching dma is found, but the campaign is inactive, an empty VAST will be returned
   - You can check the console log to observe the total ad duration served. Once the limit of 300 seconds is reached, you can restart the server to request additional ads (or wait an hour!)

10. Enter the VAST XML into an online VAST validator (e.g., https://googleads.github.io/googleads-ima-html5/vsi/) and watch the ad play!
   - If you open up the browser console log and/or network tab, you can see the callback urls from the VAST firing

11. Run ```make test``` to run the unit test(s)

### Bonus Requirements

#### Logging of Ad Requests and VAST Ad Responses

- Whenever an ad request is allowed through by the rate limiter, a unique transaction ID is created (i.e., a UNIX time stamp in milliseconds). The transaction ID will be appended to the impression, error, and tracking event callback urls included in the VAST response, so that the consumer of the VAST response fires these callback beacons, they can be easily linked back to the ad request and response data
- All ad requests and responses will be recorded to the ```ad_transaction``` table, along with the associated transaction ID, DMA from the client ad request, and the associated campaign ID.

Verification Steps:
1. Ensure that at least one successful ad request has been sent to the VAST ad response API endpoint (i.e., endpoint returns either a populated VAST response or an empty VAST response with no ads due to inactive campaign)

2. Query the ad_transaction table using your method of choice. Options include using the psql shell, or DB Administration application of choice. Here are the required psql commands, if that's your style:

```
psql -U {username}
psql \c adserver
SELECT * FROM ad_transction
```

#### Logging of Impressions and Other Event Callbacks

- As previously mentioned, when an ad request is allowed through by the rate limiter, the transaction ID is appended to the event callback urls. The callback urls point towards a ```/beacons``` endpoint that has been exposed
- This endpoint extracts the transaction ID and callback name from the query parameters, and logs them to the DB along with the full beacon url

Verification Steps:
1. After submitting an ad request and successfully receiving a VAST response, grab the transaction ID from the callback urls in the VAST, or check the console log (if receiving an empty VAST due to an inactive campaign, the transaction ID will only be available in the console log)

2. Since we don't have a video player client to fire the callbacks (yet...), you can use a GET curl command or your API testing tool of choice to manually fire the callbacks
   - Use the transaction ID obtained in step 1 above for the ```t``` query parameter
   - Choose an event callback name (i.e., defaultImpression, error, start, fristQuartile, midPoint, thirdQuartile, complete). The endpoint currently does not have any validation and accepts any string

3. Query the ```ad_beacon``` table to see the logged event callback beacons:

```
psql -U {username}
psql \c adserver
SELECT * FROM ad_beacon
```

4. See the API Documentation section below for additional details

## API Documentation


### Campaign API


#### Get Campaign by Campaign ID


```
GET /campaigns/{id}
```
Url must contain the ID of a campaign, as an int.


Example:
```
GET /campaigns/30
```


Success:
- Returns a 200 HTTP status code and a JSON object of the Campaign if the campaign exists in the response body


Error:
- Returns a 404 HTTP error code if the campaign does not exist
- Returns a 400 HTTP error code if a non-integer is passed to the endpoint


#### Get All Campaigns


```
GET /campaigns/
```


Success:
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
POST /campaigns/
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
PUT /campaigns/6


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
DELETE /campaigns/7
```


Success:
- Returns a 200 HTTP status code if the campaign exists and was deleted successfully


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
- Returns a 429 HTTP error code if the client IP address exceeds the limit of 300 seconds of ad duration served. The limit will reset one hour after the first request is received from the client IP address
- Returns a 500 HTTP error code for all other issues


### Ad Impression Log API

#### Log Event callback


```
GET /beacons?cn={callbackName}&t={transactionId}
```
Url must contain query parameters for:
- callbackName, as a string
- transactionId, as an integer


Example:
```
GET /beacons?cn=defaultImpression&t=1760668810092
```

Success:
- Returns a 200 HTTP status code


Error:
- Returns a 400 HTTP error code if the transactionId is not an integer
- Returns a 500 HTTP error code for all other issues

## Assumptions and Limitations


- Only one publisher network (e.g., BobsAwesomeCatVideos.com is our only ad space publisher). Therefore, when any ad request is received, we only filter by targetDma and by if the campaign is active
- DMA code/id will be hardcoded into the ad request
- Ad and Creative data will live in the Campaign table. This is for ease of testing and updating the exposed Campaign API endpoints - i.e., testers can add/modify campaigns and see those changes reflected in the VAST response
   - All other VAST data will be hardcoded
   - The VAST response returned will only contain nodes outlined in the VAST 3.0 Inline Linear Example (https://github.com/InteractiveAdvertisingBureau/VAST_Samples/blob/master/VAST%203.0%20Samples/Inline_Linear_Tag-test.xml)
   - Excludes certain optional nodes (e.g., Wrapper, Non-Linear, Companion)
- If DMA in request does not match any active campaign target DMA, we will return an empty VAST
- A campaign will only have one DMA
- Only one ad/creative will be returned in the VAST response and all creatives will have a duration of 15 seconds
- IP address from HTTP request will be used to identify "unique users" for rate limiting


## Guiding Principles:
1. Use clean, layered architecture / separation of controls
2. Leverage dependency inversion principle, promote loose coupling
3. Strike a balance between familiar tech and learning opportunities


## Implementation Decisions, Considerations, Tradeoffs


- Tech stack
   - Goal was to use as few external packages/dependencies as possible
   - I chose Postgres instead of MySQL simply because I haven't worked with it before and wanted a new learning experience
- Architecture/project set up
   - Repository Pattern: Transport, Service, Repository Layers
       - The transport layer handles the API routing and project initialization logic. It is completely separate and independent from the business and database logic
       - Abstracting the database allows us to defer the choice of specific database to a later point, and allows for more efficient testing with mocking
       - Inspiration: https://threedots.tech/post/repository-pattern-in-go/
   - I ended up "combining" the service layer into the repository layer. However, based on the complexity of the vast response service logic, it may have made more sense to separate out the business logic from the model
   - While I had to spend extra time during the planning phase of the project to better understand the repository pattern and get more comfortable with dependency injection, I believe the time investment was worth it to make the application more extensible, flexible, and testable
- API routing
   - To align with my goal of minimizing external dependencies, I chose to stick with the net/http package from the Go standard library, and implement any middleware manually (i.e., rateLimiter)
   - Another good option would have been to use Chi (https://go-chi.io/). Despite being an external package, it is extremely popular. If I had ended up implementing more middleware (e.g., authentication, structured logging) it may have made sense to implement to improve code clarity
- Database tables
   - Here is where I must explain my most "controversial" decision: combining the ad and creative data with the campaign data into a single table
   - In a real world scenario, these would be segregated into different tables. However, given: a. the limited time frame of the project, and b. the need for  the requirements to be verified/tested in 10 - 20 minutes, I decided that it would be best to allow the tester to modify an ad/creative data using the Campaign API endpoints
   - A better approach, given more time for implementation and testing, would have been to keep the ad and creative data in its own table and expose another set of API endpoints for the ad/creatives. Then, when it comes to query, use a foreign key to join the campaign and ad/creative tables
- Campaign API
   - Went with standard CRUD operation APIs, while ensuring that no business logic is contained in the transport layer (i.e., the api folder within the main package)
- Ad Response API
   - I decided to use url query parameters for the API endpoint as it would be a good foundation for a future, more developed ad request with numerous query parameters (e.g., parameters for network profile, video asset dimensions, capability flags, etc.)
   - Due to the lack of specificity in the project requirements and the limited timeframe, I decided to keep the targeting and ad selection logic extremely simple, and made several assumptions as long as they did not conflict with the core project requirements (see the Assumptions & Limitations section for details)
   - One such assumption I implemented was that each campaign would only have one target DMA, and that the target DMAs would not repeat. Therefore, the logic for looking up campaigns/ads by target DMA is a simple "getByDma" method
       - However, if I was working under the assumption that a single DMA code could be associated with multiple campaigns, I would have instead implemented a "getAllByDma" method, that returns all campaigns with that DMA code
       - Then, combined with an updated assumption that a VAST response can contain more than one ad, I would implement a method that returns a slice of Ad structs, containing one Ad struct for each campaign/ad returned from the DB. Then that slice would be passed to the method(s) that construct the VAST response
- Rate Limiting
   - I decided to implement a fixed-window algorithm for the rate limiter
       - Another option would have been to use a sliding window algorithm. While sliding windows are more accurate and less susceptible to bursts (e.g., a burst of 20 requests at minute 59, and an additional 20 requests at minute 61), this algorithm would have been more complex to implement
   - The requirement that the rate limiting be based on ad duration served instead of a count of requests to the server was a really interesting challenge. I also worked under the assumption that, at a future date, the vast response service would need to know the current ad duration served for a given user, and incorporate that in the ad selection logic so as to not exceed that limit
       - However, in this version of the adserver, the currentAdDurationServed is not used in the ad selection logic
- Unit Testing
    - I implemented a  WIP unit testing structure for the APIs. One of the benefits of using a repository pattern is the ability to mock data for testing purposes, which I took a stab at
    - However, there are conflicting views from resources online that claim that it is in fact better to use a real database for testing
        - I think this actually makes more sense, as the interface methods implemented for the mock repo classes can more closely resemble the real methods, and test the SQL queries 
        - This approach, combined with having a separate, test DB in a docker container, I think would be better than mocking the data

## Implementation Decision Making Hierarchy
1. Leverage past experience from similar projects
2. Check official documentation
3. Use **reputable** StackOverflow posts, guides, and tutorials


Leverage 2 & 3 to find more optimal approaches from 1


## Future Enhancements and Considerations

### Engineering/Dev Related
- Implement rate limiting using Redis cache instead of in-memory
   - Distributed cache allows for rate limiting for users across multiple server instances
- Authentication & Authorization
   - Especially required for managing campaigns
- Database caching for frequently used creatives
   - Use Redis to cache creatives frequently selected for insertion in ad responses
- Explore optimistic concurrency control/locks to avoid race conditions
   - Example: use timestamp/version column
- Log ad impressions, quartile beacons, etc. using event streaming
   - Example: Kafka, RabbitMQ
   - Can also use ELK stack (Elasticsearch, Logstash, Kibana)
- Campaign Active vs. Inactive
   - Run chronjob daily to designate campaign as active/inactive based on date that chronjob is run
   - May optimize ad response process by reducing the number of steps/calculations required

### Adtech Related
- Expand ad request query parameters (e.g., slot parameters)
- Implement external API for geolocation look up using RemoteAddr from adrequest HTTP request header
- Add table to register publisher networks
- Add additional tables mapped to each VAST node
- Wrap VAST file inside of an ad response
   - Allow for multiple ads to be included in a single response, reducing number of requests received
   - Allow for temporal slots (i.e., preroll, midroll, postroll)
   - Allow for slotImpressions
       - slotImpressions can be used to calculate avails vs. filled
   - Or could use the Extension node in VAST for slot beacons/info
- Invalid traffic (IVT) filtering
   - Block traffic from known bots and automated tools
   - Suspicious non-human patterns (e.g., high number of clicks or requests)

   #### Thanks for reading!