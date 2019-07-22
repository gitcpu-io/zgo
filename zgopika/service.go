package zgopika

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

//Pika 对外
type Pikaer interface {
	New(label ...string) (*zgopika, error)
	//Post
	Set(ctx context.Context, key string, value interface{}) (string, error)
	//设置分布式锁
	SetMutex(ctx context.Context, key string, ttl int, value interface{}) (string, error)
	//SETNX if Not eXists 1 如果key被设置了; 0 如果key没有被设置
	Setnx(ctx context.Context, key string, value interface{}) (int, error)
	//对key设置ttl为秒的过期; OK表示成功
	Setex(ctx context.Context, key string, ttl int, value interface{}) (string, error)
	Expire(ctx context.Context, key string, time int) (int, error)
	Incrby(ctx context.Context, key string, val int) (interface{}, error)
	Hset(ctx context.Context, key string, name string, value interface{}) (int, error)
	Hmset(ctx context.Context, key string, values interface{}) (string, error)

	Lpush(ctx context.Context, key string, value interface{}) (int, error)
	Rpush(ctx context.Context, key string, value interface{}) (int, error)
	Sadd(ctx context.Context, key string, value interface{}) (int, error)
	Srem(ctx context.Context, key string, value interface{}) (int, error)

	//Get
	Exists(ctx context.Context, key string) (interface{}, error)
	Get(ctx context.Context, key string) (interface{}, error)
	Keys(ctx context.Context, pattern string) (interface{}, error)
	Hget(ctx context.Context, key string, name string) (interface{}, error)
	Ttl(ctx context.Context, key string) (interface{}, error)
	Type(ctx context.Context, key string) (interface{}, error)
	Hlen(ctx context.Context, key string) (int, error)
	Hdel(ctx context.Context, key string, name interface{}) (int, error)
	Hgetall(ctx context.Context, key string) (interface{}, error)
	Hincrby(ctx context.Context, key string, field string, inc int64) (int64, error)

	Del(ctx context.Context, key string) (interface{}, error)
	Llen(ctx context.Context, key string) (int, error)
	Lrange(ctx context.Context, key string, start int, stop int) (interface{}, error)
	Ltrim(ctx context.Context, key string, start int, stop int) (interface{}, error)
	Lpop(ctx context.Context, key string) (interface{}, error)
	Lrem(ctx context.Context, key string, count int, value string) (int, error)
	Rpoplpush(ctx context.Context, key1 string, key2 string) (interface{}, error)

	Rpop(ctx context.Context, key string) (interface{}, error)

	Scard(ctx context.Context, key string) (int, error)
	Smembers(ctx context.Context, key string) (interface{}, error)
	Sismember(ctx context.Context, key string, value interface{}) (int, error)
	Srandmember(ctx context.Context, key string) (string, error)

	Zrank(ctx context.Context, key string, member interface{}) (int, error)
	Zscore(ctx context.Context, key string, member interface{}) (string, error)
	Zrange(ctx context.Context, key string, start int, stop int, withscores bool) (interface{}, error)
	Zrevrange(ctx context.Context, key string, start int, stop int, withscores bool) (interface{}, error)
	Zrangebyscore(ctx context.Context, key string, start int, stop int, withscores bool, limitOffset, limitCount int) (interface{}, error)
	Zrevrangebyscore(ctx context.Context, key string, start int, stop int, withscores bool, limitOffset, limitCount int) (interface{}, error)
	ZINCRBY(ctx context.Context, key string, increment int, member interface{}) (string, error)
	Zadd(ctx context.Context, key string, score interface{}, member interface{}) (int, error)
	Zrem(ctx context.Context, key string, member ...interface{}) (int, error)

	//Rename(ctx context.Context, key string, newkey string) (int, error) //pika 暂不支持
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
func InitPika(hsmIn map[string][]*config.ConnDetail, label ...string) chan *zgopika {
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

	InitPikaResource(hsm)

	//自动为变量初始化对象
	initLabel := ""
	for k, _ := range hsm {
		if k != "" {
			initLabel = k
			break
		}
	}
	out := make(chan *zgopika)
	go func() {
		in, err := GetPika(initLabel)
		if err != nil {
			out <- nil
		}
		out <- in
		close(out)
	}()
	return out
}

func (n *zgopika) New(label ...string) (*zgopika, error) {
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

func (r *zgopika) Set(ctx context.Context, key string, value interface{}) (string, error) {
	return r.res.Set(ctx, key, value)
}

func (r *zgopika) SetMutex(ctx context.Context, key string, ttl int, value interface{}) (string, error) {
	return r.res.SetMutex(ctx, key, ttl, value)
}

func (r *zgopika) Setnx(ctx context.Context, key string, value interface{}) (int, error) {
	return r.res.Setnx(ctx, key, value)
}

func (r *zgopika) Setex(ctx context.Context, key string, ttl int, value interface{}) (string, error) {
	return r.res.Setex(ctx, key, ttl, value)
}

func (r *zgopika) Expire(ctx context.Context, key string, time int) (int, error) {
	return r.res.Expire(ctx, key, time)
}

func (r *zgopika) Incrby(ctx context.Context, key string, val int) (interface{}, error) {
	return r.res.Incrby(ctx, key, val)
}

func (r *zgopika) Hset(ctx context.Context, key string, name string, value interface{}) (int, error) {
	return r.res.Hset(ctx, key, name, value)
}

func (r *zgopika) Hmset(ctx context.Context, key string, values interface{}) (string, error) {
	return r.res.Hmset(ctx, key, values)
}

func (r *zgopika) Lpush(ctx context.Context, key string, value interface{}) (int, error) {
	return r.res.Lpush(ctx, key, value)
}

func (r *zgopika) Rpush(ctx context.Context, key string, value interface{}) (int, error) {
	return r.res.Rpush(ctx, key, value)
}

func (r *zgopika) Rpoplpush(ctx context.Context, key1 string, key2 string) (interface{}, error) {
	return r.res.Rpoplpush(ctx, key1, key2)
}

func (r *zgopika) Sadd(ctx context.Context, key string, value interface{}) (int, error) {
	return r.res.Sadd(ctx, key, value)
}

func (r *zgopika) Srem(ctx context.Context, key string, value interface{}) (int, error) {
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

func (r *zgopika) Hlen(ctx context.Context, key string) (int, error) {
	return r.res.Hlen(ctx, key)
}

func (r *zgopika) Hdel(ctx context.Context, key string, name interface{}) (int, error) {
	return r.res.Hdel(ctx, key, name)
}

func (r *zgopika) Hgetall(ctx context.Context, key string) (interface{}, error) {
	return r.res.Hgetall(ctx, key)
}

func (r *zgopika) Hincrby(ctx context.Context, key, field string, inc int64) (int64, error) {
	return r.res.Hincrby(ctx, key, field, inc)
}

func (r *zgopika) Del(ctx context.Context, key string) (interface{}, error) {
	return r.res.Del(ctx, key)
}

func (r *zgopika) Llen(ctx context.Context, key string) (int, error) {
	return r.res.Llen(ctx, key)
}

func (r *zgopika) Lrange(ctx context.Context, key string, start int, stop int) (interface{}, error) {
	return r.res.Lrange(ctx, key, start, stop)
}

func (r *zgopika) Ltrim(ctx context.Context, key string, start int, stop int) (interface{}, error) {
	return r.res.Ltrim(ctx, key, start, stop)
}

func (r *zgopika) Lpop(ctx context.Context, key string) (interface{}, error) {
	return r.res.Lpop(ctx, key)
}

func (r *zgopika) Lrem(ctx context.Context, key string, count int, value string) (int, error) {
	return r.res.Lrem(ctx, key, count, value)
}

func (r *zgopika) Rpop(ctx context.Context, key string) (interface{}, error) {
	return r.res.Lpop(ctx, key)
}

func (r *zgopika) Scard(ctx context.Context, key string) (int, error) {
	return r.res.Scard(ctx, key)
}

func (r *zgopika) Smembers(ctx context.Context, key string) (interface{}, error) {
	return r.res.Smembers(ctx, key)
}

func (r *zgopika) Sismember(ctx context.Context, key string, value interface{}) (int, error) {
	return r.res.Sismember(ctx, key, value)
}
func (r *zgopika) Srandmember(ctx context.Context, key string) (string, error) {
	return r.res.Srandmember(ctx, key)
}
func (r *zgopika) Zrank(ctx context.Context, key string, member interface{}) (int, error) {
	return r.res.Zrank(ctx, key, member)
}

func (r *zgopika) Zscore(ctx context.Context, key string, member interface{}) (string, error) {
	return r.res.Zscore(ctx, key, member)
}

func (r *zgopika) Zrange(ctx context.Context, key string, start int, stop int, withscores bool) (interface{}, error) {
	return r.res.Zrange(ctx, key, start, stop, withscores)
}

func (r *zgopika) Zrevrange(ctx context.Context, key string, start int, stop int, withscores bool) (interface{}, error) {
	return r.res.Zrevrange(ctx, key, start, stop, withscores)
}

func (r *zgopika) Zrangebyscore(ctx context.Context, key string, start int, stop int, withscores bool, limitOffet, limitCount int) (interface{}, error) {
	return r.res.Zrangebyscore(ctx, key, start, stop, withscores, limitOffet, limitCount)
}

func (r *zgopika) Zrevrangebyscore(ctx context.Context, key string, start int, stop int, withscores bool, limitOffet, limitCount int) (interface{}, error) {
	return r.res.Zrevrangebyscore(ctx, key, start, stop, withscores, limitOffet, limitCount)
}

func (r *zgopika) ZINCRBY(ctx context.Context, key string, increment int, member interface{}) (string, error) {
	return r.res.ZINCRBY(ctx, key, increment, member)
}

func (r *zgopika) Zadd(ctx context.Context, key string, score interface{}, member interface{}) (int, error) {
	return r.res.Zadd(ctx, key, score, member)
}

func (r *zgopika) Zrem(ctx context.Context, key string, member ...interface{}) (int, error) {
	return r.res.Zrem(ctx, key, member)
}

//func (r *zgopika) Rename(ctx context.Context, key string, newkey string) (int, error) {
//	return r.res.Rename(ctx, key, newkey)
//}
