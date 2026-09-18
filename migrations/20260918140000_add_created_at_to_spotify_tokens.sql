-- +goose Up
ALTER TABLE spotify_tokens ADD COLUMN created_at date;

-- +goose Down
ALTER TABLE spotify_tokens DROP COLUMN created_at;
