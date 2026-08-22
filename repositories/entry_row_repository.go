package repositories

import (
	"cli-music-reviewer/models/entities"
)

type EntryRowRepository interface {
	GetActiveRows() []entities.EntryRow
}
