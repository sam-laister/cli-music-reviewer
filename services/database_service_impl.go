package services

import "database/sql"

const (
	driverName = "sqlite"
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
