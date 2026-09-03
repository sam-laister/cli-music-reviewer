package repositories

import (
	"cli-music-reviewer/interfaces"
	"database/sql"
	"fmt"
	"strings"
)

type EntityRepository[T interfaces.EntityInterface] struct {
	db *sql.DB
}

func NewEntityRepository[T interfaces.EntityInterface](db *sql.DB) *EntityRepository[T] {
	return &EntityRepository[T]{db: db}
}

func (r *EntityRepository[T]) FindByID(id int) (T, error) {
	var result T
	query := fmt.Sprintf("SELECT * FROM %s WHERE id = ?", result.TableName())

	row := r.db.QueryRow(query, id)
	err := result.ScanRow(row)
	return result, err
}

func (r *EntityRepository[T]) Create(entity T) (T, error) {
	cols := entity.Columns()
	placeholders := make([]string, len(cols))
	for i := range cols {
		placeholders[i] = "?"
	}
	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)",
		entity.TableName(),
		strings.Join(cols, ", "),
		strings.Join(placeholders, ", "))
	result, err := r.db.Exec(query, entity.Values()...)

	id, err := result.LastInsertId()
	if err != nil {
		return entity, err
	}

	entity.SetID(int(id))
	return entity, err
}

func (r *EntityRepository[T]) GetLatestOrNull() (T, error) {
	var zero T
	query := fmt.Sprintf("SELECT * FROM %s ORDER BY id DESC LIMIT 1", zero.TableName())

	row := r.db.QueryRow(query)
	var result T
	if err := result.ScanRow(row); err != nil {
		if err == sql.ErrNoRows {
			return zero, nil
		}
		return zero, err
	}
	return result, nil
}

func (r *EntityRepository[T]) FindBy(column string, value interface{}) ([]T, error) {
	var zero T
	query := fmt.Sprintf("SELECT * FROM %s WHERE %s = ?", zero.TableName(), column)

	rows, err := r.db.Query(query, value)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []T
	for rows.Next() {
		var item T
		if err := item.ScanRows(rows); err != nil {
			return nil, err
		}
		results = append(results, item)
	}
	return results, rows.Err()
}

func (r *EntityRepository[T]) Exists(id int) (bool, error) {
	var zero T
	query := fmt.Sprintf("SELECT 1 FROM %s WHERE id = ? LIMIT 1", zero.TableName())
	var dummy int
	err := r.db.QueryRow(query, id).Scan(&dummy)
	if err == sql.ErrNoRows {
		return false, nil
	}
	return err == nil, err
}

func (r *EntityRepository[T]) Count() (int, error) {
	var zero T
	query := fmt.Sprintf("SELECT COUNT(*) FROM %s", zero.TableName())
	var count int
	err := r.db.QueryRow(query).Scan(&count)
	return count, err
}

func (r *EntityRepository[T]) Delete(id int) error {
	var zero T
	query := fmt.Sprintf("DELETE FROM %s WHERE id = ?", zero.TableName())
	_, err := r.db.Exec(query, id)
	return err
}

func (r *EntityRepository[T]) FindAll() ([]T, error) {
	var zero T
	query := fmt.Sprintf("SELECT * FROM %s", zero.TableName())

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []T
	for rows.Next() {
		var item T
		if err := item.ScanRows(rows); err != nil {
			return nil, err
		}
		results = append(results, item)
	}
	return results, rows.Err()
}

func (r *EntityRepository[T]) Update(entity T) error {
	cols := entity.Columns()
	setClauses := make([]string, len(cols))
	for i, c := range cols {
		setClauses[i] = fmt.Sprintf("%s = ?", c)
	}

	query := fmt.Sprintf("UPDATE %s SET %s WHERE id = ?",
		entity.TableName(),
		strings.Join(setClauses, ", "))

	args := append(entity.Values(), entity.GetID())
	_, err := r.db.Exec(query, args...)
	return err
}
