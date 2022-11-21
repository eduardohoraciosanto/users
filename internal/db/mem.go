package db

import (
	"context"
	"encoding/json"
	"sync"

	"github.com/eduardohoraciosanto/users/internal/errors"
	"github.com/eduardohoraciosanto/users/internal/logger"
)

type db struct {
	data map[string]string
	m    sync.Mutex
	log  logger.Logger
}

func NewMemDB(l logger.Logger) DB {
	return &db{
		data: make(map[string]string),
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

	b, err := json.Marshal(data)
	if err != nil {
		d.log.WithField("key", key).WithError(err).Error(ctx, "unable to save record on DB")
		return errors.New(errors.InternalErrorCode, "unable to save record on DB")
	}

	d.data[key] = string(b)

	d.log.WithField("key", key).Info(ctx, "Data saved on DB")
	return nil
}

func (d *db) Get(ctx context.Context, key string, here interface{}) error {
	d.m.Lock()
	defer d.m.Unlock()

	data, ok := d.data[key]
	if !ok {
		d.log.WithField("key", key).Error(ctx, "record not found on DB")
		return errors.New(errors.DBErrorNotFoundCode, "record not found on DB")
	}
	d.log.WithField("key", key).Info(ctx, "Got data from DB")

	err := json.Unmarshal([]byte(data), here)
	if err != nil {
		d.log.WithField("key", key).WithError(err).Error(ctx, "unable to place data on interface")
		return errors.New(errors.InternalErrorCode, "unable to place data on interface")
	}

	return nil
}

func (d *db) Delete(ctx context.Context, key string) error {
	d.m.Lock()
	defer d.m.Unlock()

	delete(d.data, key)

	d.log.WithField("key", key).Info(ctx, "Deleted data from DB")
	return nil
}
