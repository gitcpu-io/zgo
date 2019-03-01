package zgopika

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

//Pika 对外
type Pikaer interface {
	NewPika(label ...string) (*zgopika, error)
	//Post
	Set(ctx context.Context, key string, value string, time int) (interface{}, error)
	Expire(ctx context.Context, key string, time int) (interface{}, error)
	Hset(ctx context.Context, key string, name string, value string) (interface{}, error)
	Lpush(ctx context.Context, key string, value string) (interface{}, error)
	Rpush(ctx context.Context, key string, value string) (interface{}, error)
	Sadd(ctx context.Context, key string, value string) (interface{}, error)
	Srem(ctx context.Context, key string, value string) (interface{}, error)
	//Get
	Exists(ctx context.Context, key string) (interface{}, error)
	Get(ctx context.Context, key string) (interface{}, error)
	Keys(ctx context.Context, pattern string) (interface{}, error)
	Hget(ctx context.Context, key string, name string) (interface{}, error)
	Ttl(ctx context.Context, key string) (interface{}, error)
	Type(ctx context.Context, key string) (interface{}, error)
	Hlen(ctx context.Context, key string) (interface{}, error)
	Hdel(ctx context.Context, key string, name string) (interface{}, error)
	Hgetall(ctx context.Context, key string) (interface{}, error)
	Del(ctx context.Context, key string) (interface{}, error)
	Llen(ctx context.Context, key string) (interface{}, error)
	Lrange(ctx context.Context, key string, start int, stop int) (interface{}, error)
	Lpop(ctx context.Context, key string) (interface{}, error)
	Rpop(ctx context.Context, key string) (interface{}, error)
	Scard(ctx context.Context, key string) (interface{}, error)
	Smembers(ctx context.Context, key string) (interface{}, error)
	Sismember(ctx context.Context, key string, value string) (interface{}, error)
}

func Pika(l string) Pikaer {
	return &zgopika{
		res: NewPikaResource(l),
	}
}

//zgopika实现了Redis的接口
type zgopika struct {
	res PikaResourcer //使用resource另外的一个接口
}

//InitPika 初始化连接pika
func InitPika(hsm map[string][]*config.ConnDetail) {
	muLabel.Lock()
	defer muLabel.Unlock()

	currentLabels = hsm
	InitPikaResource(hsm)
}

func (n *zgopika) NewPika(label ...string) (*zgopika, error) {
	return GetPika(label...)
}

//GetRedis zgo内部获取一个连接pika
func GetPika(label ...string) (*zgopika, error) {
	l, err := comm.GetCurrentLabel(label, muLabel, currentLabels)
	if err != nil {
		return nil, err
	}
	return &zgopika{
		res: NewPikaResource(l), //interface
	}, nil
}

func (r *zgopika) Set(ctx context.Context, key string, value string, time int) (interface{}, error) {
	return r.res.Set(ctx, key, value, time)
}

func (r *zgopika) Expire(ctx context.Context, key string, time int) (interface{}, error) {
	return r.res.Expire(ctx, key, time)
}

func (r *zgopika) Hset(ctx context.Context, key string, name string, value string) (interface{}, error) {
	return r.res.Hset(ctx, key, name, value)
}

func (r *zgopika) Lpush(ctx context.Context, key string, value string) (interface{}, error) {
	return r.res.Lpush(ctx, key, value)
}

func (r *zgopika) Rpush(ctx context.Context, key string, value string) (interface{}, error) {
	return r.res.Lpush(ctx, key, value)
}

func (r *zgopika) Sadd(ctx context.Context, key string, value string) (interface{}, error) {
	return r.res.Sadd(ctx, key, value)
}

func (r *zgopika) Srem(ctx context.Context, key string, value string) (interface{}, error) {
	return r.res.Srem(ctx, key, value)
}

func (r *zgopika) Exists(ctx context.Context, key string) (interface{}, error) {
	return r.res.Exists(ctx, key)
}

func (r *zgopika) Get(ctx context.Context, key string) (interface{}, error) {
	return r.res.Get(ctx, key)
}

func (r *zgopika) Keys(ctx context.Context, key string) (interface{}, error) {
	return r.res.Keys(ctx, key)
}

func (r *zgopika) Hget(ctx context.Context, key string, name string) (interface{}, error) {
	return r.res.Hget(ctx, key, name)
}

func (r *zgopika) Ttl(ctx context.Context, key string) (interface{}, error) {
	return r.res.Ttl(ctx, key)
}

func (r *zgopika) Type(ctx context.Context, key string) (interface{}, error) {
	return r.res.Type(ctx, key)
}

func (r *zgopika) Hlen(ctx context.Context, key string) (interface{}, error) {
	return r.res.Hlen(ctx, key)
}

func (r *zgopika) Hdel(ctx context.Context, key string, name string) (interface{}, error) {
	return r.res.Hdel(ctx, key, name)
}

func (r *zgopika) Hgetall(ctx context.Context, key string) (interface{}, error) {
	return r.res.Hgetall(ctx, key)
}

func (r *zgopika) Del(ctx context.Context, key string) (interface{}, error) {
	return r.res.Del(ctx, key)
}

func (r *zgopika) Llen(ctx context.Context, key string) (interface{}, error) {
	return r.res.Llen(ctx, key)
}

func (r *zgopika) Lrange(ctx context.Context, key string, start int, stop int) (interface{}, error) {
	return r.res.Lrange(ctx, key, start, stop)
}

func (r *zgopika) Lpop(ctx context.Context, key string) (interface{}, error) {
	return r.res.Lpop(ctx, key)
}

func (r *zgopika) Rpop(ctx context.Context, key string) (interface{}, error) {
	return r.res.Lpop(ctx, key)
}

func (r *zgopika) Scard(ctx context.Context, key string) (interface{}, error) {
	return r.res.Scard(ctx, key)
}

func (r *zgopika) Smembers(ctx context.Context, key string) (interface{}, error) {
	return r.res.Smembers(ctx, key)
}

func (r *zgopika) Sismember(ctx context.Context, key string, value string) (interface{}, error) {
	return r.res.Sismember(ctx, key, value)
}
