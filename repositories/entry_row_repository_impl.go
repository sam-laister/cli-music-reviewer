package repositories

import (
	"cli-music-reviewer/models/entities"

	"github.com/jmoiron/sqlx"
)

type EntryRowRepositoryImpl struct {
	*EntityRepositoryImpl[*entities.EntryRow]
}

func NewEntryRowRepository(db *sqlx.DB) *EntryRowRepositoryImpl {
	return &EntryRowRepositoryImpl{
		EntityRepositoryImpl: NewEntityRepositoryImpl[*entities.EntryRow](db),
	}
}

func (r *EntryRowRepositoryImpl) GetActiveRows() ([]*entities.EntryRow, error) {
	return r.FindBy("active", true)
}

var _ EntryRowRepositoryInterface = (*EntryRowRepositoryImpl)(nil)
