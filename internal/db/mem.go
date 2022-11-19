package db

import (
	"context"
	"sync"

	"github.com/eduardohoraciosanto/users/internal/errors"
	"github.com/eduardohoraciosanto/users/internal/logger"
)

type db struct {
	data map[string]interface{}
	m    sync.Mutex
	log  logger.Logger
}

func NewMemDB(l logger.Logger) DB {
	return &db{
		data: make(map[string]interface{}),
		m:    sync.Mutex{},
		log:  l,
	}
}

func (d *db) Alive(ctx context.Context) bool {
	return true
}

func (d *db) Set(ctx context.Context, key string, data interface{}) error {
	d.m.Lock()
	defer d.m.Unlock()

	d.data[key] = data

	d.log.WithField("key", key).Info(ctx, "Data saved on DB")
	return nil
}

func (d *db) Get(ctx context.Context, key string) (interface{}, error) {
	d.m.Lock()
	defer d.m.Unlock()

	data, ok := d.data[key]
	if !ok {
		d.log.WithField("key", key).Error(ctx, "record not found on DB")
		return nil, errors.New(errors.DBErrorNotFoundCode, "record not found on DB")
	}
	d.log.WithField("key", key).Info(ctx, "Got data from DB")
	return data, nil
}

func (d *db) Delete(ctx context.Context, key string) error {
	d.m.Lock()
	defer d.m.Unlock()

	delete(d.data, key)

	d.log.WithField("key", key).Info(ctx, "Deleted data from DB")
	return nil
}
