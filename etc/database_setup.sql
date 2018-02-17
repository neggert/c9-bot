create database if not exists c9bot;

use c9bot;

create table if not exists occurrences(
    channelid bigint unsigned,
    ts timestamp
);
