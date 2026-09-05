package entities

import (
	"cli-music-reviewer/interfaces"
	"database/sql"
	"fmt"
	"time"
)

type EntryRow struct {
	GenericEntity
	Title          string
	Body           string
	Active         bool
	SpotifyID      string
	SpotifyType    string
	SpotifyLink    string
	CoverArtSmall  string
	CoverArtMedium string
	CoverArtLarge  string
}

func NewEntryRow(title, body, spotifyId, spotifyType, spotifyLink, coverArtSmall, coverArtMedium, coverArtLarge string, active bool) *EntryRow {
	return &EntryRow{
		Title:          title,
		Body:           body,
		Active:         active,
		SpotifyID:      spotifyId,
		SpotifyType:    spotifyType,
		SpotifyLink:    spotifyLink,
		CoverArtSmall:  coverArtSmall,
		CoverArtMedium: coverArtMedium,
		CoverArtLarge:  coverArtLarge,
		GenericEntity: GenericEntity{
			UpdatedAt: time.Now(),
			CreatedAt: time.Now(),
		},
	}
}

func (r *EntryRow) ScanRow(row *sql.Row) error {
	var createdAt, updatedAt string
	if err := row.Scan(&r.ID, &r.Title, &r.Body, &createdAt, &updatedAt, &r.Active, &r.SpotifyID, &r.SpotifyType, &r.SpotifyLink, &r.CoverArtSmall, &r.CoverArtMedium, &r.CoverArtLarge); err != nil {
		return err
	}
	return r.parseTimestamps(createdAt, updatedAt)
}

func (r *EntryRow) ScanRows(rows *sql.Rows) error {
	var createdAt, updatedAt string
	if err := rows.Scan(&r.ID, &r.Title, &r.Body, &createdAt, &updatedAt, &r.Active, &r.SpotifyID, &r.SpotifyType, &r.SpotifyLink, &r.CoverArtSmall, &r.CoverArtMedium, &r.CoverArtLarge); err != nil {
		return err
	}
	return r.parseTimestamps(createdAt, updatedAt)
}

// parseTimestamps converts the TEXT-column timestamps returned by
// modernc.org/sqlite (it does not auto-convert to time.Time on Scan, unlike
// mattn/go-sqlite3) into the CreatedAt/UpdatedAt fields.
func (r *EntryRow) parseTimestamps(createdAt, updatedAt string) error {
	created, err := time.Parse(DEFAULT_DB_DATE_FORMAT, createdAt)
	if err != nil {
		return fmt.Errorf("parsing created_at: %w", err)
	}
	updated, err := time.Parse(DEFAULT_DB_DATE_FORMAT, updatedAt)
	if err != nil {
		return fmt.Errorf("parsing updated_at: %w", err)
	}
	r.CreatedAt = created
	r.UpdatedAt = updated
	return nil
}

func (r *EntryRow) Values() []interface{} {
	return []interface{}{r.Title, r.Body, r.CreatedAt.Format(DEFAULT_DB_DATE_FORMAT), r.UpdatedAt.Format(DEFAULT_DB_DATE_FORMAT), r.Active, r.SpotifyID, r.SpotifyType, r.SpotifyLink, r.CoverArtSmall, r.CoverArtMedium, r.CoverArtLarge}
}

func (r *EntryRow) Columns() []string {
	return []string{"title", "body", "created_at", "updated_at", "active", "spotify_id", "spotify_type", "spotify_link", "cover_art_small", "cover_art_medium", "cover_art_large"}
}

func (r *EntryRow) TableName() string {
	return "entry_rows"
}

var _ interfaces.EntityInterface = (*EntryRow)(nil)
