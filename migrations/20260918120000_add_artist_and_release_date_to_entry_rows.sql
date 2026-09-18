-- +goose Up
ALTER TABLE entry_rows ADD COLUMN artist TEXT NOT NULL DEFAULT '';
ALTER TABLE entry_rows ADD COLUMN release_date TEXT NOT NULL DEFAULT '';

-- +goose Down
ALTER TABLE entry_rows DROP COLUMN release_date;
ALTER TABLE entry_rows DROP COLUMN artist;
