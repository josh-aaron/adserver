# adserver

Welcome to the README for the adserver project!

## Tech Stack, Dependencies
- Go 1.25.2
- Postgres v1.10.9
- github.com/joho/godotenv 
- github.com/rs/vast
- Migrate CLI Tool

## Verification Steps
1. Clone repo
2. Download/install Postgres
3. Create a .env file in the main directory (can copy .env.SAMPLE) and update with your Postgres username, password, and preferred port for the HTTP server
4. Create the adserver DB
    First connect to the psql shell
    ```
    psql -U {username}
    ```
    Or run the script 
    ```
    \i scripts/init_db.sql
    ```
5. Run the migrations to set up the tables
6. Run ```make run``` in the terminal to build and run the executable
7. Use the curl commands in curl-commands.txt to create some campaign data, or leverage an API testing tool like Postman or Thunderclient
8. Test the other available Campaign API endpoints
9. Use the adrequest endpoint to retrieve a VAST response
10. Enter the VAST XML into an online VAST validator (e.g., https://tools.springserve.com/tagtest) and watch the ad play!
11. Run ```make test``` to run the unit tests

## Assumptions and Limitations

- Only one publisher network (e.g., BobsAwesomeCatVideos.com is our only ad space publisher). Therefore, when any ad request is received, we only filter by targetDma and by if campaign is active
- DMA code/id will be hardcoded into ad request
- Ad and Creative data will live in the Campaign table. This is for ease of testing and updating using the exposed Campaign API endpoints. Testers can add/modify campaigns and see those changes reflected in the VAST response.
    - All other VAST data will be hardcoded.
    - The VAST response returned will only contain nodes outlined in the VAST 3.0 Inline Linear Example (https://github.com/InteractiveAdvertisingBureau/VAST_Samples/blob/master/VAST%203.0%20Samples/Inline_Linear_Tag-test.xml)
    - Excludes certain optional nodes (e.g., Wrapper, Non-Linear, Companion).
- If DMA in request does not match any active campaign target DMA, no ad will return
- A campaign will only have one DMA, 
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