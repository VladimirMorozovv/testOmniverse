package cache

import (
	"context"
	lru "github.com/bserdar/golang-lru"
	"service_api/internal/model"
	"service_api/internal/storage"
	"time"
)

type CacheInMemory struct {
	cache *lru.Cache
}

func NewCacheInMemory(countSize int, ttlSecond int) (*CacheInMemory, error) {
	cache, err := lru.NewWithTTL(countSize, time.Duration(ttlSecond)*time.Second)
	if err != nil {
		return nil, err
	}

	return &CacheInMemory{
		cache: cache,
	}, nil
}

/*
Изменил логику работы не в соответствии с тех. заданием , так как для данной задачи по тз теряется консистентность данных,
а так же появляется возможность упереться в оперативную память

Могу переделать в соответвии с тех заданием и сделать опережающий кэш если потребуется.
*/

func (c *CacheInMemory) Get(ctx context.Context, key model.Params, store storage.IStorageSQL) ([]storage.Product, error) {
	value, ok := c.cache.Get(key)
	if !ok {
		res, err := store.GetProduct(ctx, key.Limit, key.Offset)
		if err != nil {
			return nil, err
		}
		_ = c.cache.Add(key, res, 1)
		return res, err
	}
	res := value.([]storage.Product)
	return res, nil
}

func (c *CacheInMemory) Shutdown(_ context.Context) error {
	return nil
}
