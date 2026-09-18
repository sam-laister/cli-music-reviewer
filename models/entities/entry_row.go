package entities

type EntryRow struct {
	GenericEntity
	Title          string `db:"title"`
	Body           string `db:"body"`
	Active         bool   `db:"active"`
	SpotifyID      string `db:"spotify_id"`
	SpotifyType    string `db:"spotify_type"`
	SpotifyLink    string `db:"spotify_link"`
	CoverArtSmall  string `db:"cover_art_small"`
	CoverArtMedium string `db:"cover_art_medium"`
	CoverArtLarge  string `db:"cover_art_large"`
}

func NewEntryRow(title, body, spotifyId, spotifyType, spotifyLink, coverArtSmall, coverArtMedium, coverArtLarge string, active bool) *EntryRow {
	row := &EntryRow{
		Title:          title,
		Body:           body,
		Active:         active,
		SpotifyID:      spotifyId,
		SpotifyType:    spotifyType,
		SpotifyLink:    spotifyLink,
		CoverArtSmall:  coverArtSmall,
		CoverArtMedium: coverArtMedium,
		CoverArtLarge:  coverArtLarge,
	}
	row.MarkCreated()
	row.MarkUpdated()
	return row
}

func (r *EntryRow) TableName() string {
	return "entry_rows"
}

var _ EntityInterface = (*EntryRow)(nil)
