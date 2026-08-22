package repositories

import (
	"cli-music-reviewer/models/entities"
	"database/sql"
	"log"
)

type EntryRowRepositoryImpl struct {
	db *sql.DB
}

func (r *EntryRowRepositoryImpl) GetActiveRows() []entities.EntryRow {
	rows, err := r.db.Query("SELECT * FROM entry_rows")
	if err != nil {
		log.Fatalf("Failed to execute query: %v", err)
	}

	defer rows.Close()

	var entryRows []entities.EntryRow

	for rows.Next() {
		var u entities.EntryRow
		if err := rows.Scan(&u.ID, &u.Title, &u.Body, &u.CreatedAt, &u.UpdatedAt); err != nil {
			log.Fatalf("Failed to scan row: %v", err)
		}
		entryRows = append(entryRows, u)
	}

	return entryRows
}

func NewEntryRowRepository(db *sql.DB) *EntryRowRepositoryImpl {
	return &EntryRowRepositoryImpl{
		db: db,
	}
}
