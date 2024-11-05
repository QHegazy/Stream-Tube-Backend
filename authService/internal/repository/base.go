package repository

import (
	"context"

	"github.com/google/uuid"
)

type ResultChan[T any] struct {
	Error error
	Data  T
}

// Operations defines the CRUD operations that can be performed on a model.
type Operations[T any] interface {
	Insert(ctx context.Context, data T) <-chan ResultChan[T]
	Update(ctx context.Context, data T) <-chan ResultChan[T]
	Delete(ctx context.Context, id uuid.UUID) <-chan ResultChan[T]
	Query(ctx context.Context, id uuid.UUID) <-chan ResultChan[T]
}
