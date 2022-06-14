create table if not exists shorten_link_urls (
    id bigint primary key,
    client_id varchar(500),
    short_url varchar(500),
    long_url varchar(500),
    expire_time timestamp
);
create table if not exists shorten_link_clients (
    id varchar (500) primary key,
    api_key varchar(500)
);
alter table shorten_link_urls
add index index_id using btree(id);
alter table shorten_link_urls
add index index_short_url using hash(short_url);
alter table shorten_link_urls
add index index_long_url using hash(long_url);