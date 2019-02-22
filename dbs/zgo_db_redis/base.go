package zgo_db_redis

import (
	"github.com/mediocregopher/radix"
)

type RedisResource struct {
	RedisClient *radix.Pool
}

func NewRedisResource() *RedisResource {
	return &RedisResource{RedisClient: client}
}

func (r *RedisResource) GetRedisClient() *radix.Pool {
	return r.RedisClient
}
