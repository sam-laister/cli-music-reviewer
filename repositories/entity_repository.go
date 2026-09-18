package repositories

import "cli-music-reviewer/models/entities"

type EntityRepositoryInterface[T entities.EntityInterface] interface {
	Create(entity T) (T, error)
	Update(entity T) error
	Delete(id uint64) error
	FindByID(id uint64) (T, error)
	FindAll() ([]T, error)
	FindBy(column string, value any) ([]T, error)
	FindOneBy(column string, value any) (T, error)
	GetLatestOrNull() (T, error)
	Exists(id uint64) (bool, error)
	Count() (int, error)
}
