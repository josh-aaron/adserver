# adserver


Welcome to the README for the adserver project! Please also refer to the WORKLOG file, which maps each git commit to a high level unit of work/implementation step.


## Tech Stack & Dependencies
- Go 1.25.2
- Postgres v1.10.9
- github.com/joho/godotenv
- github.com/rs/vast
- golang-migrate (CLI Tool)


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


7. Test the API endpoints!
   - Use the curl commands in curl-commands.txt to create some campaign data, or leverage an API testing tool like Postman or Thunderclient and snag the campaign JSON objects in curl-commands.txt
   - See the "API Documentation" section below for details
   - You can find market names mapped to DMA codes here: https://www.spstechnical.com/DMACodes.htm
8. Test the other available Campaign API endpoints
9. Use the ad response endpoint to retrieve a VAST response
10. Enter the VAST XML into an online VAST validator (e.g., https://tools.springserve.com/tagtest) and watch the ad play!
11. Run ```make test``` to run the unit tests


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
- IP address from HTTP request will be used to identify "unique users" for rate limiting


## Guiding Principles:
1. Use clean, layered architecture / separation of controls
2. Leverage dependency inversion principle, promote loose coupling
3. Strike a balance between familiar tech and learning opportunities


## Implementation Decisions, Considerations, Tradeoffs (WIP)


- Tech stack
   - Goal was to use as few external packages/dependencies as possible
   - I chose Postgres instead of MySQL simply because I haven't worked with it before and wanted a new learning experience
- Architecture/project set up
   - Repository Pattern: Transport, Service, Repository Layer
       - The transport layer handles the API routing and project initialization logic. It is completely separate and independent from the business and database logic
       - Abstracting the database allows us to defer the choice of specific database to a later point, and allows for more efficient testing with mocking
       - Inspiration: https://threedots.tech/post/repository-pattern-in-go/
   - I ended up "combining" the service layer into the repository layer. However, based on the complexity of the vast response service logic, it may have made more sense to separate out the business logic from the model
   - While I had to spend extra time during the planning phase of the project to better understand the repository pattern and get more comfortable with dependency injection, I believe the time investment was worth it to make the application more extensible, flexible, and testable
- API routing
   - To align with my goal of minimizing external dependencies, I chose to stick with the net/http package from the Go standard library, and implement any middleware manually (i.e., rateLimiter)
   - Another good option would have been to use Chi (https://go-chi.io/). Despite being an external package, it is extremely popular. If I had ended up implementing more middleware (e.g., authentication, structured logging) it may have made sense to implement to improve code clarity
- Database tables
   - Here is where I must explain my most "controversial" decision: combining the ad and creative data with the campaign data into a single table.
   - In a real world scenario, these would be segregated into different tables. However, given: a. the limited time frame of the project, and b. the need for  the requirements to be verified/tested in 10 - 20 minutes, I decided that it would be best to allow the tester to modify an ad/creative data using the Campaign API endpoints.
   - A better approach, given more time for implementation and testing, would have been to keep the ad and creative data in its own table and expose another set of API endpoints for the ad/creatives. Then, when it comes to query, use a foreign key to join the campaign and ad/creative tables
- Campaign API
   - Went with standard CRUD operation APIs, while ensuring that no business logic is contained in the transport layer (i.e., the api folder within the main package)
- Ad Response API
   - I decided to use url query parameters for the API endpoint as it would be a good foundation for a future, more developed ad request with numerous query parameters (e.g., parameters for network profile, video asset dimensions, capability flags, etc.)
   - Due to the lack of specificity in the project requirements and the limited timeframe, I decided to keep the targeting and ad selection logic extremely simple, and made several assumptions as long as they did not conflict with the project requirements (see the Assumptions & Limitations section for details)
   - One such assumption I implemented was that each campaign would only have one target DMA, and that the target DMAs would not repeat. Therefore, the logic for looking up campaigns/ads by target DMA is a simple "getByDma" method.
       - However, if I was working under the assumption that a single DMA code could be associated with multiple campaigns, I would have instead implemented a "getAllByDma" method, that returns all campaigns with that DMA code.
       - Then, combined with an updated assumption that a VAST response can contain more than one ad, I would implement a method that returns a slice of type Ad structs, containing one Ad struct for each campaign/ad returned from the DB. Then that slice would be passed to the method(s) that construct the VAST response
- Rate Limiting
   - I decided to implement a fixed-window algorithm for the rate limiter
       - Another option would have been to use a sliding window algorithm. While sliding windows are more accurate and less susceptible to bursts (e.g., a burst of 20 requests at minute 59, and an additional 20 requests at minute 61), this algorithm would have been more complex to implement
   - The requirement that the rate limiting be based on ad duration served instead of a count of requests to the server was a really interesting challenge. I also worked under the assumption that, at a future date, the vast response construction service would need to know the current ad duration served for a given user, and incorporate that in the ad selection logic so as to not exceed that limit
       - However, in this version of the adserver, the currentAdDurationServed is not used in the ad selection logic
- Unit Testing
    - I implemented a yet to be finished set of unit tests for the APIs. One of the benefits of using a repository pattern is the ability to mock data for testing purposes, which I took a stab at
    - However, there are conflicting views from resources online that claim that it is in fact better to use a real database for testing
        - I think this actually makes more sense, as the interface methods implemented for the mock repo classes can more closely resemble the real methods, and test the SQL queries 
        - This approach, combined with having a separate, test DB in a docker container, I think would be better than mocking the data

## Implementation Decision Making Hierarchy
1. Leverage past experience from similar projects
2. Check official documentation
3. Use **reputable** StackOverflow posts, guides, and tutorials


Leverage 2 & 3 to find more optimal approaches from 1


## Future Enhancements and Considerations


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
