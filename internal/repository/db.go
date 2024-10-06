package repository

import (
	"context"
)

type Repository[T any] interface {
	Create(ctx context.Context, entity *T) (int64, error)
	Get(ctx context.Context, id int64) (*T, error)
}
