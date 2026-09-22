package reviews

type entryArtworkLoadedMsg struct {
	requestID int
	url       string
	rendered  string
	err       error
}
