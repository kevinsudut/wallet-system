package lrucache

import (
	"time"

	"github.com/karlseguin/ccache/v3"
)

func (c lruCache) Get(key string) *ccache.Item[interface{}] {
	return c.cache.Get(key)
}

func (c lruCache) Set(key string, value interface{}, duration time.Duration) {
	c.cache.Set(key, value, duration)
}

func (c lruCache) Delete(key string) bool {
	return c.cache.Delete(key)
}

func (c lruCache) Fetch(key string, duration time.Duration, fetch func() (interface{}, error)) (*ccache.Item[interface{}], error) {
	return c.cache.Fetch(key, duration, fetch)
}
