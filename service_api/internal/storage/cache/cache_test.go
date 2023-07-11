package cache

import (
	"context"
	"service_api/internal/model"
	"service_api/internal/storage"
	"testing"
)

func TestGet(t *testing.T) {
	cache, err := NewCacheInMemory(1000, 10)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	postgres := &storage.IStorageSQLMock{
		GetProductFunc: func(ctx context.Context, limit int, offset int) ([]storage.Product, error) {
			result := make([]storage.Product, 0)
			result = append(result, storage.Product{
				Id:    "pr1",
				Price: 12,
			})
			return result, nil
		},
	}
	postgresFail := &storage.IStorageSQLMock{
		GetProductFunc: func(ctx context.Context, limit int, offset int) ([]storage.Product, error) {

			return nil, nil
		},
	}
	_, err = cache.Get(context.Background(), model.Params{Limit: 10, Offset: 10}, postgres)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	res2, err := cache.Get(context.Background(), model.Params{Limit: 10, Offset: 10}, postgresFail)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	if len(res2) == 0 {
		t.Error("invalid job cache")
		t.FailNow()
	}

}
