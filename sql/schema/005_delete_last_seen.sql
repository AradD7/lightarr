-- +goose Up
ALTER TABLE devices
DROP COLUMN last_seen;

-- +goose Down
ALTER TABLE devices
ADD COLUMN last_seen TEXT NOT NULL;
