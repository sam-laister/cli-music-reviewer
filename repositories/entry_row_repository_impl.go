package repositories

import (
	"cli-music-reviewer/models/entities"
	"database/sql"
)

type EntryRowRepositoryImpl struct {
	*EntityRepository[*entities.EntryRow]
}

func NewEntryRowRepository(db *sql.DB) *EntryRowRepositoryImpl {
	return &EntryRowRepositoryImpl{
		EntityRepository: NewEntityRepository[*entities.EntryRow](db),
	}
}

func (r *EntryRowRepositoryImpl) GetActiveRows() ([]*entities.EntryRow, error) {
	return r.EntityRepository.FindBy("active", true)
}

var _ EntryRowRepositoryInterface = (*EntryRowRepositoryImpl)(nil)
