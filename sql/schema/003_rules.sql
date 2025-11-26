-- +goose Up
CREATE TABLE rules (
    id          TEXT PRIMARY KEY,
    condition   TEXT NOT NULL,
    action      TEXT NOT NULL
);

-- +goose Down
DROP TABLE rules;
