package zgoredis

import (
	"context"
	"git.zhugefang.com/gocore/zgo.git/config"
	"github.com/mediocregopher/radix"
	"sync"
)

//NsqResourcer 给service使用
type RedisResourcer interface {
	GetConnChan(label string) chan *radix.Pool
	Do(ctx context.Context, rcv interface{}, cmd string, args ...string) (interface{}, error)
}

type redisResource struct {
	label    string
	mu       sync.RWMutex
	connpool ConnPooler
}

func InitRedisResource(hsm map[string][]config.ConnDetail) {
	InitConnPool(hsm)
}

func NewRedisResource(label string) RedisResourcer {
	return &redisResource{
		label:    label,
		connpool: NewConnPool(label)}
}

//GetConnChan 返回存放连接的chan
func (r *redisResource) GetConnChan(label string) chan *radix.Pool {
	return r.connpool.GetConnChan(label)
}

func (r *redisResource) Do(ctx context.Context, rcv interface{}, cmd string, args ...string) (interface{}, error) {
	s := <-r.connpool.GetConnChan(r.label)
	return nil, s.Do(radix.Cmd(rcv, cmd, args...))
}
