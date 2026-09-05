package services

import (
	"cli-music-reviewer/models/entities"
	"time"
)

func DateToString(date time.Time) string {
	return date.UTC().Format(entities.DEFAULT_DB_DATE_FORMAT)
}
