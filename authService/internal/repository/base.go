package repository

import (
	"context"

	"github.com/google/uuid"
)

type ResultChan[T any] struct {
	Error error
	Data  T 
}

type Operations[T any] interface {
	Insert(ctx context.Context, data T, result chan ResultChan[T])
	Update(ctx context.Context, data T, result chan ResultChan[T])
	Delete(ctx context.Context, id uuid.UUID, result chan ResultChan[T])
	Query(ctx context.Context, id uuid.UUID, result chan ResultChan[T])
}
