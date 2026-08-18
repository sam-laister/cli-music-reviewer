package repositories

import "cli-music-reviewer/models"

type EntryRowRepository interface {
	GetActiveRows() []models.EntryRow
}
