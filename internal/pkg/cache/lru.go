package cache

import (
	"time"

	lru "github.com/hashicorp/golang-lru/v2"
)

type Cache interface {
	Get(key string) (any, bool)
	Set(key string, value any)
	Delete(key string)
	Purge()
}

type lruCache struct {
	cache *lru.Cache[string, any]
}

// NewLRUCache creates a new thread-safe LRU cache with the given size.
func NewLRUCache(size int) (Cache, error) {
	c, err := lru.New[string, any](size)
	if err != nil {
		return nil, err
	}
	return &lruCache{cache: c}, nil
}

func (c *lruCache) Get(key string) (any, bool) {
	return c.cache.Get(key)
}

func (c *lruCache) Set(key string, value any) {
	c.cache.Add(key, value)
}

func (c *lruCache) Delete(key string) {
	c.cache.Remove(key)
}

func (c *lruCache) Purge() {
	c.cache.Purge()
}

// Item with TTL (Optional wrapper if needed later, but for now simple LRU is fine)
type Item struct {
	Value      any
	Expiration int64
}

func (c *lruCache) SetWithTTL(key string, value any, ttl time.Duration) {
	c.cache.Add(key, Item{
		Value:      value,
		Expiration: time.Now().Add(ttl).UnixNano(),
	})
}
