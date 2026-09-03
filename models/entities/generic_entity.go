package entities

type GenericEntity struct {
	ID int
}

func (r *GenericEntity) GetID() int {
	return r.ID
}

func (r *GenericEntity) SetID(id int) {
	r.ID = id
}
