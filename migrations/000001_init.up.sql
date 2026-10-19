CREATE SCHEMA urlShortenerApp;

CREATE TABLE urlShortenerApp.urls (
    id        SERIAL                                    PRIMARY KEY,
    full_url  VARCHAR(1000)     NOT NULL    UNIQUE,
    code      VARCHAR(10)       NOT NULL    UNIQUE      CHECK(code ~ '^[a-zA-Z0-9_]{10}$')
);