package lrucache

import (
	"time"

	"github.com/karlseguin/ccache/v3"
)

type LRUCacheItf interface {
	Get(key string) *ccache.Item[string]
	Set(key string, value string, duration time.Duration)
	Delete(key string) bool
	Fetch(key string, duration time.Duration, fetch func() (string, error)) (*ccache.Item[string], error)
}
