package zgoredis

import (
	"context"
	"git.zhugefang.com/gocore/zgo.git/comm"
	"sync"
)

var (
	currentLabels = make(map[string][]string)
	muLabel       sync.RWMutex
)

//Redis 对外
type Rediser interface {
	NewRedis(label ...string) (*zgoredis, error)
	Do(ctx context.Context, rcv interface{}, cmd string, args ...string) (interface{}, error)
}

func Redis(l string) Rediser {
	return &zgoredis{
		res: NewRedisResource(l),
	}
}

//zgoredis实现了Redis的接口
type zgoredis struct {
	res RedisResourcer //使用resource另外的一个接口
}

//InitRedis 初始化连接redis
func InitRedis(hsm map[string][]string) {
	muLabel.Lock()
	defer muLabel.Unlock()

	currentLabels = hsm
	InitRedisResource(hsm)
}

func (n *zgoredis) NewRedis(label ...string) (*zgoredis, error) {
	return GetRedis(label...)
}

//GetRedis zgo内部获取一个连接redis
func GetRedis(label ...string) (*zgoredis, error) {
	l, err := comm.GetCurrentLabel(label, muLabel, currentLabels)
	if err != nil {
		return nil, err
	}
	return &zgoredis{
		res: NewRedisResource(l), //interface
	}, nil
}

func (m *zgoredis) Do(ctx context.Context, rcv interface{}, cmd string, args ...string) (interface{}, error) {
	return m.res.Do(ctx, rcv, cmd, args...)
}
