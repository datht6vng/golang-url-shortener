CREATE TABLE IF NOT EXISTS shorten_link_urls (
    id BIGINT PRIMARY KEY,
    client_id VARCHAR(256),
    short_url VARCHAR(256),
    long_url VARCHAR(500),
    expire_time DATETIME
);

CREATE TABLE IF NOT EXISTS shorten_link_clients (
    client_id VARCHAR (256) PRIMARY KEY,
    api_key VARCHAR(256),
    liciense_key VARCHAR(256),
    max_link BIGINT
);

ALTER TABLE shorten_link_urls
    ADD INDEX index_id USING BTREE (id);
ALTER TABLE shorten_link_urls
    ADD INDEX index_short_url USING HASH (short_url);
ALTER TABLE shorten_link_urls
    ADD INDEX index_long_url USING HASH (long_url);

-- CREATE TABLE IF NOT EXISTS shorten_link_generate_counter (
--     client_id VARCHAR(256),
--     create_date DATE,
--     number_link_generated INT,
--     PRIMARY KEY (client_id, create_date)
-- );

