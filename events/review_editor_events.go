package events

import "cli-music-reviewer/models/entities"

type ReviewSaveRequestedMsg struct {
	Entry *entities.EntryRow
}

type ReviewEditCancelledMsg struct{}
