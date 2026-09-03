package events

type EntryCreateRequestedMsg struct{}

type EntryEditRequestedMsg struct {
	Index int
}

type EntryCreateSubmittedMsg struct {
	Title string
}

type EntryCreateCancelledMsg struct{}
