package lrucache

import "github.com/karlseguin/ccache/v3"

type lruCache struct {
	cache *ccache.Cache[string]
}

func Init() LRUCacheItf {
	return &lruCache{
		cache: ccache.New(ccache.Configure[string]()),
	}
}
