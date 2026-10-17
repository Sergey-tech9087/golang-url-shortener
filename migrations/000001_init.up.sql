CREATE SCHEMA urlShortenerApp;

CREATE TABLE urlShortenerApp.urls (
    id              SERIAL                      PRIMARY KEY,
    full_url        VARCHAR(100)    UNIQUE      NOT NULL        CHECK (char_length(full_url) = 10),
    short_url       VARCHAR(100)                NOT NULL
);