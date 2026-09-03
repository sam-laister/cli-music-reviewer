package entities

import (
	"cli-music-reviewer/interfaces"
	"database/sql"
)

type EntryRow struct {
	GenericEntity
	Title          string
	Body           string
	UpdatedAt      string
	CreatedAt      string
	Active         bool
	SpotifyID      string
	SpotifyType    string
	SpotifyLink    string
	CoverArtSmall  string
	CoverArtMedium string
	CoverArtLarge  string
}

func (r *EntryRow) ScanRow(row *sql.Row) error {
	return row.Scan(&r.ID, &r.Title, &r.Body, &r.CreatedAt, &r.UpdatedAt, &r.Active, &r.SpotifyID, &r.SpotifyType, &r.SpotifyLink, &r.CoverArtSmall, &r.CoverArtMedium, &r.CoverArtLarge)
}

func (r *EntryRow) ScanRows(rows *sql.Rows) error {
	return rows.Scan(&r.ID, &r.Title, &r.Body, &r.CreatedAt, &r.UpdatedAt, &r.Active, &r.SpotifyID, &r.SpotifyType, &r.SpotifyLink, &r.CoverArtSmall, &r.CoverArtMedium, &r.CoverArtLarge)
}

func (r *EntryRow) Values() []interface{} {
	return []interface{}{r.Title, r.Body, r.CreatedAt, r.UpdatedAt, r.Active, r.SpotifyID, r.SpotifyType, r.SpotifyLink, r.CoverArtSmall, r.CoverArtMedium, r.CoverArtLarge}
}

func (r *EntryRow) Columns() []string {
	return []string{"title", "body", "created_at", "updated_at", "active", "spotify_id", "spotify_type", "spotify_link", "cover_art_small", "cover_art_medium", "cover_art_large"}
}

func (r *EntryRow) TableName() string {
	return "entry_rows"
}

var _ interfaces.EntityInterface = (*EntryRow)(nil)
