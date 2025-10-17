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