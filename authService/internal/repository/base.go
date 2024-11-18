package repository

import (
	"authService/config"
	"authService/db"
	"context"

	"github.com/google/uuid"
)

type Operations[T any] interface {
    Insert(ctx context.Context, data T) (T, error)
    Update(ctx context.Context, data T) (T, error)
    Delete(ctx context.Context, id uuid.UUID) error
    Query(ctx context.Context, id uuid.UUID) (T, error)
}

func init(){
	db.InitDB(config.GetDBConfig())
}