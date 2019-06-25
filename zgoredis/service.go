package zgoredis

import (
	"context"
	"errors"
	"git.zhugefang.com/gocore/zgo/comm"
	"git.zhugefang.com/gocore/zgo/config"
	"github.com/mediocregopher/radix"
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
	Hget(ctx context.Context, key string, name string) (interface{}, error)
	Ttl(ctx context.Context, key string) (interface{}, error)
	Type(ctx context.Context, key string) (interface{}, error)
	Hlen(ctx context.Context, key string) (int, error)
	Hdel(ctx context.Context, key string, name string) (interface{}, error)
	Hgetall(ctx context.Context, key string) (interface{}, error)
	Hincrby(ctx context.Context, key string, field string, inc int64) (int64, error)

	Del(ctx context.Context, key string) (interface{}, error)
	Llen(ctx context.Context, key string) (int, error)
	Lrange(ctx context.Context, key string, start int, stop int) (interface{}, error)
	Ltrim(ctx context.Context, key string, start int, stop int) (interface{}, error)
	Lpop(ctx context.Context, key string) (interface{}, error)
	Rpop(ctx context.Context, key string) (interface{}, error)
	Scard(ctx context.Context, key string) (int, error)
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

	// Publish 发布
	Publish(ctx context.Context, key string, value string) (int, error)

	// Subscribe订阅
	Subscribe(ctx context.Context, chanName string) (chan radix.PubSubMessage, error)

	// PSubscribe 模式订阅，模糊匹配channel的名字
	// PSubscribe(context.TODO(), "my*")
	PSubscribe(ctx context.Context, patterns ...string) (chan radix.PubSubMessage, error)

	// Unsubscribe取消订阅
	Unsubscribe(ctx context.Context, chanName string) (chan radix.PubSubMessage, error)

	// PUnsubscribe 取消模式订阅，模糊匹配channel的名字
	PUnsubscribe(ctx context.Context, patterns ...string) (chan radix.PubSubMessage, error)

	//streams流处理
	//	m := make(map[string]string)
	//	m["aaa"] = "aaa123"
	//	m["bbb"] = "bbb123"
	//	m["ccc"] = "ccc123"
	//
	//XAdd(context.TODO(), "key-101", "19000000000000", m)
	//id 的值必须后加入的要大于之前加入的，m可以是map[string]string
	XAdd(ctx context.Context, key string, id string, values interface{}) (string, error)

	XLen(ctx context.Context, key string) (int32, error)

	//	ids := []string{
	//		"1561371910929-0",
	//		"1561372099154-0",
	//		}
	//XDel(context.TODO(), "key-101", ids)
	//ids可以string数组，删除一个或多个
	XDel(ctx context.Context, key string, ids []string) (int32, error)

	//按范围取XRange(context.TODO(), "key-101", "1561372594375-0", "1561372671389-0")，count可选，输入1表示取1条
	//start可以是 - ; end可以是 + ; 表示全部
	XRange(ctx context.Context, key string, start, end string, count ...int) ([]map[string]map[string]string, error)

	//按范围反向取XRevrange(context.TODO(), "key-101", "1561372671389-0", "1561372594375-0")，count可选，输入1表示取1条
	//start可以是 + ; end可以是 - ; 表示全部
	XRevRange(ctx context.Context, key string, start, end string, count ...int) ([]map[string]map[string]string, error)

	//创建消费者
	//XGroupCreate(context.TODO(), "key-101", "group-101", "$")  $表示最后一个ID从最新的开始消费，也可以是0，表示从头开始，
	XGroupCreate(ctx context.Context, key string, groupName string, id string) (string, error)

	//删除消费者
	//XGroupDestroy(context.TODO(), "key-101", "group-101")
	XGroupDestroy(ctx context.Context, key string, groupName string) (int32, error)

	//	ids := []string{
	//		"19000000000010-0",
	//		"1561372594375-0",
	//	}
	//XAck(context.TODO(), "key-101", "group-101", ids)
	//ids可以string数组，确认一个或多个
	XAck(ctx context.Context, key string, groupName string, ids []string) (int32, error)

	/*
			var streamName = "key-101"
			var groupName = "group-102"

			zgo.Redis.XGroupCreate(context.TODO(), streamName, groupName, "0") //从0开始, $从最新开始

			streams := []string{
				streamName,
				"lol",
			}
			streamReader,err := zgo.Redis.NewStreamReader(streams, groupName, "101")
			if err != nil {
				zgo.Log.Error(err)
				return
		    }
			if streamReader.Err() != nil {
				zgo.Log.Error(streamReader.Err())
				return
			}

			for {
				if _, entries, ok := streamReader.Next(); ok == true {
					if len(entries) > 0 {
						fmt.Println(groupName, "===", entries)
					}
				}
			}
	*/
	//通过创建stream选项，来创建一个流的reader，然后在for中读到最新写进去的，默认xack为true, block为true
	NewStreamReader(streams []string, group, consumer string) (radix.StreamReader, error)
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

func (r *zgoredis) SetMutex(ctx context.Context, key string, ttl int, value interface{}) (string, error) {
	return r.res.SetMutex(ctx, key, ttl, value)
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

func (r *zgoredis) Incrby(ctx context.Context, key string, val int) (interface{}, error) {
	return r.res.Incrby(ctx, key, val)
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

func (r *zgoredis) Hlen(ctx context.Context, key string) (int, error) {
	return r.res.Hlen(ctx, key)
}

func (r *zgoredis) Hdel(ctx context.Context, key string, name string) (interface{}, error) {
	return r.res.Hdel(ctx, key, name)
}

func (r *zgoredis) Hgetall(ctx context.Context, key string) (interface{}, error) {
	return r.res.Hgetall(ctx, key)
}

func (r *zgoredis) Hincrby(ctx context.Context, key, field string, inc int64) (int64, error) {
	return r.res.Hincrby(ctx, key, field, inc)
}

func (r *zgoredis) Del(ctx context.Context, key string) (interface{}, error) {
	return r.res.Del(ctx, key)
}

func (r *zgoredis) Llen(ctx context.Context, key string) (int, error) {
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

func (r *zgoredis) Scard(ctx context.Context, key string) (int, error) {
	return r.res.Scard(ctx, key)
}

func (r *zgoredis) Smembers(ctx context.Context, key string) (interface{}, error) {
	return r.res.Smembers(ctx, key)
}

func (r *zgoredis) Sismember(ctx context.Context, key string, value interface{}) (int, error) {
	return r.res.Sismember(ctx, key, value)
}

func (r *zgoredis) Zrank(ctx context.Context, key string, member interface{}) (int, error) {
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

func (r *zgoredis) Zrangebyscore(ctx context.Context, key string, start int, stop int, withscores bool, limitOffet, limitCount int) (interface{}, error) {
	return r.res.Zrangebyscore(ctx, key, start, stop, withscores, limitOffet, limitCount)
}

func (r *zgoredis) Zrevrangebyscore(ctx context.Context, key string, start int, stop int, withscores bool, limitOffet, limitCount int) (interface{}, error) {
	return r.res.Zrevrangebyscore(ctx, key, start, stop, withscores, limitOffet, limitCount)
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

func (r *zgoredis) Publish(ctx context.Context, key string, value string) (int, error) {
	return r.res.Publish(ctx, key, value)
}

func (r *zgoredis) Subscribe(ctx context.Context, chanName string) (chan radix.PubSubMessage, error) {
	return r.res.Subscribe(ctx, chanName)
}

func (r *zgoredis) PSubscribe(ctx context.Context, patterns ...string) (chan radix.PubSubMessage, error) {
	return r.res.PSubscribe(ctx, patterns...)
}

func (r *zgoredis) Unsubscribe(ctx context.Context, chanName string) (chan radix.PubSubMessage, error) {
	return r.res.Unsubscribe(ctx, chanName)
}

func (r *zgoredis) PUnsubscribe(ctx context.Context, patterns ...string) (chan radix.PubSubMessage, error) {
	return r.res.PUnsubscribe(ctx, patterns...)
}

func (r *zgoredis) XAdd(ctx context.Context, key string, id string, values interface{}) (string, error) {
	return r.res.XAdd(ctx, key, id, values)
}

func (r *zgoredis) XLen(ctx context.Context, key string) (int32, error) {
	return r.res.XLen(ctx, key)
}

func (r *zgoredis) XDel(ctx context.Context, key string, ids []string) (int32, error) {
	if len(ids) <= 0 {
		return 0, errors.New("ID不能为空")
	}
	return r.res.XDel(ctx, key, ids)
}

func (r *zgoredis) XRange(ctx context.Context, key string, start, end string, count ...int) ([]map[string]map[string]string, error) {
	return r.res.XRange(ctx, key, start, end, count...)
}

func (r *zgoredis) XRevRange(ctx context.Context, key string, start, end string, count ...int) ([]map[string]map[string]string, error) {
	return r.res.XRevRange(ctx, key, start, end, count...)
}

func (r *zgoredis) XGroupCreate(ctx context.Context, key string, groupName string, id string) (string, error) {
	return r.res.XGroupCreate(ctx, key, groupName, id)
}

func (r *zgoredis) XGroupDestroy(ctx context.Context, key string, groupName string) (int32, error) {
	return r.res.XGroupDestroy(ctx, key, groupName)
}

func (r *zgoredis) XAck(ctx context.Context, key string, groupName string, ids []string) (int32, error) {
	if len(ids) <= 0 {
		return 0, errors.New("ID不能为空")
	}
	return r.res.XAck(ctx, key, groupName, ids)
}

func (r *zgoredis) NewStreamReader(streams []string, group, consumer string) (radix.StreamReader, error) {
	if len(streams) <= 0 {
		return nil, errors.New("stream key 至少有一个")
	}
	m := make(map[string]*radix.StreamEntryID)
	for _, v := range streams {
		m[v] = nil
	}
	opts := radix.StreamReaderOpts{
		Streams:  m,
		Group:    group,
		Consumer: consumer,
	}
	return r.res.NewStreamReader(opts), nil
}
