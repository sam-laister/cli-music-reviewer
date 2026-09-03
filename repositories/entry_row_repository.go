package repositories

import (
	"cli-music-reviewer/models/entities"
)

type EntryRowRepositoryInterface interface {
	EntityRepositoryInterface[*entities.EntryRow]
	GetActiveRows() ([]*entities.EntryRow, error)
}
