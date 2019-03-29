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
	GetConnChan(label string) chan *radix.Pool
	//Post
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
	//hget
	Hget(ctx context.Context, key string, name string) (interface{}, error)
	Ttl(ctx context.Context, key string) (interface{}, error)
	Type(ctx context.Context, key string) (interface{}, error)
	Hlen(ctx context.Context, key string) (interface{}, error)
	Hdel(ctx context.Context, key string, name interface{}) (int, error)
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

	Zrank(ctx context.Context, key string, member interface{}) (int, error)
	Zscore(ctx context.Context, key string, member interface{}) (string, error)
	Zrange(ctx context.Context, key string, start int, stop int, withscores bool) (interface{}, error)
	Zrevrange(ctx context.Context, key string, start int, stop int, withscores bool) (interface{}, error)
	ZINCRBY(ctx context.Context, key string, increment int, member interface{}) (string, error)
	Zadd(ctx context.Context, key string, score interface{}, member interface{}) (int, error)
	Zrem(ctx context.Context, key string, member ...interface{}) (int, error)
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
func (r *redisResource) GetConnChan(label string) chan *radix.Pool {
	return r.connpool.GetConnChan(label)
}

func (r *redisResource) Set(ctx context.Context, key string, value interface{}) (string, error) {
	s := <-r.connpool.GetConnChan(r.label)
	var res string
	if err := s.Do(radix.FlatCmd(&res, "SET", key, value)); err != nil {
		return "", err
	} else {
		return res, err
	}
	return res, nil
}

func (r *redisResource) Setnx(ctx context.Context, key string, value interface{}) (int, error) {
	s := <-r.connpool.GetConnChan(r.label)
	var res int
	if err := s.Do(radix.FlatCmd(&res, "SETNX", key, value)); err != nil {
		return 0, err
	} else {
		return res, err
	}
	return res, nil
}

func (r *redisResource) Setex(ctx context.Context, key string, ttl int, value interface{}) (string, error) {
	s := <-r.connpool.GetConnChan(r.label)
	var res string
	if err := s.Do(radix.FlatCmd(&res, "SETEX", key, ttl, value)); err != nil {
		return "", err
	} else {
		return res, err
	}
	return res, nil
}

func (r *redisResource) Expire(ctx context.Context, key string, time int) (int, error) {
	s := <-r.connpool.GetConnChan(r.label)
	var res int
	if err := s.Do(radix.FlatCmd(&res, "Expire", key, time)); err != nil {
		return 0, err
	} else {
		return res, err
	}
	return res, nil
}

func (r *redisResource) Hset(ctx context.Context, key string, name string, value interface{}) (int, error) {
	s := <-r.connpool.GetConnChan(r.label)
	var res int
	if err := s.Do(radix.FlatCmd(&res, "Hset", key, name, value)); err != nil {
		return 0, err
	} else {
		return res, err
	}
	return res, nil
}

func (r *redisResource) Hmset(ctx context.Context, key string, values interface{}) (string, error) {
	s := <-r.connpool.GetConnChan(r.label)
	var res string
	if err := s.Do(radix.FlatCmd(&res, "HMSET", key, values)); err != nil {
		return "", err
	} else {
		return res, err
	}
	return res, nil
}

func (r *redisResource) Lpush(ctx context.Context, key string, value interface{}) (int, error) {
	s := <-r.connpool.GetConnChan(r.label)
	var res int
	if err := s.Do(radix.FlatCmd(&res, "Lpush", key, value)); err != nil {
		return 0, err
	} else {
		return res, err
	}
	return res, nil
}

func (r *redisResource) Rpush(ctx context.Context, key string, value interface{}) (int, error) {
	s := <-r.connpool.GetConnChan(r.label)
	var res int
	if err := s.Do(radix.FlatCmd(&res, "Rpush", key, value)); err != nil {
		return 0, err
	} else {
		return res, err
	}
	return res, nil
}

func (r *redisResource) Sadd(ctx context.Context, key string, value interface{}) (int, error) {
	s := <-r.connpool.GetConnChan(r.label)
	var res int
	if err := s.Do(radix.FlatCmd(&res, "Sadd", key, value)); err != nil {
		return 0, err
	} else {
		return res, err
	}
	return res, nil
}

func (r *redisResource) Srem(ctx context.Context, key string, value interface{}) (int, error) {
	s := <-r.connpool.GetConnChan(r.label)
	var res int
	if err := s.Do(radix.FlatCmd(&res, "Srem", key, value)); err != nil {
		return 0, err
	} else {
		return res, err
	}
	return res, nil
}

func (r *redisResource) Exists(ctx context.Context, key string) (int, error) {
	s := <-r.connpool.GetConnChan(r.label)
	var flag int
	if err := s.Do(radix.Cmd(&flag, "Exists", key)); err != nil {
		return 0, err
	} else {
		return flag, err
	}
}

func (r *redisResource) Get(ctx context.Context, key string) (interface{}, error) {
	s := <-r.connpool.GetConnChan(r.label)
	var fooVal string
	if err := s.Do(radix.FlatCmd(&fooVal, "Get", key)); err != nil {
		return nil, err
	} else {
		return fooVal, err
	}
}

func (r *redisResource) Keys(ctx context.Context, key string) (interface{}, error) {
	s := <-r.connpool.GetConnChan(r.label)
	if key == "*" {
		return nil, errors.New("forbidden")
	}
	var bazEls []string
	if err := s.Do(radix.FlatCmd(&bazEls, "Keys", key)); err != nil {
		return nil, err
	} else {
		return bazEls, err
	}
}

func (r *redisResource) Hget(ctx context.Context, key string, name string) (interface{}, error) {
	s := <-r.connpool.GetConnChan(r.label)
	var fooVal string
	if err := s.Do(radix.FlatCmd(&fooVal, "Hget", key, name)); err != nil {
		return nil, err
	} else {
		return fooVal, err
	}
}

func (r *redisResource) Ttl(ctx context.Context, key string) (interface{}, error) {
	s := <-r.connpool.GetConnChan(r.label)
	var intervltime int
	if err := s.Do(radix.FlatCmd(&intervltime, "Ttl", key)); err != nil {
		return nil, err
	} else {
		return intervltime, err
	}
}

func (r *redisResource) Type(ctx context.Context, key string) (interface{}, error) {
	s := <-r.connpool.GetConnChan(r.label)
	var dataType interface{}
	if err := s.Do(radix.FlatCmd(&dataType, "Type", key)); err != nil {
		return nil, err
	} else {
		return dataType, err
	}
}

func (r *redisResource) Hlen(ctx context.Context, key string) (interface{}, error) {
	s := <-r.connpool.GetConnChan(r.label)
	var dataLen int
	if err := s.Do(radix.FlatCmd(&dataLen, "Hlen", key)); err != nil {
		return nil, err
	} else {
		return dataLen, err
	}
}

func (r *redisResource) Hdel(ctx context.Context, key string, name interface{}) (int, error) {
	s := <-r.connpool.GetConnChan(r.label)
	var flag int
	if err := s.Do(radix.FlatCmd(&flag, "Hdel", key, name)); err != nil {
		return 0, err
	} else {
		return flag, err
	}
}

func (r *redisResource) Hgetall(ctx context.Context, key string) (interface{}, error) {
	s := <-r.connpool.GetConnChan(r.label)
	var buzMap map[string]string
	if err := s.Do(radix.FlatCmd(&buzMap, "Hgetall", key)); err != nil {
		return nil, err
	} else {
		return buzMap, err
	}
}

func (r *redisResource) Del(ctx context.Context, key string) (interface{}, error) {
	s := <-r.connpool.GetConnChan(r.label)
	var flag int
	if err := s.Do(radix.FlatCmd(&flag, "del", key)); err != nil {
		return nil, err
	} else {
		return flag, err
	}
}

func (r *redisResource) Llen(ctx context.Context, key string) (interface{}, error) {
	s := <-r.connpool.GetConnChan(r.label)
	var dataLen int
	if err := s.Do(radix.FlatCmd(&dataLen, "Llen", key)); err != nil {
		return nil, err
	} else {
		return dataLen, err
	}
}

func (r *redisResource) Lrange(ctx context.Context, key string, start int, stop int) (interface{}, error) {
	s := <-r.connpool.GetConnChan(r.label)
	var listContent []string
	if err := s.Do(radix.FlatCmd(&listContent, "Lrange", key, start, stop)); err != nil {
		return nil, err
	} else {
		return listContent, err
	}
}

func (r *redisResource) Ltrim(ctx context.Context, key string, start int, stop int) (interface{}, error) {
	s := <-r.connpool.GetConnChan(r.label)
	var listContent interface{}
	if err := s.Do(radix.FlatCmd(&listContent, "Ltrim", key, start, stop)); err != nil {
		return nil, err
	} else {
		return listContent, err
	}
}

func (r *redisResource) Lpop(ctx context.Context, key string) (interface{}, error) {
	s := <-r.connpool.GetConnChan(r.label)
	var listContent string
	if err := s.Do(radix.FlatCmd(&listContent, "Lpop", key)); err != nil {
		return nil, err
	} else {
		return listContent, err
	}
}

func (r *redisResource) Rpop(ctx context.Context, key string) (interface{}, error) {
	s := <-r.connpool.GetConnChan(r.label)
	var listContent string
	if err := s.Do(radix.FlatCmd(&listContent, "Rpop", key)); err != nil {
		return nil, err
	} else {
		return listContent, err
	}
}

func (r *redisResource) Scard(ctx context.Context, key string) (interface{}, error) {
	s := <-r.connpool.GetConnChan(r.label)
	var setLen int
	if err := s.Do(radix.FlatCmd(&setLen, "Scard", key)); err != nil {
		return nil, err
	} else {
		return setLen, err
	}
}

func (r *redisResource) Smembers(ctx context.Context, key string) (interface{}, error) {
	s := <-r.connpool.GetConnChan(r.label)
	var setContent []string
	if err := s.Do(radix.FlatCmd(&setContent, "Smembers", key)); err != nil {
		return nil, err
	} else {
		return setContent, err
	}
}

func (r *redisResource) Sismember(ctx context.Context, key string, value interface{}) (int, error) {
	s := <-r.connpool.GetConnChan(r.label)
	var flag int
	if err := s.Do(radix.FlatCmd(&flag, "Sismember", key, value)); err != nil {
		return 0, err
	} else {
		return flag, err
	}
}

func (r *redisResource) Zrank(ctx context.Context, key string, member interface{}) (int, error) {
	s := <-r.connpool.GetConnChan(r.label)
	var rank int
	if err := s.Do(radix.FlatCmd(&rank, "Zrank", key, member)); err != nil {
		return 0, err
	} else {
		return rank, err
	}
}

func (r *redisResource) Zscore(ctx context.Context, key string, member interface{}) (string, error) {
	s := <-r.connpool.GetConnChan(r.label)
	var score string
	if err := s.Do(radix.FlatCmd(&score, "Zscore", key, member)); err != nil {
		return "", err
	} else {
		return score, err
	}
}

func (r *redisResource) Zrange(ctx context.Context, key string, start int, stop int, withscores bool) (interface{}, error) {
	s := <-r.connpool.GetConnChan(r.label)
	var setContent []string
	if withscores {
		if err := s.Do(radix.FlatCmd(&setContent, "Zrange", key, start, stop, "WITHSCORES")); err != nil {
			return nil, err
		} else {
			return setContent, err
		}
	} else {
		if err := s.Do(radix.FlatCmd(&setContent, "Zrange", key, start, stop)); err != nil {
			return nil, err
		} else {
			return setContent, err
		}
	}

}

func (r *redisResource) Zrevrange(ctx context.Context, key string, start int, stop int, withscores bool) (interface{}, error) {
	s := <-r.connpool.GetConnChan(r.label)
	var setContent []string
	if withscores {
		if err := s.Do(radix.FlatCmd(&setContent, "Zrevrange", key, start, stop, "WITHSCORES")); err != nil {
			return nil, err
		} else {
			return setContent, err
		}
	} else {
		if err := s.Do(radix.FlatCmd(&setContent, "Zrevrange", key, start, stop)); err != nil {
			return nil, err
		} else {
			return setContent, err
		}
	}

}

func (r *redisResource) ZINCRBY(ctx context.Context, key string, increment int, member interface{}) (string, error) {
	s := <-r.connpool.GetConnChan(r.label)
	var score string
	if err := s.Do(radix.FlatCmd(&score, "Zincrby", key, increment, member)); err != nil {
		return "", err
	} else {
		return score, err
	}
}

func (r *redisResource) Zadd(ctx context.Context, key string, score interface{}, member interface{}) (int, error) {
	s := <-r.connpool.GetConnChan(r.label)
	var flag int
	if err := s.Do(radix.FlatCmd(&flag, "Zadd", key, score, member)); err != nil {
		return 0, err
	} else {
		return flag, err
	}
}

func (r *redisResource) Zrem(ctx context.Context, key string, member ...interface{}) (int, error) {
	s := <-r.connpool.GetConnChan(r.label)
	var flag int
	if err := s.Do(radix.FlatCmd(&flag, "Zrem", key, member...)); err != nil {
		return 0, err
	} else {
		return flag, err
	}
}
