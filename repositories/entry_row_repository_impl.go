package repositories

import (
	"cli-music-reviewer/models"
	"database/sql"
	"log"
)

type EntryRowRepositoryImpl struct {
	db *sql.DB
}

func (r *EntryRowRepositoryImpl) GetActiveRows() []models.EntryRow {
	rows, err := r.db.Query("SELECT * FROM entry_rows")
	if err != nil {
		log.Fatalf("Failed to execute query: %v", err)
	}

	defer rows.Close()

	var entry_rows []models.EntryRow

	for rows.Next() {
		var u models.EntryRow
		if err := rows.Scan(&u.ID, &u.Title, &u.Body, &u.CreatedAt, &u.UpdatedAt); err != nil {
			log.Fatalf("Failed to scan row: %v", err)
		}
		entry_rows = append(entry_rows, u)
	}

	return entry_rows
}

func NewEntryRowRepository(db *sql.DB) *EntryRowRepositoryImpl {
	return &EntryRowRepositoryImpl{
		db: db,
	}
}
