package services

import "errors"

// Connect ensures a valid Spotify token is available: it reuses or refreshes
// a stored token via EnsureAuthorized, falling back to the full browser
// Authorize/callback flow only when no token has ever been stored.
func Connect(spotifyHandler SpotifyHandler, httpHandler HttpHandler) error {
	err := spotifyHandler.EnsureAuthorized()
	if err == nil {
		return nil
	}
	if !errors.Is(err, ErrNoStoredToken) {
		return err
	}

	if err := httpHandler.Setup(); err != nil {
		return err
	}

	if err := spotifyHandler.Authorize(); err != nil {
		return err
	}

	return httpHandler.Wait()
}
