package services

import (
	"time"

	"github.com/dgraph-io/ristretto/v2"
	"github.com/go-raptor/raptor/v4"
)

// CacheService holds hot read results in memory. Every entry has a TTL, so a write that forgets
// to invalidate costs at most that long.
type CacheService struct {
	raptor.Service

	cache *ristretto.Cache[string, any]
}

func (s *CacheService) Setup() error {
	var err error
	s.cache, err = ristretto.NewCache(&ristretto.Config[string, any]{
		NumCounters:        1_000,
		MaxCost:            100, // entries, each set with cost 1
		BufferItems:        64,
		IgnoreInternalCost: true,
	})
	return err
}

func (s *CacheService) Cleanup() error {
	s.cache.Close()
	return nil
}

// Invalidate drops a key, so the next read loads it again.
func (s *CacheService) Invalidate(key string) {
	s.cache.Del(key)
}

// cached returns the value under key, calling load and storing its result on a miss. A failed
// load is not cached. Callers must treat the value as read-only: every reader shares it.
func cached[T any](c *CacheService, key string, ttl time.Duration, load func() (T, error)) (T, error) {
	if v, ok := c.cache.Get(key); ok {
		if typed, ok := v.(T); ok {
			return typed, nil
		}
	}
	v, err := load()
	if err != nil {
		return v, err
	}
	c.cache.SetWithTTL(key, v, 1, ttl)
	c.cache.Wait() // sets are buffered; make this one visible to the next Get
	return v, nil
}
