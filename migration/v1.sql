CREATE TABLE IF NOT EXISTS shorten_link_urls (
    id BIGINT PRIMARY KEY,
    client_id VARCHAR(256),
    short_url VARCHAR(256),
    long_url TEXT,
    expire_time timestamp
);

CREATE TABLE IF NOT EXISTS shorten_link_clients (
    client_id VARCHAR (256) PRIMARY KEY,
    client_name VARCHAR (256),
    api_key VARCHAR(256)
);

ALTER TABLE shorten_link_urls
    ADD INDEX index_id USING BTREE (client_id);
ALTER TABLE shorten_link_urls
    ADD INDEX index_short_url USING HASH (short_url);
ALTER TABLE shorten_link_urls
    ADD INDEX index_long_url USING HASH (long_url(2048));