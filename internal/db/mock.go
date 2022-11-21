package db

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type DBMock struct {
	mock.Mock
}

func (d *DBMock) Set(ctx context.Context, key string, data interface{}) error {
	args := d.Called(ctx, key, data)
	return args.Error(0)
}
func (d *DBMock) Get(ctx context.Context, key string, here interface{}) error {
	args := d.Called(ctx, key, here)
	return args.Error(0)
}
func (d *DBMock) Delete(ctx context.Context, key string) error {
	args := d.Called(ctx, key)
	return args.Error(0)
}
func (d *DBMock) Alive(ctx context.Context) bool {
	args := d.Called(ctx)
	return args.Bool(0)
}
