-- +goose Up
ALTER TABLE entry_rows ADD COLUMN spotify_id TEXT NOT NULL DEFAULT '';
ALTER TABLE entry_rows ADD COLUMN spotify_type TEXT NOT NULL DEFAULT '';
ALTER TABLE entry_rows ADD COLUMN spotify_link TEXT NOT NULL DEFAULT '';
ALTER TABLE entry_rows ADD COLUMN cover_art_small TEXT NOT NULL DEFAULT '';
ALTER TABLE entry_rows ADD COLUMN cover_art_medium TEXT NOT NULL DEFAULT '';
ALTER TABLE entry_rows ADD COLUMN cover_art_large TEXT NOT NULL DEFAULT '';

-- +goose Down
ALTER TABLE entry_rows DROP COLUMN cover_art_large;
ALTER TABLE entry_rows DROP COLUMN cover_art_medium;
ALTER TABLE entry_rows DROP COLUMN cover_art_small;
ALTER TABLE entry_rows DROP COLUMN spotify_link;
ALTER TABLE entry_rows DROP COLUMN spotify_type;
ALTER TABLE entry_rows DROP COLUMN spotify_id;
