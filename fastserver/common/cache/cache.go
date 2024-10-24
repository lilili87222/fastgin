package cache

import (
	"github.com/patrickmn/go-cache"
	"time"
)

// var Cache, _ = bigcache.New(context.Background(), bigcache.DefaultConfig(10*time.Minute))
var Cache = cache.New(30*time.Minute, 3*time.Hour)

func GetString(key string) string {
	v, e := Cache.Get(key)
	if !e {
		return ""
	}
	return v.(string)
}
