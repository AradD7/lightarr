-- +goose Up
ALTER TABLE devices
ADD COLUMN product TEXT NOT NULL;

-- +goose Down
ALTER TABLE devices
DROP COLUMN product;
