package services

type SpotifyHandler interface {
	Authorize() error
	Token(code string) (*TokenResponse, error)
	// EnsureAuthorized checks the stored token, refreshing it if expired.
	// Returns ErrNoStoredToken if no token has ever been stored, meaning the
	// caller must run the full Authorize/callback flow instead.
	EnsureAuthorized() error
}
