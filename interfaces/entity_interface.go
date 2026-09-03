package interfaces

import "database/sql"

type EntityInterface interface {
	SetID(int)
	GetID() int
	TableName() string
	ScanRow(*sql.Row) error
	ScanRows(*sql.Rows) error
	Values() []interface{} // in column order
	Columns() []string     // column names to provide sync between values
}
