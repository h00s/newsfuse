package services

import (
	"sync"
	"time"

	"github.com/dgraph-io/ristretto/v2"
	"github.com/go-raptor/raptor/v4"
)

// CacheService holds hot read results in memory. Every entry has a TTL, so a write that forgets
// to invalidate costs at most that long.
type CacheService struct {
	raptor.Service

	cache *ristretto.Cache[string, any]

	// generations counts invalidations per key, so a load that was running when its key was
	// invalidated knows its result is already stale and doesn't store it.
	mu          sync.Mutex
	generations map[string]uint64
}

func (s *CacheService) Setup() error {
	s.generations = map[string]uint64{}
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
	s.mu.Lock()
	s.generations[key]++
	s.mu.Unlock()
	s.cache.Del(key)
}

// cached returns the value under key, calling load and storing its result on a miss. A failed
// load is not cached, and neither is one that raced an Invalidate of its key. Callers must treat
// the value as read-only: every reader shares it.
func cached[T any](c *CacheService, key string, ttl time.Duration, load func() (T, error)) (T, error) {
	if v, ok := c.cache.Get(key); ok {
		if typed, ok := v.(T); ok {
			return typed, nil
		}
	}

	c.mu.Lock()
	generation := c.generations[key]
	c.mu.Unlock()

	v, err := load()
	if err != nil {
		return v, err
	}

	// Held through the set, so an Invalidate either happened before the check (and the value is
	// dropped) or waits and deletes the value after it is stored.
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.generations[key] == generation {
		c.cache.SetWithTTL(key, v, 1, ttl)
		c.cache.Wait() // sets are buffered; make this one visible to the next Get
	}
	return v, nil
}
