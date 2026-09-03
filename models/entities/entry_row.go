package entities

import (
	"cli-music-reviewer/interfaces"
	"database/sql"
)

type EntryRow struct {
	GenericEntity
	Title     string
	Body      string
	UpdatedAt string
	CreatedAt string
	Active    bool
}

func (r *EntryRow) ScanRow(row *sql.Row) error {
	return row.Scan(&r.ID, &r.Title, &r.Body, &r.CreatedAt, &r.UpdatedAt, &r.Active)
}

func (r *EntryRow) ScanRows(rows *sql.Rows) error {
	return rows.Scan(&r.ID, &r.Title, &r.Body, &r.CreatedAt, &r.UpdatedAt, &r.Active)
}

func (r *EntryRow) Values() []interface{} {
	return []interface{}{r.Title, r.Body, r.CreatedAt, r.UpdatedAt, r.Active}
}

func (r *EntryRow) Columns() []string {
	return []string{"title", "body", "created_at", "updated_at", "active"}
}

func (r *EntryRow) TableName() string {
	return "entry_rows"
}

var _ interfaces.EntityInterface = (*EntryRow)(nil)
