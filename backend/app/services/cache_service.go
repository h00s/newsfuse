package services

import (
	"github.com/dgraph-io/ristretto/v2"
	"github.com/go-raptor/raptor/v4"
)

type CacheService struct {
	raptor.Service

	cache *ristretto.Cache[string, []byte]
}

func NewCacheService() *CacheService {
	cs := &CacheService{}

	cs.OnInit(func() error {
		var err error

		cs.cache, err = ristretto.NewCache(&ristretto.Config[string, []byte]{
			NumCounters: 500,
			MaxCost:     100000,
			BufferItems: 64,
		})

		return err
	})

	return cs
}

func (cs *CacheService) Get(key string) ([]byte, bool) {
	value, found := cs.cache.Get(key)
	if !found {
		return nil, false
	}
	return value, true
}

func (cs *CacheService) Set(key string, value []byte) {
	cs.cache.Set(key, value, 1)
}
