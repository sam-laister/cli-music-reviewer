package services

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

const (
	driverName = "sqlite3"
)

type DatabaseServiceImpl struct {
	dsnURI string
}

func (s *DatabaseServiceImpl) Connect() (*sqlx.DB, error) {
	return sqlx.Open(driverName, s.dsnURI)
}

func NewDatabaseService(dsnURI string) *DatabaseServiceImpl {
	return &DatabaseServiceImpl{
		dsnURI: dsnURI,
	}
}
