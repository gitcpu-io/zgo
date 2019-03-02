package zgoredis

import (
	"context"
	"git.zhugefang.com/gocore/zgo.git/comm"
	"git.zhugefang.com/gocore/zgo.git/config"
	"sync"
)

var (
	currentLabels = make(map[string][]*config.ConnDetail)
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
func InitRedis(hsm map[string][]*config.ConnDetail) chan *zgoredis {
	muLabel.Lock()
	defer muLabel.Unlock()

	currentLabels = hsm
	InitRedisResource(hsm)

	//自动为变量初始化对象
	initLabel := ""
	for k, _ := range hsm {
		if k != "" {
			initLabel = k
			break
		}
	}
	out := make(chan *zgoredis)
	go func() {
		in, err := GetRedis(initLabel)
		if err != nil {
			out <- nil
		}
		out <- in
		close(out)
	}()
	return out
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

func (r *zgoredis) Do(ctx context.Context, rcv interface{}, cmd string, args ...string) (interface{}, error) {
	return r.res.Do(ctx, rcv, cmd, args...)
}
