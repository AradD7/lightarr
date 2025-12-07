-- +goose Up
CREATE TABLE devices (
    id          TEXT PRIMARY KEY,
    name        TEXT NOT NULL
);

CREATE TABLE accounts (
    id      INTEGER PRIMARY KEY,
    title   TEXT NOT NULL,
    thumb   TEXT
);
-- +goose Down
DROP TABLE devices;
DROP TABLE accounts;
