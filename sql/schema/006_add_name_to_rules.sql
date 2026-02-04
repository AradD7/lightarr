
-- +goose Up
ALTER TABLE rules
ADD COLUMN name TEXT;

-- +goose Down
ALTER TABLE rules
DROP COLUMN name;
