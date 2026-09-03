package entities

import (
	"cli-music-reviewer/interfaces"
	"database/sql"
	"time"
)

type SpotifyToken struct {
	GenericEntity
	AccessToken  string
	RefreshToken string
	ExpiresAt    time.Time
	UpdatedAt    time.Time
}

func (t *SpotifyToken) TableName() string {
	return "spotify_tokens"
}

func (t *SpotifyToken) Columns() []string {
	return []string{"access_token", "refresh_token", "expires_at", "updated_at"}
}

func (t *SpotifyToken) Values() []interface{} {
	return []interface{}{t.AccessToken, t.RefreshToken, t.ExpiresAt, t.UpdatedAt}
}

func (t *SpotifyToken) ScanRow(row *sql.Row) error {
	return row.Scan(&t.ID, &t.AccessToken, &t.RefreshToken, &t.ExpiresAt, &t.UpdatedAt)
}

func (t *SpotifyToken) ScanRows(rows *sql.Rows) error {
	return rows.Scan(&t.ID, &t.AccessToken, &t.RefreshToken, &t.ExpiresAt, &t.UpdatedAt)
}

var _ interfaces.EntityInterface = (*SpotifyToken)(nil)
