-- +goose Up
CREATE TABLE players (
    id          INTEGER PRIMARY KEY,
    name        TEXT NOT NULL,
    last_seen   TEXT NOT NULL
);

CREATE TABLE accounts (
    id      INTEGER PRIMARY KEY,
    title   TEXT NOT NULL,
    thumb   TEXT
);
-- +goose Down
DROP TABLE players;
DROP TABLE accounts;
