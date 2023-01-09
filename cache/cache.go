package cache

import (
	"errors"
	"sync"
	"time"
)

var (
	ErrInvalidKey = errors.New("key does not exist")
)

type Cache interface {
	Get(key string) (string, error)
	Set(key, value string) (string, error)
	Del(key string) (string, error)
	Expire(key string, ttl time.Duration) (string, error)
	deleteAfter(key string, ttl time.Duration)
}

type CacheStorage struct {
	mu    sync.RWMutex
	cache map[string]string
}

func New() Cache {
	return &CacheStorage{
		mu:    sync.RWMutex{},
		cache: make(map[string]string),
	}
}

func (c *CacheStorage) Get(key string) (string, error) {
	defer c.mu.RLock()
	c.mu.RLock()

	val, ok := c.cache[key]
	if !ok {
		return "", ErrInvalidKey
	}

	return val, nil
}

func (c *CacheStorage) deleteAfter(key string, ttl time.Duration) {
	<-time.After(ttl)
	delete(c.cache, key)
}

func (c *CacheStorage) Expire(key string, ttl time.Duration) (string, error) {
	_, ok := c.cache[key]
	if !ok {
		return "", ErrInvalidKey
	}
	go c.deleteAfter(key, ttl)

	return "OK", nil
}

func (c *CacheStorage) Set(key, value string) (string, error) {
	defer c.mu.RLock()
	c.mu.RLock()

	c.cache[key] = value

	return "OK", nil
}

func (c *CacheStorage) Del(key string) (string, error) {
	defer c.mu.RLock()
	c.mu.RLock()

	delete(c.cache, key)

	return "OK", nil
}
