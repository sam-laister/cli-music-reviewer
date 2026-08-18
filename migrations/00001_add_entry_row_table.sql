-- +goose Up
CREATE TABLE entry_rows (id INTEGER PRIMARY KEY, title TEXT, body TEXT, created_at date, updated_at date);

-- +goose Down
DROP TABLE entry_rows;
