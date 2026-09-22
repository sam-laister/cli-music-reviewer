package options

// spotifyAuthCheckedMsg carries the result of an EnsureAuthorized call made
// in response to a refresh request.
type spotifyAuthCheckedMsg struct {
	err error
}

// spotifyConnectedMsg carries the result of a full Authorize/callback flow
// started in response to an authorize request.
type spotifyConnectedMsg struct {
	err error
}
