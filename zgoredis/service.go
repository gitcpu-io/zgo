package zgoredis

import (
	"context"
	"git.zhugefang.com/gocore/zgo/comm"
	"git.zhugefang.com/gocore/zgo/config"
	"sync"
)

var (
	currentLabels = make(map[string][]*config.ConnDetail)
	muLabel       sync.RWMutex
)

//Redis 对外
type Rediser interface {
	New(label ...string) (*zgoredis, error)
	//Post
	//hmset setnx setex
	Set(ctx context.Context, key string, value interface{}) (string, error)
	//SETNX if Not eXists 1 如果key被设置了; 0 如果key没有被设置
	Setnx(ctx context.Context, key string, value interface{}) (int, error)
	//对key设置ttl为秒的过期; OK表示成功
	Setex(ctx context.Context, key string, ttl int, value interface{}) (string, error)
	Expire(ctx context.Context, key string, time int) (int, error)
	Hset(ctx context.Context, key string, name string, value interface{}) (int, error)
	Hmset(ctx context.Context, key string, values interface{}) (string, error)
	Lpush(ctx context.Context, key string, value interface{}) (int, error)
	Rpush(ctx context.Context, key string, value interface{}) (int, error)
	Sadd(ctx context.Context, key string, value interface{}) (int, error)
	Srem(ctx context.Context, key string, value interface{}) (int, error)
	//Get
	Exists(ctx context.Context, key string) (int, error)
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
	Ltrim(ctx context.Context, key string, start int, stop int) (interface{}, error)
	Lpop(ctx context.Context, key string) (interface{}, error)
	Rpop(ctx context.Context, key string) (interface{}, error)
	Scard(ctx context.Context, key string) (interface{}, error)
	Smembers(ctx context.Context, key string) (interface{}, error)
	Sismember(ctx context.Context, key string, value interface{}) (int, error)
	Zrank(ctx context.Context, key string, member interface{}) (interface{}, error)
	Zscore(ctx context.Context, key string, member interface{}) (string, error)
	Zrange(ctx context.Context, key string, start int, stop int, withscores bool) (interface{}, error)
	Zrevrange(ctx context.Context, key string, start int, stop int, withscores bool) (interface{}, error)
	ZINCRBY(ctx context.Context, key string, increment int, member interface{}) (string, error)
	Zadd(ctx context.Context, key string, score interface{}, member interface{}) (int, error)
	Zrem(ctx context.Context, key string, member ...interface{}) (int, error)
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
func InitRedis(hsmIn map[string][]*config.ConnDetail, label ...string) chan *zgoredis {
	muLabel.Lock()
	defer muLabel.Unlock()

	var hsm map[string][]*config.ConnDetail

	if len(label) > 0 && len(currentLabels) > 0 { //此时是destory操作,传入的hsm是nil
		//fmt.Println("--destory--前",currentLabels)
		for _, v := range label {
			delete(currentLabels, v)
		}
		hsm = currentLabels
		//fmt.Println("--destory--后",currentLabels)

	} else { //这是第一次创建操作或etcd中变更时init again操作
		hsm = hsmIn
		//currentLabels = hsm	//this operation is error
		for k, v := range hsm { //so big bug can't set hsm to currentLabels，must be for, may be have old label
			currentLabels[k] = v
		}
	}

	if len(hsm) == 0 {
		return nil
	}

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

func (n *zgoredis) New(label ...string) (*zgoredis, error) {
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

func (r *zgoredis) Set(ctx context.Context, key string, value interface{}) (string, error) {
	return r.res.Set(ctx, key, value)
}

func (r *zgoredis) Setnx(ctx context.Context, key string, value interface{}) (int, error) {
	return r.res.Setnx(ctx, key, value)
}

func (r *zgoredis) Setex(ctx context.Context, key string, ttl int, value interface{}) (string, error) {
	return r.res.Setex(ctx, key, ttl, value)
}

func (r *zgoredis) Expire(ctx context.Context, key string, time int) (int, error) {
	return r.res.Expire(ctx, key, time)
}

func (r *zgoredis) Hset(ctx context.Context, key string, name string, value interface{}) (int, error) {
	return r.res.Hset(ctx, key, name, value)
}

func (r *zgoredis) Hmset(ctx context.Context, key string, values interface{}) (string, error) {
	return r.res.Hmset(ctx, key, values)
}

func (r *zgoredis) Lpush(ctx context.Context, key string, value interface{}) (int, error) {
	return r.res.Lpush(ctx, key, value)
}

func (r *zgoredis) Rpush(ctx context.Context, key string, value interface{}) (int, error) {
	return r.res.Rpush(ctx, key, value)
}

func (r *zgoredis) Sadd(ctx context.Context, key string, value interface{}) (int, error) {
	return r.res.Sadd(ctx, key, value)
}

func (r *zgoredis) Srem(ctx context.Context, key string, value interface{}) (int, error) {
	return r.res.Srem(ctx, key, value)
}

func (r *zgoredis) Exists(ctx context.Context, key string) (int, error) {
	return r.res.Exists(ctx, key)
}

func (r *zgoredis) Get(ctx context.Context, key string) (interface{}, error) {
	return r.res.Get(ctx, key)
}

func (r *zgoredis) Keys(ctx context.Context, key string) (interface{}, error) {
	return r.res.Keys(ctx, key)
}

func (r *zgoredis) Hget(ctx context.Context, key string, name string) (interface{}, error) {
	return r.res.Hget(ctx, key, name)
}

func (r *zgoredis) Ttl(ctx context.Context, key string) (interface{}, error) {
	return r.res.Ttl(ctx, key)
}

func (r *zgoredis) Type(ctx context.Context, key string) (interface{}, error) {
	return r.res.Type(ctx, key)
}

func (r *zgoredis) Hlen(ctx context.Context, key string) (interface{}, error) {
	return r.res.Hlen(ctx, key)
}

func (r *zgoredis) Hdel(ctx context.Context, key string, name string) (interface{}, error) {
	return r.res.Hdel(ctx, key, name)
}

func (r *zgoredis) Hgetall(ctx context.Context, key string) (interface{}, error) {
	return r.res.Hgetall(ctx, key)
}

func (r *zgoredis) Del(ctx context.Context, key string) (interface{}, error) {
	return r.res.Del(ctx, key)
}

func (r *zgoredis) Llen(ctx context.Context, key string) (interface{}, error) {
	return r.res.Llen(ctx, key)
}

func (r *zgoredis) Lrange(ctx context.Context, key string, start int, stop int) (interface{}, error) {
	return r.res.Lrange(ctx, key, start, stop)
}

func (r *zgoredis) Ltrim(ctx context.Context, key string, start int, stop int) (interface{}, error) {
	return r.res.Ltrim(ctx, key, start, stop)
}

func (r *zgoredis) Lpop(ctx context.Context, key string) (interface{}, error) {
	return r.res.Lpop(ctx, key)
}

func (r *zgoredis) Rpop(ctx context.Context, key string) (interface{}, error) {
	return r.res.Rpop(ctx, key)
}

func (r *zgoredis) Scard(ctx context.Context, key string) (interface{}, error) {
	return r.res.Scard(ctx, key)
}

func (r *zgoredis) Smembers(ctx context.Context, key string) (interface{}, error) {
	return r.res.Smembers(ctx, key)
}

func (r *zgoredis) Sismember(ctx context.Context, key string, value interface{}) (int, error) {
	return r.res.Sismember(ctx, key, value)
}

func (r *zgoredis) Zrank(ctx context.Context, key string, member interface{}) (interface{}, error) {
	return r.res.Zrank(ctx, key, member)
}

func (r *zgoredis) Zscore(ctx context.Context, key string, member interface{}) (string, error) {
	return r.res.Zscore(ctx, key, member)
}

func (r *zgoredis) Zrange(ctx context.Context, key string, start int, stop int, withscores bool) (interface{}, error) {
	return r.res.Zrange(ctx, key, start, stop, withscores)
}

func (r *zgoredis) Zrevrange(ctx context.Context, key string, start int, stop int, withscores bool) (interface{}, error) {
	return r.res.Zrevrange(ctx, key, start, stop, withscores)
}

func (r *zgoredis) ZINCRBY(ctx context.Context, key string, increment int, member interface{}) (string, error) {
	return r.res.ZINCRBY(ctx, key, increment, member)
}

func (r *zgoredis) Zadd(ctx context.Context, key string, score interface{}, member interface{}) (int, error) {
	return r.res.Zadd(ctx, key, score, member)
}

func (r *zgoredis) Zrem(ctx context.Context, key string, member ...interface{}) (int, error) {
	return r.res.Zrem(ctx, key, member)
}
