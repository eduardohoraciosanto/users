package db

import (
	"context"
)

type DB interface {
	Set(ctx context.Context, key string, data interface{}) error
	Get(ctx context.Context, key string, here interface{}) error
	Delete(ctx context.Context, key string) error
	Alive(ctx context.Context) bool
}
