package zgo

import "zgo/logic/zgo_redis"

//todo  same as mongo or utils

var Redis *zgo_redis.Redis

func init() {
	Redis = zgo_redis.NewRedis()
}
