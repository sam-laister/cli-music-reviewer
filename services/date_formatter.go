package services

import (
	"cli-music-reviewer/models/entities"
	"time"
)

func DateToString(date time.Time) string {
	return date.UTC().Format(entities.DEFAULT_DB_DATE_FORMAT)
}

// StringToDate parses a date stored via DateToString. The modernc.org/sqlite
// driver returns TEXT columns as plain strings rather than converting them to
// time.Time on Scan, so callers reading a stored timestamp must parse it
// explicitly.
func StringToDate(value string) (time.Time, error) {
	return time.Parse(entities.DEFAULT_DB_DATE_FORMAT, value)
}
