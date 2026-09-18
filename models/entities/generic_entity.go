package entities

import "time"

type GenericEntity struct {
	ID        uint64    `db:"id"`
	UpdatedAt time.Time `db:"updated_at"`
	CreatedAt time.Time `db:"created_at"`
}

func (r *GenericEntity) GetID() uint64 {
	return r.ID
}

func (r *GenericEntity) SetID(id uint64) {
	r.ID = id
}

func (r *GenericEntity) MarkUpdated() {
	r.UpdatedAt = time.Now()
}

func (r *GenericEntity) MarkCreated() {
	r.CreatedAt = time.Now()
}

func (r *GenericEntity) GetUpdatedAt() time.Time {
	return r.CreatedAt
}

func (r *GenericEntity) GetCreatedAt() time.Time {
	return r.CreatedAt
}
