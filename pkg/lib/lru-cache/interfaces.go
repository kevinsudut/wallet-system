package lrucache

import (
	"time"

	"github.com/karlseguin/ccache/v3"
)

type LRUCacheItf interface {
	Get(key string) *ccache.Item[interface{}]
	Set(key string, value interface{}, duration time.Duration)
	Delete(key string) bool
	Fetch(key string, duration time.Duration, fetch func() (interface{}, error)) (*ccache.Item[interface{}], error)
}
