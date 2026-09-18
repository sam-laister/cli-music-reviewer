package dtos

// GetSavedAlbumsResponseDTO is the response body for
// GET https://api.spotify.com/v1/me/albums (a paged list of the current
// user's saved albums).
type GetSavedAlbumsResponseDTO struct {
	Href     string          `json:"href"`
	Limit    int             `json:"limit"`
	Next     *string         `json:"next"`
	Offset   int             `json:"offset"`
	Previous *string         `json:"previous"`
	Total    int             `json:"total"`
	Items    []SavedAlbumDTO `json:"items"`
}

type SavedAlbumDTO struct {
	AddedAt string   `json:"added_at"`
	Album   AlbumDTO `json:"album"`
}

type AlbumDTO struct {
	AlbumType            string                 `json:"album_type"`
	TotalTracks          int                    `json:"total_tracks"`
	AvailableMarkets     []string               `json:"available_markets"`
	ExternalURLs         ExternalURLsDTO        `json:"external_urls"`
	Href                 string                 `json:"href"`
	ID                   string                 `json:"id"`
	Images               []ImageDTO             `json:"images"`
	Name                 string                 `json:"name"`
	ReleaseDate          string                 `json:"release_date"`
	ReleaseDatePrecision string                 `json:"release_date_precision"`
	Restrictions         *AlbumRestrictionsDTO  `json:"restrictions,omitempty"`
	Type                 string                 `json:"type"`
	URI                  string                 `json:"uri"`
	Artists              []SimplifiedArtistDTO  `json:"artists"`
	Tracks               SimplifiedTrackPageDTO `json:"tracks"`
	Copyrights           []CopyrightDTO         `json:"copyrights"`
	ExternalIDs          ExternalIDsDTO         `json:"external_ids"`
	Genres               []string               `json:"genres"`
	Label                string                 `json:"label"`
	Popularity           int                    `json:"popularity"`
}

type AlbumRestrictionsDTO struct {
	Reason string `json:"reason"`
}

type ExternalURLsDTO struct {
	Spotify string `json:"spotify"`
}

type ExternalIDsDTO struct {
	ISRC string `json:"isrc"`
	EAN  string `json:"ean"`
	UPC  string `json:"upc"`
}

type ImageDTO struct {
	URL    string `json:"url"`
	Height *int   `json:"height"`
	Width  *int   `json:"width"`
}

type CopyrightDTO struct {
	Text string `json:"text"`
	Type string `json:"type"`
}

type SimplifiedArtistDTO struct {
	ExternalURLs ExternalURLsDTO `json:"external_urls"`
	Href         string          `json:"href"`
	ID           string          `json:"id"`
	Name         string          `json:"name"`
	Type         string          `json:"type"`
	URI          string          `json:"uri"`
}

// SimplifiedTrackPageDTO is the paged list of a saved album's tracks.
type SimplifiedTrackPageDTO struct {
	Href     string               `json:"href"`
	Limit    int                  `json:"limit"`
	Next     *string              `json:"next"`
	Offset   int                  `json:"offset"`
	Previous *string              `json:"previous"`
	Total    int                  `json:"total"`
	Items    []SimplifiedTrackDTO `json:"items"`
}

type SimplifiedTrackDTO struct {
	Artists          []SimplifiedArtistDTO `json:"artists"`
	AvailableMarkets []string              `json:"available_markets"`
	DiscNumber       int                   `json:"disc_number"`
	DurationMs       int                   `json:"duration_ms"`
	Explicit         bool                  `json:"explicit"`
	ExternalURLs     ExternalURLsDTO       `json:"external_urls"`
	Href             string                `json:"href"`
	ID               string                `json:"id"`
	IsPlayable       bool                  `json:"is_playable,omitempty"`
	Name             string                `json:"name"`
	PreviewURL       *string               `json:"preview_url"`
	TrackNumber      int                   `json:"track_number"`
	Type             string                `json:"type"`
	URI              string                `json:"uri"`
	IsLocal          bool                  `json:"is_local"`
}
