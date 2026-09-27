package services

import (
	"testing"
	"time"
)

func newTestCache(t *testing.T) *CacheService {
	t.Helper()
	c := &CacheService{}
	if err := c.Setup(); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = c.Cleanup() })
	return c
}

func TestCachedServesTheStoredValue(t *testing.T) {
	c := newTestCache(t)
	loads := 0
	load := func() (int, error) { loads++; return loads, nil }

	_, _ = cached(c, "k", time.Minute, load)
	got, _ := cached(c, "k", time.Minute, load)

	if got != 1 || loads != 1 {
		t.Errorf("second read = %d after %d loads, want the stored 1 after 1 load", got, loads)
	}
}

func TestCachedReloadsAfterInvalidate(t *testing.T) {
	c := newTestCache(t)
	loads := 0
	load := func() (int, error) { loads++; return loads, nil }

	_, _ = cached(c, "k", time.Minute, load)
	c.Invalidate("k")
	got, _ := cached(c, "k", time.Minute, load)

	if got != 2 {
		t.Errorf("read after Invalidate = %d, want a fresh load (2)", got)
	}
}

// An ingest that invalidates while a read is loading must win: the read's result predates the
// new rows, so storing it would serve a stale page until the TTL ran out.
func TestCachedDoesNotStoreALoadThatRacedAnInvalidate(t *testing.T) {
	c := newTestCache(t)
	loads := 0
	racing := func() (int, error) {
		loads++
		c.Invalidate("k") // the ingest lands while this load is running
		return loads, nil
	}

	first, _ := cached(c, "k", time.Minute, racing)
	got, _ := cached(c, "k", time.Minute, func() (int, error) { loads++; return loads, nil })

	if first != 1 {
		t.Errorf("racing read = %d, want its own load (1)", first)
	}
	if got != 2 {
		t.Errorf("next read = %d, want a fresh load (2), not the stale value stored after the invalidate", got)
	}
}
