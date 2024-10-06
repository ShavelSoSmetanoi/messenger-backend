package postgres

import (
	"context"
	"database/sql"
)

type GenericRepository[T any] struct {
	db *sql.DB
}

func (r *GenericRepository[T]) Create(ctx context.Context, entity *T) (int64, error) {
	// Логика вставки данных в базу
	return 0, nil
}

func (r *GenericRepository[T]) Get(ctx context.Context, id int64) (*T, error) {
	// Логика получения данных из базы
	return nil, nil
}
