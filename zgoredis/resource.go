package zgoredis

import (
	"context"
	"errors"
	"git.zhugefang.com/gocore/zgo/config"
	"github.com/mediocregopher/radix"
	"sync"
)

//NsqResourcer 给service使用
type RedisResourcer interface {
	GetConnChan(label string) chan interface{}
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
	Exists(ctx context.Context, key string) (int, error)
	Get(ctx context.Context, key string) (interface{}, error)
	Keys(ctx context.Context, pattern string) (interface{}, error)
	//hget
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
	Rpop(ctx context.Context, key string) (interface{}, error)

	Scard(ctx context.Context, key string) (int, error)
	Sunion(ctx context.Context, key string, key1 ...interface{}) (interface{}, error)
	Smembers(ctx context.Context, key string) (interface{}, error)
	Sismember(ctx context.Context, key string, value interface{}) (int, error)

	Zrank(ctx context.Context, key string, member interface{}) (int, error)
	Zscore(ctx context.Context, key string, member interface{}) (string, error)
	Zrange(ctx context.Context, key string, start int, stop int, withscores bool) (interface{}, error)
	Zrevrange(ctx context.Context, key string, start int, stop int, withscores bool) (interface{}, error)
	Zrangebyscore(ctx context.Context, key string, start int, stop int, withscores bool, limitOffset, limitCount int) (interface{}, error)
	Zrevrangebyscore(ctx context.Context, key string, start int, stop int, withscores bool, limitOffset, limitCount int) (interface{}, error)
	ZINCRBY(ctx context.Context, key string, increment int, member interface{}) (string, error)
	Zadd(ctx context.Context, key string, score interface{}, member interface{}) (int, error)
	Zrem(ctx context.Context, key string, member ...interface{}) (int, error)
	Zremrangebyscore(ctx context.Context, key string, start int, stop int) (int, error)

	Publish(ctx context.Context, key string, value string) (int, error)
	Subscribe(ctx context.Context, chanName string) (chan radix.PubSubMessage, error)
	Unsubscribe(ctx context.Context, chanName string) (chan radix.PubSubMessage, error)
	PSubscribe(ctx context.Context, patterns ...string) (chan radix.PubSubMessage, error)
	PUnsubscribe(ctx context.Context, patterns ...string) (chan radix.PubSubMessage, error)

	//streams
	XAdd(ctx context.Context, key string, id string, values interface{}) (string, error)
	XLen(ctx context.Context, key string) (int32, error)
	XDel(ctx context.Context, key string, ids []string) (int32, error)
	XRange(ctx context.Context, key string, start, end string, count ...int) ([]map[string]map[string]string, error)
	XRevRange(ctx context.Context, key string, start, end string, count ...int) ([]map[string]map[string]string, error)
	XGroupCreate(ctx context.Context, key string, groupName string, id string) (string, error)
	XGroupDestroy(ctx context.Context, key string, groupName string) (int32, error)
	XAck(ctx context.Context, key string, groupName string, ids []string) (int32, error)
	NewStreamReader(opts radix.StreamReaderOpts) radix.StreamReader
}

type redisResource struct {
	label    string
	mu       sync.RWMutex
	connpool ConnPooler
}

func InitRedisResource(hsm map[string][]*config.ConnDetail) {
	InitConnPool(hsm)
}

func NewRedisResource(label string) RedisResourcer {
	return &redisResource{
		label:    label,
		connpool: NewConnPool(label)}
}

//GetConnChan 返回存放连接的chan
func (r *redisResource) GetConnChan(label string) chan interface{} {
	return r.connpool.GetConnChan(label)
}

func (r *redisResource) Set(ctx context.Context, key string, value interface{}) (string, error) {

	var result string
	flatCmd := radix.FlatCmd(&result, "SET", key, value)
	err := r.deal(flatCmd)
	return result, err
}

func (r *redisResource) SetMutex(ctx context.Context, key string, ttl int, value interface{}) (string, error) {

	var result string
	flatCmd := radix.FlatCmd(&result, "SET", key, value, "EX", ttl, "NX")
	err := r.deal(flatCmd)
	return result, err
}

func (r *redisResource) Setnx(ctx context.Context, key string, value interface{}) (int, error) {

	var result int
	flatCmd := radix.FlatCmd(&result, "SETNX", key, value)
	err := r.deal(flatCmd)
	return result, err
}

func (r *redisResource) Setex(ctx context.Context, key string, ttl int, value interface{}) (string, error) {

	var result string
	flatCmd := radix.FlatCmd(&result, "SETEX", key, ttl, value)
	err := r.deal(flatCmd)
	return result, err
}

func (r *redisResource) Expire(ctx context.Context, key string, ttl int) (int, error) {

	var result int
	flatCmd := radix.FlatCmd(&result, "Expire", key, ttl)
	err := r.deal(flatCmd)
	return result, err
}

func (r *redisResource) Incrby(ctx context.Context, key string, val int) (interface{}, error) {

	var result string
	flatCmd := radix.FlatCmd(&result, "Incrby", key, val)
	err := r.deal(flatCmd)
	return result, err
}

func (r *redisResource) Hset(ctx context.Context, key string, field string, value interface{}) (int, error) {

	var result int
	flatCmd := radix.FlatCmd(&result, "Hset", key, field, value)
	err := r.deal(flatCmd)
	return result, err
}

func (r *redisResource) Hmset(ctx context.Context, key string, values interface{}) (string, error) {

	var result string
	flatCmd := radix.FlatCmd(&result, "HMSET", key, values)
	err := r.deal(flatCmd)
	return result, err
}

func (r *redisResource) Lpush(ctx context.Context, key string, value interface{}) (int, error) {

	var result int
	flatCmd := radix.FlatCmd(&result, "Lpush", key, value)
	err := r.deal(flatCmd)
	return result, err
}

func (r *redisResource) Rpush(ctx context.Context, key string, value interface{}) (int, error) {

	var result int
	flatCmd := radix.FlatCmd(&result, "Rpush", key, value)
	err := r.deal(flatCmd)
	return result, err
}

func (r *redisResource) Sadd(ctx context.Context, key string, value interface{}) (int, error) {

	var result int
	flatCmd := radix.FlatCmd(&result, "Sadd", key, value)
	err := r.deal(flatCmd)
	return result, err
}

func (r *redisResource) Srem(ctx context.Context, key string, value interface{}) (int, error) {

	var result int
	flatCmd := radix.FlatCmd(&result, "Srem", key, value)
	err := r.deal(flatCmd)
	return result, err
}

func (r *redisResource) Exists(ctx context.Context, key string) (int, error) {

	var result int
	flatCmd := radix.FlatCmd(&result, "Exists", key)
	err := r.deal(flatCmd)
	return result, err
}

func (r *redisResource) Get(ctx context.Context, key string) (interface{}, error) {

	var result string
	flatCmd := radix.FlatCmd(&result, "Get", key)
	err := r.deal(flatCmd)
	return result, err
}

func (r *redisResource) Keys(ctx context.Context, key string) (interface{}, error) {

	if key == "*" {
		return nil, errors.New("forbidden")
	}
	var result []string
	flatCmd := radix.FlatCmd(&result, "Keys", key)
	err := r.deal(flatCmd)
	return result, err
}

func (r *redisResource) Hget(ctx context.Context, key string, name string) (interface{}, error) {

	var result string
	flatCmd := radix.FlatCmd(&result, "Hget", key, name)
	err := r.deal(flatCmd)
	return result, err
}

func (r *redisResource) Ttl(ctx context.Context, key string) (interface{}, error) {

	var result int
	flatCmd := radix.FlatCmd(&result, "ttl", key)
	err := r.deal(flatCmd)
	return result, err
}

func (r *redisResource) Type(ctx context.Context, key string) (interface{}, error) {

	var result interface{}
	flatCmd := radix.FlatCmd(&result, "Type", key)
	err := r.deal(flatCmd)
	return result, err
}

func (r *redisResource) Hlen(ctx context.Context, key string) (int, error) {

	var result int
	flatCmd := radix.FlatCmd(&result, "Hlen", key)
	err := r.deal(flatCmd)
	return result, err
}

func (r *redisResource) Hdel(ctx context.Context, key string, name interface{}) (int, error) {

	var result int
	flatCmd := radix.FlatCmd(&result, "Hdel", key, name)
	err := r.deal(flatCmd)
	return result, err
}

func (r *redisResource) Hgetall(ctx context.Context, key string) (interface{}, error) {

	var result map[string]string
	flatCmd := radix.FlatCmd(&result, "Hgetall", key)
	err := r.deal(flatCmd)
	return result, err
}

func (r *redisResource) Hincrby(ctx context.Context, key, field string, inc int64) (int64, error) {

	var result int64
	flatCmd := radix.FlatCmd(&result, "HINCRBY", key, field, inc)
	err := r.deal(flatCmd)
	return result, err
}

func (r *redisResource) Del(ctx context.Context, key string) (interface{}, error) {

	var result int
	flatCmd := radix.FlatCmd(&result, "del", key)
	err := r.deal(flatCmd)
	return result, err
}

func (r *redisResource) Llen(ctx context.Context, key string) (int, error) {

	var result int
	flatCmd := radix.FlatCmd(&result, "Llen", key)
	err := r.deal(flatCmd)
	return result, err
}

func (r *redisResource) Lrange(ctx context.Context, key string, start int, stop int) (interface{}, error) {

	var result []string
	flatCmd := radix.FlatCmd(&result, "Lrange", key, start, stop)
	err := r.deal(flatCmd)
	return result, err
}

func (r *redisResource) Ltrim(ctx context.Context, key string, start int, stop int) (interface{}, error) {

	var result interface{}
	flatCmd := radix.FlatCmd(&result, "Ltrim", key, start, stop)
	err := r.deal(flatCmd)
	return result, err
}

func (r *redisResource) Lpop(ctx context.Context, key string) (interface{}, error) {

	var result string
	flatCmd := radix.FlatCmd(&result, "Lpop", key)
	err := r.deal(flatCmd)
	return result, err
}

func (r *redisResource) Rpop(ctx context.Context, key string) (interface{}, error) {

	var result string
	flatCmd := radix.FlatCmd(&result, "Rpop", key)
	err := r.deal(flatCmd)
	return result, err
}

func (r *redisResource) Scard(ctx context.Context, key string) (int, error) {

	var result int
	flatCmd := radix.FlatCmd(&result, "Scard", key)
	err := r.deal(flatCmd)
	return result, err
}
func (r *redisResource) Sunion(ctx context.Context, key string, key1 ...interface{}) (interface{}, error) {
	var result []string
	flatCmd := radix.FlatCmd(&result, "Sunion", key, key1...)
	err := r.deal(flatCmd)
	return result, err
}

func (r *redisResource) Smembers(ctx context.Context, key string) (interface{}, error) {

	var result []string
	flatCmd := radix.FlatCmd(&result, "Smembers", key)
	err := r.deal(flatCmd)
	return result, err
}

func (r *redisResource) Sismember(ctx context.Context, key string, value interface{}) (int, error) {

	var result int
	flatCmd := radix.FlatCmd(&result, "Sismember", key, value)
	err := r.deal(flatCmd)
	return result, err
}

func (r *redisResource) Zrank(ctx context.Context, key string, member interface{}) (int, error) {

	var result int
	flatCmd := radix.FlatCmd(&result, "Zrank", key, member)
	err := r.deal(flatCmd)
	return result, err
}

func (r *redisResource) Zscore(ctx context.Context, key string, member interface{}) (string, error) {

	var result string
	flatCmd := radix.FlatCmd(&result, "Zscore", key, member)
	err := r.deal(flatCmd)
	return result, err
}

func (r *redisResource) Zrange(ctx context.Context, key string, start int, stop int, withscores bool) (interface{}, error) {

	var result []string
	var flatCmd radix.CmdAction
	if withscores {
		flatCmd = radix.FlatCmd(&result, "Zrange", key, start, stop, "WITHSCORES")
	} else {
		flatCmd = radix.FlatCmd(&result, "Zrange", key, start, stop)
	}
	err := r.deal(flatCmd)
	return result, err
}

func (r *redisResource) Zrevrange(ctx context.Context, key string, start int, stop int, withscores bool) (interface{}, error) {

	var result []string
	var flatCmd radix.CmdAction
	if withscores {
		flatCmd = radix.FlatCmd(&result, "Zrevrange", key, start, stop, "WITHSCORES")
	} else {
		flatCmd = radix.FlatCmd(&result, "Zrevrange", key, start, stop)
	}
	err := r.deal(flatCmd)
	return result, err
}

func (r *redisResource) Zrangebyscore(ctx context.Context, key string, start int, stop int, withscores bool, limitOffset, limitCount int) (interface{}, error) {
	var result []string
	var flatCmd radix.CmdAction
	if withscores {
		flatCmd = radix.FlatCmd(&result, "ZRANGEBYSCORE", key, start, stop, "WITHSCORES", "LIMIT", limitOffset, limitCount)
	} else {
		flatCmd = radix.FlatCmd(&result, "ZRANGEBYSCORE", key, start, stop, "LIMIT", limitOffset, limitCount)
	}
	err := r.deal(flatCmd)
	return result, err
}

func (r *redisResource) Zrevrangebyscore(ctx context.Context, key string, start int, stop int, withscores bool, limitOffset, limitCount int) (interface{}, error) {

	var result []string
	var flatCmd radix.CmdAction
	if withscores {
		flatCmd = radix.FlatCmd(&result, "ZREVRANGEBYSCORE", key, start, stop, "WITHSCORES", "LIMIT", limitOffset, limitCount)
	} else {
		flatCmd = radix.FlatCmd(&result, "ZREVRANGEBYSCORE", key, start, stop, "LIMIT", limitOffset, limitCount)
	}
	err := r.deal(flatCmd)
	return result, err
}

func (r *redisResource) ZINCRBY(ctx context.Context, key string, increment int, member interface{}) (string, error) {

	var result string
	flatCmd := radix.FlatCmd(&result, "Zincrby", key, increment, member)
	err := r.deal(flatCmd)
	return result, err
}

func (r *redisResource) Zadd(ctx context.Context, key string, score interface{}, member interface{}) (int, error) {

	var result int
	flatCmd := radix.FlatCmd(&result, "Zadd", key, score, member)
	err := r.deal(flatCmd)
	return result, err
}

func (r *redisResource) Zrem(ctx context.Context, key string, member ...interface{}) (int, error) {
	var result int
	flatCmd := radix.FlatCmd(&result, "Zrem", key, member...)
	err := r.deal(flatCmd)
	return result, err
}

func (r *redisResource) Zremrangebyscore(ctx context.Context, key string, start int, stop int) (int, error) {
	var result int
	flatCmd := radix.FlatCmd(&result, "Zremrangebyscore", key, start, stop)
	err := r.deal(flatCmd)
	return result, err
}

// Publish 发布
func (r *redisResource) Publish(ctx context.Context, key string, value string) (int, error) {

	var result int
	flatCmd := radix.FlatCmd(&result, "PUBLISH", key, value)
	err := r.deal(flatCmd)
	return result, err
}

// Subscribe订阅
func (r *redisResource) Subscribe(ctx context.Context, chanName string) (chan radix.PubSubMessage, error) {
	s := <-r.connpool.GetCChan(r.label)
	ps := radix.PubSub(*s)
	msgCh := make(chan radix.PubSubMessage)
	if err := ps.Subscribe(msgCh, chanName); err == nil {
		return msgCh, err
	} else {
		return nil, err
	}

}

// PSubscribe 模式订阅，模糊匹配channel的名字
func (r *redisResource) PSubscribe(ctx context.Context, patterns ...string) (chan radix.PubSubMessage, error) {
	s := <-r.connpool.GetCChan(r.label)
	ps := radix.PubSub(*s)
	msgCh := make(chan radix.PubSubMessage)
	if err := ps.PSubscribe(msgCh, patterns...); err == nil {
		return msgCh, err
	} else {
		return nil, err
	}

}

// Unsubscribe 取消订阅
func (r *redisResource) Unsubscribe(ctx context.Context, chanName string) (chan radix.PubSubMessage, error) {
	s := <-r.connpool.GetCChan(r.label)
	ps := radix.PubSub(*s)
	msgCh := make(chan radix.PubSubMessage)
	if err := ps.Unsubscribe(msgCh, chanName); err == nil {
		return msgCh, err
	} else {
		return nil, err
	}

}

// PUnsubscribe 取消模式订阅，模糊匹配channel的名字
func (r *redisResource) PUnsubscribe(ctx context.Context, patterns ...string) (chan radix.PubSubMessage, error) {
	s := <-r.connpool.GetCChan(r.label)
	ps := radix.PubSub(*s)
	msgCh := make(chan radix.PubSubMessage)
	if err := ps.PUnsubscribe(msgCh, patterns...); err == nil {
		return msgCh, err
	} else {
		return nil, err
	}

}

func (r *redisResource) XAdd(ctx context.Context, key string, id string, values interface{}) (string, error) {

	var result string
	flatCmd := radix.FlatCmd(&result, "XADD", key, id, values)
	err := r.deal(flatCmd)
	return result, err
}

func (r *redisResource) XLen(ctx context.Context, key string) (int32, error) {

	var result int32
	flatCmd := radix.FlatCmd(&result, "XLEN", key)
	err := r.deal(flatCmd)
	return result, err
}

func (r *redisResource) XDel(ctx context.Context, key string, ids []string) (int32, error) {

	var result int32
	flatCmd := radix.FlatCmd(&result, "XDEL", key, ids)
	err := r.deal(flatCmd)
	return result, err
}

func (r *redisResource) XRange(ctx context.Context, key string, start, end string, count ...int) ([]map[string]map[string]string, error) {

	var result []map[string]map[string]string
	var flatCmd radix.CmdAction
	if len(count) > 0 {
		flatCmd = radix.FlatCmd(&result, "XRANGE", key, start, end, "COUNT", count[0])
	} else {
		flatCmd = radix.FlatCmd(&result, "XRANGE", key, start, end)
	}
	err := r.deal(flatCmd)
	return result, err
}

func (r *redisResource) XRevRange(ctx context.Context, key string, start, end string, count ...int) ([]map[string]map[string]string, error) {

	var result []map[string]map[string]string
	var flatCmd radix.CmdAction
	if len(count) > 0 {
		flatCmd = radix.FlatCmd(&result, "XREVRANGE", key, start, end, "COUNT", count[0])
	} else {
		flatCmd = radix.FlatCmd(&result, "XREVRANGE", key, start, end)
	}
	err := r.deal(flatCmd)
	return result, err
}

func (r *redisResource) XGroupCreate(ctx context.Context, key string, groupName string, id string) (string, error) {

	var result string
	flatCmd := radix.FlatCmd(&result, "XGROUP", "CREATE", key, groupName, id)
	err := r.deal(flatCmd)
	return result, err
}

func (r *redisResource) XGroupDestroy(ctx context.Context, key string, groupName string) (int32, error) {

	var result int32
	flatCmd := radix.FlatCmd(&result, "XGROUP", "DESTROY", key, groupName)
	err := r.deal(flatCmd)
	return result, err
}

func (r *redisResource) XAck(ctx context.Context, key string, groupName string, ids []string) (int32, error) {

	var result int32
	flatCmd := radix.FlatCmd(&result, "XACK", key, groupName, ids)
	err := r.deal(flatCmd)
	return result, err
}

func (r *redisResource) NewStreamReader(opts radix.StreamReaderOpts) radix.StreamReader {
	sn := <-r.connpool.GetConnChan(r.label)
	switch s := sn.(type) {
	case *radix.Pool:
		return radix.NewStreamReader(s, opts)
	case *radix.Cluster:
		return radix.NewStreamReader(s, opts)
	}
	return nil
}

func (r *redisResource) deal(flatCmd radix.CmdAction) error {
	sn := <-r.connpool.GetConnChan(r.label)
	var err error
	switch s := sn.(type) {
	case *radix.Pool:
		err = s.Do(flatCmd)
	case *radix.Cluster:
		err = s.Do(flatCmd)
	}
	return err
}
