package app

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/go-redis/redis"
)

// DefaultExpiration is the default expiration of items stored in redis
const DefaultExpiration = 10 * time.Minute

// CacheItem is an item that can be stored in the cache
type CacheItem struct {
	Content []byte     `json:"content"`
	ModTime *time.Time `json:"modtime"`
}

// Cache is a wrapper around redis for easy consumption
type Cache struct {
	client *redis.Client
}

// NewCache generates a new redis client wrapper
func NewCache(client *redis.Client) *Cache {
	return &Cache{client}
}

// Get an item from the cache
func (c *Cache) Get(key string) (*CacheItem, bool) {
	value := &CacheItem{}

	err := c.client.Get(key).Scan(value)
	if err != nil {
		return nil, false
	}

	return value, true
}

// Set an item in the cache
func (c *Cache) Set(key string, value *CacheItem) error {
	if value == nil {
		return errors.New("no cache item given")
	}

	bytes, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return c.client.Set(key, bytes, DefaultExpiration).Err()
}
