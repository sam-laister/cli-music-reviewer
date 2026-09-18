package repositories

import (
	"cli-music-reviewer/models/entities"
	"database/sql"
	"fmt"
	"reflect"
	"strings"

	"github.com/jmoiron/sqlx"
)

type EntityRepositoryImpl[T entities.EntityInterface] struct {
	db *sqlx.DB
}

func NewEntityRepositoryImpl[T entities.EntityInterface](db *sqlx.DB) *EntityRepositoryImpl[T] {
	return &EntityRepositoryImpl[T]{db: db}
}

func newEntity[T entities.EntityInterface]() T {
	var zero T
	t := reflect.TypeOf(zero)
	if t == nil || t.Kind() != reflect.Pointer {
		return zero
	}
	return reflect.New(t.Elem()).Interface().(T)
}

// dbColumns returns the `db`-tagged column names for T, excluding "id",
// by walking the struct fields (embedded fields included).
func dbColumns[T entities.EntityInterface]() []string {
	t := reflect.TypeOf(newEntity[T]()).Elem()

	var cols []string
	var walk func(reflect.Type)
	walk = func(t reflect.Type) {
		for i := 0; i < t.NumField(); i++ {
			field := t.Field(i)
			if field.Anonymous {
				walk(field.Type)
				continue
			}
			tag := field.Tag.Get("db")
			if tag == "" || tag == "-" || tag == "id" {
				continue
			}
			cols = append(cols, tag)
		}
	}
	walk(t)
	return cols
}

func (r *EntityRepositoryImpl[T]) FindByID(id uint64) (T, error) {
	result := newEntity[T]()
	query := fmt.Sprintf("SELECT * FROM %s WHERE id = ?", result.TableName())
	err := r.db.Get(result, query, id)
	return result, err
}

func (r *EntityRepositoryImpl[T]) Create(entity T) (T, error) {
	cols := dbColumns[T]()
	placeholders := make([]string, len(cols))
	for i, c := range cols {
		placeholders[i] = ":" + c
	}
	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)",
		entity.TableName(),
		strings.Join(cols, ", "),
		strings.Join(placeholders, ", "))
	result, err := r.db.NamedExec(query, entity)
	if err != nil {
		return entity, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return entity, err
	}

	entity.SetID(uint64(id))
	return entity, nil
}

func (r *EntityRepositoryImpl[T]) GetLatestOrNull() (T, error) {
	var zero T
	query := fmt.Sprintf("SELECT * FROM %s ORDER BY id DESC LIMIT 1", zero.TableName())

	result := newEntity[T]()
	if err := r.db.Get(result, query); err != nil {
		if err == sql.ErrNoRows {
			return zero, nil
		}
		return zero, err
	}
	return result, nil
}

func (r *EntityRepositoryImpl[T]) FindBy(column string, value any) ([]T, error) {
	var zero T
	query := fmt.Sprintf("SELECT * FROM %s WHERE %s = ?", zero.TableName(), column)

	var results []T
	err := r.db.Select(&results, query, value)
	return results, err
}

func (r *EntityRepositoryImpl[T]) FindOneBy(column string, value any) (T, error) {
	result := newEntity[T]()
	query := fmt.Sprintf("SELECT * FROM %s WHERE %s = ? LIMIT 1", result.TableName(), column)
	err := r.db.Get(result, query, value)
	return result, err
}

func (r *EntityRepositoryImpl[T]) Exists(id uint64) (bool, error) {
	var zero T
	query := fmt.Sprintf("SELECT 1 FROM %s WHERE id = ? LIMIT 1", zero.TableName())
	var dummy int
	err := r.db.QueryRow(query, id).Scan(&dummy)
	if err == sql.ErrNoRows {
		return false, nil
	}
	return err == nil, err
}

func (r *EntityRepositoryImpl[T]) Count() (int, error) {
	var zero T
	query := fmt.Sprintf("SELECT COUNT(*) FROM %s", zero.TableName())
	var count int
	err := r.db.QueryRow(query).Scan(&count)
	return count, err
}

func (r *EntityRepositoryImpl[T]) Delete(id uint64) error {
	var zero T
	query := fmt.Sprintf("DELETE FROM %s WHERE id = ?", zero.TableName())
	_, err := r.db.Exec(query, id)
	return err
}

func (r *EntityRepositoryImpl[T]) FindAll() ([]T, error) {
	var zero T
	query := fmt.Sprintf("SELECT * FROM %s", zero.TableName())

	var results []T
	err := r.db.Select(&results, query)
	return results, err
}

func (r *EntityRepositoryImpl[T]) Update(entity T) error {
	cols := dbColumns[T]()
	setClauses := make([]string, len(cols))
	for i, c := range cols {
		setClauses[i] = fmt.Sprintf("%s = :%s", c, c)
	}

	query := fmt.Sprintf("UPDATE %s SET %s WHERE id = :id",
		entity.TableName(),
		strings.Join(setClauses, ", "))

	_, err := r.db.NamedExec(query, entity)
	return err
}
