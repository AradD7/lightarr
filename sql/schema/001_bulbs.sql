-- +goose Up
CREATE TABLE bulbs (
    mac         TEXT PRIMARY KEY,
    created_at  TIMESTAMP NOT NULL,
    updated_at  TIMESTAMP NOT NULL,
    ip          TEXT NOT NULL UNIQUE,
    name        TEXT NOT NULL
);

-- +goose Down
DROP TABLE bulbs;
