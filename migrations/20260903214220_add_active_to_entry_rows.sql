-- +goose Up
ALTER TABLE entry_rows ADD COLUMN active BOOLEAN NOT NULL DEFAULT TRUE;

-- +goose Down
ALTER TABLE entry_rows DROP COLUMN active;