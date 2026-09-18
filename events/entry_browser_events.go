package events

type EntryCreateRequestedMsg struct{}

type EntryEditRequestedMsg struct {
	Index int
}

type EntryCreateSubmittedMsg struct {
	Title          string
	Artist         string
	ReleaseDate    string
	SpotifyID      string
	SpotifyType    string
	SpotifyLink    string
	CoverArtSmall  string
	CoverArtMedium string
	CoverArtLarge  string
}

type EntryCreateCancelledMsg struct{}
