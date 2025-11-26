-- +goose Up
CREATE TABLE players (
    uuid            TEXT PRIMARY KEY,
    title           TEXT NOT NULL,
    public_address  TEXT NOT NULL
);

CREATE TABLE accounts (
    id      INTEGER PRIMARY KEY,
    title   TEXT NOT NULL
);
-- +goose Down
DROP TABLE players;
DROP TABLE accounts;
