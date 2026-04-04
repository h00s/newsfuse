package services

import (
	"github.com/dgraph-io/ristretto/v2"
	"github.com/go-raptor/raptor/v4"
)

type CacheService struct {
	raptor.Service

	cache *ristretto.Cache[string, []byte]
}

func (s *CacheService) Setup() error {
	var err error
	s.cache, err = ristretto.NewCache(&ristretto.Config[string, []byte]{
		NumCounters: 500,
		MaxCost:     100000,
		BufferItems: 64,
	})

	return err
}

func (s *CacheService) Get(key string) ([]byte, bool) {
	value, found := s.cache.Get(key)
	if !found {
		return nil, false
	}
	return value, true
}

func (s *CacheService) Set(key string, value []byte) {
	s.cache.Set(key, value, 1)
}
