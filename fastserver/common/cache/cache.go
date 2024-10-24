package cache

import (
	"fastgin/modules/sys/model"
	"fmt"
	"github.com/patrickmn/go-cache"

	"time"
)

var Cache = cache.New(30*time.Minute, 3*time.Hour)     //存放user之外的数据，key是uuid
var UserCache = cache.New(30*time.Minute, 3*time.Hour) //存放uid->user对应关系

func SetUser(u *model.User) {
	UserCache.Set(u.GetUidString(), u, cache.DefaultExpiration)
}
func GetUser(uid uint64) *model.User {
	v, e := UserCache.Get(fmt.Sprintf("%d", uid))
	if !e {
		return nil
	}
	return v.(*model.User)
}
func GetString(key string) string {
	v, e := Cache.Get(key)
	if !e {
		return ""
	}
	return v.(string)
}
