package services

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

const (
	driverName = "sqlite3"
)

type DatabaseServiceImpl struct {
	dsnURI string
}

func (s *DatabaseServiceImpl) Connect() (*sql.DB, error) {
	return sql.Open(driverName, s.dsnURI)
}

func NewDatabaseService(dsnURI string) *DatabaseServiceImpl {
	return &DatabaseServiceImpl{
		dsnURI: dsnURI,
	}
}
