CREATE TABLE IF NOT EXISTS campaign(
    id bigserial PRIMARY KEY,
    name varchar(255) NOT NULL,
    start_date varchar(255) NOT NULL,
    end_date varchar(255) NOT NULL,
    target_dma_id bigserial NOT NULL,
    ad_id bigserial NOT NULL,
    ad_name varchar(255) NOT NULL,
    ad_duration bigserial NOT NULL,
    ad_creative_id bigserial NOT NULL,
    ad_creative_url varchar(255) NOT NULL
);

CREATE TABLE IF NOT EXISTS ad_transaction(
    transaction_id bigint PRIMARY KEY NOT NULL,
    ad_request text NOT NULL,
    vast_response text,
    client_dma_id bigserial NOT NULL,
    campaign_id bigserial REFERENCES campaign (id)
);

CREATE TABLE IF NOT EXISTS ad_beacon(
    id bigserial PRIMARY KEY,
    transaction_id bigint NOT NULL REFERENCES ad_transaction (transaction_id),
    beacon_url text NOT NULL,
    beacon_name text NOT NULL
);