package queries

import (
	"shinkyuShotokan/packages/cache"
)

var Cache *cache.MemoryCache

func InitCache(c *cache.MemoryCache) {
	Cache = c
}
