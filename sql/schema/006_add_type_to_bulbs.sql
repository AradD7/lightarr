-- +goose Up
ALTER TABLE bulbs
ADD COLUMN type TEXT NOT NULL DEFAULT 'normal';

-- +goose Down
ALTER TABLE bulbs
DROP COLUMN type;
