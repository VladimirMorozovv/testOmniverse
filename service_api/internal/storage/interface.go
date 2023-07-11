package storage

import (
	"context"
	"service_api/internal/model"
)

type Product struct {
	Id    string
	Price int
}

//go:generate moq -out storage_sql_mock.go -rm . IStorageSQL
type IStorageSQL interface {
	GetProduct(ctx context.Context, limit int, offset int) ([]Product, error)
}

type ICache interface {
	Get(ctx context.Context, key model.Params, store IStorageSQL) ([]Product, error)
	Shutdown(ctx context.Context) error
}
