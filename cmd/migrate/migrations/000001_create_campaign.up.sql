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