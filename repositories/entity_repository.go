package repositories

import "cli-music-reviewer/interfaces"

type EntityRepositoryInterface[T interfaces.EntityInterface] interface {
	Create(entity T) (T, error)
	Update(entity T) error
	Delete(id int) error
	FindByID(id int) (T, error)
	FindAll() ([]T, error)
	FindBy(column string, value interface{}) ([]T, error)
	GetLatestOrNull() (T, error)
	Exists(id int) (bool, error)
	Count() (int, error)
}
