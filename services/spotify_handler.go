package services

import "cli-music-reviewer/models/dtos"

type SpotifyHandler interface {
	Authorize() error
	Token(code string) (*TokenResponse, error)
	// EnsureAuthorized checks the stored token, refreshing it if expired.
	// Returns ErrNoStoredToken if no token has ever been stored, meaning the
	// caller must run the full Authorize/callback flow instead.
	EnsureAuthorized() error
	// GetSavedAlbums fetches a page of the current user's saved albums.
	// market may be empty to omit the query parameter.
	GetSavedAlbums(limit, offset int, market string) (*dtos.GetSavedAlbumsResponseDTO, error)
}
