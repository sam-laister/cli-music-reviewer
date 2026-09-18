package config

const (
	AppTitle = "CLI Music Reviewer"

	// SpotifyCallbackPort is the local port the OAuth callback server binds
	// to during the Spotify authorize flow; it must match the redirect_uri
	// registered with the Spotify app.
	SpotifyCallbackPort = "8888"
)
