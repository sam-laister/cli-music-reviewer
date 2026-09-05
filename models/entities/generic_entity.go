package entities

import "time"

type GenericEntity struct {
	ID        int
	UpdatedAt time.Time
	CreatedAt time.Time
}

func (r *GenericEntity) GetID() int {
	return r.ID
}

func (r *GenericEntity) SetID(id int) {
	r.ID = id
}
