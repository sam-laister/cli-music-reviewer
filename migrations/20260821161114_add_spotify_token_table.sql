-- +goose Up
CREATE TABLE spotify_tokens (id INTEGER PRIMARY KEY, access_token TEXT, refresh_token TEXT, expires_at date, updated_at date);

-- +goose Down
DROP TABLE spotify_tokens;
