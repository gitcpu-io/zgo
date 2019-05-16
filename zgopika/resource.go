package zgopika

import (
	"context"
	"errors"
	"git.zhugefang.com/gocore/zgo/config"
	"github.com/mediocregopher/radix"
	"sync"
)

//NsqResourcer 给service使用
type PikaResourcer interface {
	GetConnChan(label string) chan *radix.Pool
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
}

type pikaResource struct {
	label    string
	mu       sync.RWMutex
	connpool ConnPooler
}

func InitPikaResource(hsm map[string][]*config.ConnDetail) {
	InitConnPool(hsm)
}

func NewPikaResource(label string) PikaResourcer {
	return &pikaResource{
		label:    label,
		connpool: NewConnPool(label),
	}

}

//GetConnChan 返回存放连接的chan
func (p *pikaResource) GetConnChan(label string) chan *radix.Pool {
	return p.connpool.GetConnChan(label)
}

func (p *pikaResource) Set(ctx context.Context, key string, value interface{}) (string, error) {
	s := <-p.connpool.GetConnChan(p.label)
	prefix := p.connpool.GetPrefix(p.label)
	key = prefix + key
	var res string
	if err := s.Do(radix.FlatCmd(&res, "SET", key, value)); err != nil {
		return "", err
	} else {
		return res, err
	}
	return res, nil
}

func (p *pikaResource) SetMutex(ctx context.Context, key string, ttl int, value interface{}) (string, error) {
	s := <-p.connpool.GetConnChan(p.label)
	prefix := p.connpool.GetPrefix(p.label)
	key = prefix + key
	var res string
	if err := s.Do(radix.FlatCmd(&res, "SET", key, value, "EX", ttl, "NX")); err != nil {
		return "", err
	} else {
		return res, err
	}
	return res, nil
}

func (p *pikaResource) Setnx(ctx context.Context, key string, value interface{}) (int, error) {
	s := <-p.connpool.GetConnChan(p.label)
	prefix := p.connpool.GetPrefix(p.label)
	key = prefix + key
	var res int
	if err := s.Do(radix.FlatCmd(&res, "SETNX", key, value)); err != nil {
		return 0, err
	} else {
		return res, err
	}
	return res, nil
}

func (p *pikaResource) Setex(ctx context.Context, key string, ttl int, value interface{}) (string, error) {
	s := <-p.connpool.GetConnChan(p.label)
	prefix := p.connpool.GetPrefix(p.label)
	key = prefix + key
	var res string
	if err := s.Do(radix.FlatCmd(&res, "SETEX", key, ttl, value)); err != nil {
		return "", err
	} else {
		return res, err
	}
	return res, nil
}

func (p *pikaResource) Expire(ctx context.Context, key string, time int) (int, error) {
	s := <-p.connpool.GetConnChan(p.label)
	prefix := p.connpool.GetPrefix(p.label)
	key = prefix + key
	var res int
	if err := s.Do(radix.FlatCmd(&res, "Expire", key, time)); err != nil {
		return 0, err
	} else {
		return res, err
	}
	return res, nil
}

func (p *pikaResource) Incrby(ctx context.Context, key string, val int) (interface{}, error) {
	s := <-p.connpool.GetConnChan(p.label)
	prefix := p.connpool.GetPrefix(p.label)
	key = prefix + key
	var res string
	if err := s.Do(radix.FlatCmd(&res, "Incrby", key, val)); err != nil {
		return "", err
	} else {
		return res, err
	}
	return res, nil
}

func (p *pikaResource) Hset(ctx context.Context, key string, name string, value interface{}) (int, error) {
	s := <-p.connpool.GetConnChan(p.label)
	prefix := p.connpool.GetPrefix(p.label)
	key = prefix + key
	var res int
	if err := s.Do(radix.FlatCmd(&res, "Hset", key, name, value)); err != nil {
		return 0, err
	} else {
		return res, err
	}
	return res, nil
}

func (r *pikaResource) Hmset(ctx context.Context, key string, values interface{}) (string, error) {
	s := <-r.connpool.GetConnChan(r.label)
	var res string
	if err := s.Do(radix.FlatCmd(&res, "HMSET", key, values)); err != nil {
		return "", err
	} else {
		return res, err
	}
	return res, nil
}

func (p *pikaResource) Lpush(ctx context.Context, key string, value interface{}) (int, error) {
	s := <-p.connpool.GetConnChan(p.label)
	prefix := p.connpool.GetPrefix(p.label)
	key = prefix + key
	var res int
	if err := s.Do(radix.FlatCmd(&res, "Lpush", key, value)); err != nil {
		return 0, err
	} else {
		return res, err
	}
	return res, nil
}

func (p *pikaResource) Rpush(ctx context.Context, key string, value interface{}) (int, error) {
	s := <-p.connpool.GetConnChan(p.label)
	prefix := p.connpool.GetPrefix(p.label)
	key = prefix + key
	var res int
	if err := s.Do(radix.FlatCmd(&res, "Rpush", key, value)); err != nil {
		return 0, err
	} else {
		return res, err
	}
	return res, nil
}

func (p *pikaResource) Sadd(ctx context.Context, key string, value interface{}) (int, error) {
	s := <-p.connpool.GetConnChan(p.label)
	prefix := p.connpool.GetPrefix(p.label)
	key = prefix + key
	var res int
	if err := s.Do(radix.FlatCmd(&res, "Sadd", key, value)); err != nil {
		return 0, err
	} else {
		return res, err
	}
	return res, nil
}

func (p *pikaResource) Srem(ctx context.Context, key string, value interface{}) (int, error) {
	s := <-p.connpool.GetConnChan(p.label)
	prefix := p.connpool.GetPrefix(p.label)
	key = prefix + key
	var res int
	if err := s.Do(radix.FlatCmd(&res, "Srem", key, value)); err != nil {
		return 0, err
	} else {
		return res, err
	}
	return res, nil
}

func (p *pikaResource) Exists(ctx context.Context, key string) (interface{}, error) {
	s := <-p.connpool.GetConnChan(p.label)
	prefix := p.connpool.GetPrefix(p.label)
	key = prefix + key
	var flag int
	if err := s.Do(radix.Cmd(&flag, "Exists", key)); err != nil {
		return nil, err
	} else {
		return flag, err
	}
}

func (p *pikaResource) Get(ctx context.Context, key string) (interface{}, error) {
	s := <-p.connpool.GetConnChan(p.label)
	prefix := p.connpool.GetPrefix(p.label)
	key = prefix + key
	var fooVal string
	if err := s.Do(radix.FlatCmd(&fooVal, "Get", key)); err != nil {
		return nil, err
	} else {
		return fooVal, err
	}
}

func (p *pikaResource) Keys(ctx context.Context, key string) (interface{}, error) {
	s := <-p.connpool.GetConnChan(p.label)
	prefix := p.connpool.GetPrefix(p.label)
	if key == "*" {
		return nil, errors.New("forbidden")
	}
	key = prefix + key
	var bazEls []string
	if err := s.Do(radix.FlatCmd(&bazEls, "Keys", key)); err != nil {
		return nil, err
	} else {
		return bazEls, err
	}
}

func (p *pikaResource) Hget(ctx context.Context, key string, name string) (interface{}, error) {
	s := <-p.connpool.GetConnChan(p.label)
	prefix := p.connpool.GetPrefix(p.label)
	key = prefix + key
	var fooVal string
	if err := s.Do(radix.FlatCmd(&fooVal, "Hget", key, name)); err != nil {
		return nil, err
	} else {
		return fooVal, err
	}
}

func (p *pikaResource) Ttl(ctx context.Context, key string) (interface{}, error) {
	s := <-p.connpool.GetConnChan(p.label)
	prefix := p.connpool.GetPrefix(p.label)
	key = prefix + key
	var intervltime int
	if err := s.Do(radix.FlatCmd(&intervltime, "Ttl", key)); err != nil {
		return nil, err
	} else {
		return intervltime, err
	}
}

func (p *pikaResource) Type(ctx context.Context, key string) (interface{}, error) {
	s := <-p.connpool.GetConnChan(p.label)
	prefix := p.connpool.GetPrefix(p.label)
	key = prefix + key
	var dataType interface{}
	if err := s.Do(radix.FlatCmd(&dataType, "Type", key)); err != nil {
		return nil, err
	} else {
		return dataType, err
	}
}

func (p *pikaResource) Hlen(ctx context.Context, key string) (int, error) {
	s := <-p.connpool.GetConnChan(p.label)
	prefix := p.connpool.GetPrefix(p.label)
	key = prefix + key
	var dataLen int
	if err := s.Do(radix.FlatCmd(&dataLen, "Hlen", key)); err != nil {
		return 0, err
	} else {
		return dataLen, err
	}
}

func (p *pikaResource) Hdel(ctx context.Context, key string, name interface{}) (int, error) {
	s := <-p.connpool.GetConnChan(p.label)
	prefix := p.connpool.GetPrefix(p.label)
	key = prefix + key
	var flag int
	if err := s.Do(radix.FlatCmd(&flag, "Hdel", key, name)); err != nil {
		return 0, err
	} else {
		return flag, err
	}
}

func (p *pikaResource) Hgetall(ctx context.Context, key string) (interface{}, error) {
	s := <-p.connpool.GetConnChan(p.label)
	prefix := p.connpool.GetPrefix(p.label)
	key = prefix + key
	var buzMap map[string]string
	if err := s.Do(radix.FlatCmd(&buzMap, "Hgetall", key)); err != nil {
		return nil, err
	} else {
		return buzMap, err
	}
}

func (p *pikaResource) Hincrby(ctx context.Context, key, field string, inc int64) (int64, error) {
	s := <-p.connpool.GetConnChan(p.label)
	prefix := p.connpool.GetPrefix(p.label)
	key = prefix + key
	var reply int64
	if err := s.Do(radix.FlatCmd(&reply, "HINCRBY", key, field, inc)); err != nil {
		return 0, err
	} else {
		return reply, err
	}
}

func (p *pikaResource) Del(ctx context.Context, key string) (interface{}, error) {
	s := <-p.connpool.GetConnChan(p.label)
	prefix := p.connpool.GetPrefix(p.label)
	key = prefix + key
	var flag int
	if err := s.Do(radix.FlatCmd(&flag, "del", key)); err != nil {
		return nil, err
	} else {
		return flag, err
	}
}

func (p *pikaResource) Llen(ctx context.Context, key string) (int, error) {
	s := <-p.connpool.GetConnChan(p.label)
	prefix := p.connpool.GetPrefix(p.label)
	key = prefix + key
	var dataLen int
	if err := s.Do(radix.FlatCmd(&dataLen, "Llen", key)); err != nil {
		return 0, err
	} else {
		return dataLen, err
	}
}

func (p *pikaResource) Lrange(ctx context.Context, key string, start int, stop int) (interface{}, error) {
	s := <-p.connpool.GetConnChan(p.label)
	prefix := p.connpool.GetPrefix(p.label)
	key = prefix + key
	var listContent []string
	if err := s.Do(radix.FlatCmd(&listContent, "Lrange", key, start, stop)); err != nil {
		return nil, err
	} else {
		return listContent, err
	}
}

func (p *pikaResource) Ltrim(ctx context.Context, key string, start int, stop int) (interface{}, error) {
	s := <-p.connpool.GetConnChan(p.label)
	prefix := p.connpool.GetPrefix(p.label)
	key = prefix + key
	var listContent interface{}
	if err := s.Do(radix.FlatCmd(&listContent, "Ltrim", key, start, stop)); err != nil {
		return nil, err
	} else {
		return listContent, err
	}
}

func (p *pikaResource) Lpop(ctx context.Context, key string) (interface{}, error) {
	s := <-p.connpool.GetConnChan(p.label)
	prefix := p.connpool.GetPrefix(p.label)
	key = prefix + key
	var listContent string
	if err := s.Do(radix.FlatCmd(&listContent, "Lpop", key)); err != nil {
		return nil, err
	} else {
		return listContent, err
	}
}

func (p *pikaResource) Rpop(ctx context.Context, key string) (interface{}, error) {
	s := <-p.connpool.GetConnChan(p.label)
	prefix := p.connpool.GetPrefix(p.label)
	key = prefix + key
	var listContent string
	if err := s.Do(radix.FlatCmd(&listContent, "Rpop", key)); err != nil {
		return nil, err
	} else {
		return listContent, err
	}
}

func (p *pikaResource) Scard(ctx context.Context, key string) (int, error) {
	s := <-p.connpool.GetConnChan(p.label)
	prefix := p.connpool.GetPrefix(p.label)
	key = prefix + key
	var setLen int
	if err := s.Do(radix.FlatCmd(&setLen, "Scard", key)); err != nil {
		return 0, err
	} else {
		return setLen, err
	}
}

func (p *pikaResource) Smembers(ctx context.Context, key string) (interface{}, error) {
	s := <-p.connpool.GetConnChan(p.label)
	prefix := p.connpool.GetPrefix(p.label)
	key = prefix + key
	var setContent []string
	if err := s.Do(radix.FlatCmd(&setContent, "Smembers", key)); err != nil {
		return nil, err
	} else {
		return setContent, err
	}
}

func (p *pikaResource) Sismember(ctx context.Context, key string, value interface{}) (int, error) {
	s := <-p.connpool.GetConnChan(p.label)
	prefix := p.connpool.GetPrefix(p.label)
	key = prefix + key
	var flag int
	if err := s.Do(radix.FlatCmd(&flag, "Sismember", key, value)); err != nil {
		return 0, err
	} else {
		return flag, err
	}
}

func (p *pikaResource) Zrank(ctx context.Context, key string, member interface{}) (int, error) {
	s := <-p.connpool.GetConnChan(p.label)
	prefix := p.connpool.GetPrefix(p.label)
	key = prefix + key
	var rank int
	if err := s.Do(radix.FlatCmd(&rank, "Zrank", key, member)); err != nil {
		return 0, err
	} else {
		return rank, err
	}
}

func (p *pikaResource) Zscore(ctx context.Context, key string, member interface{}) (string, error) {
	s := <-p.connpool.GetConnChan(p.label)
	prefix := p.connpool.GetPrefix(p.label)
	key = prefix + key
	var score string
	if err := s.Do(radix.FlatCmd(&score, "Zscore", key, member)); err != nil {
		return "", err
	} else {
		return score, err
	}
}

func (p *pikaResource) Zrange(ctx context.Context, key string, start int, stop int, withscores bool) (interface{}, error) {
	s := <-p.connpool.GetConnChan(p.label)
	prefix := p.connpool.GetPrefix(p.label)
	key = prefix + key
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

func (p *pikaResource) Zrevrange(ctx context.Context, key string, start int, stop int, withscores bool) (interface{}, error) {
	s := <-p.connpool.GetConnChan(p.label)
	prefix := p.connpool.GetPrefix(p.label)
	key = prefix + key
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

func (p *pikaResource) Zrangebyscore(ctx context.Context, key string, start int, stop int, withscores bool, limitOffset, limitCount int) (interface{}, error) {
	s := <-p.connpool.GetConnChan(p.label)
	prefix := p.connpool.GetPrefix(p.label)
	key = prefix + key
	var setContent []string
	if withscores {
		if err := s.Do(radix.FlatCmd(&setContent, "ZRANGEBYSCORE", key, start, stop, "WITHSCORES", "LIMIT", limitOffset, limitCount)); err != nil {
			return nil, err
		} else {
			return setContent, err
		}
	} else {
		if err := s.Do(radix.FlatCmd(&setContent, "ZRANGEBYSCORE", key, start, stop, "LIMIT", limitOffset, limitCount)); err != nil {
			return nil, err
		} else {
			return setContent, err
		}
	}

}

func (p *pikaResource) Zrevrangebyscore(ctx context.Context, key string, start int, stop int, withscores bool, limitOffset, limitCount int) (interface{}, error) {
	s := <-p.connpool.GetConnChan(p.label)
	prefix := p.connpool.GetPrefix(p.label)
	key = prefix + key
	var setContent []string
	if withscores {
		if err := s.Do(radix.FlatCmd(&setContent, "ZREVRANGEBYSCORE", key, start, stop, "WITHSCORES", "LIMIT", limitOffset, limitCount)); err != nil {
			return nil, err
		} else {
			return setContent, err
		}
	} else {
		if err := s.Do(radix.FlatCmd(&setContent, "ZREVRANGEBYSCORE", key, start, stop, "LIMIT", limitOffset, limitCount)); err != nil {
			return nil, err
		} else {
			return setContent, err
		}
	}

}

func (p *pikaResource) ZINCRBY(ctx context.Context, key string, increment int, member interface{}) (string, error) {
	s := <-p.connpool.GetConnChan(p.label)
	prefix := p.connpool.GetPrefix(p.label)
	key = prefix + key
	var score string
	if err := s.Do(radix.FlatCmd(&score, "Zincrby", key, increment, member)); err != nil {
		return "", err
	} else {
		return score, err
	}
}

func (p *pikaResource) Zadd(ctx context.Context, key string, score interface{}, member interface{}) (int, error) {
	s := <-p.connpool.GetConnChan(p.label)
	prefix := p.connpool.GetPrefix(p.label)
	key = prefix + key
	var flag int
	if err := s.Do(radix.FlatCmd(&flag, "Zadd", key, score, member)); err != nil {
		return 0, err
	} else {
		return flag, err
	}
}

func (p *pikaResource) Zrem(ctx context.Context, key string, member ...interface{}) (int, error) {
	s := <-p.connpool.GetConnChan(p.label)
	prefix := p.connpool.GetPrefix(p.label)
	key = prefix + key
	var flag int
	if err := s.Do(radix.FlatCmd(&flag, "Zrem", key, member...)); err != nil {
		return 0, err
	} else {
		return flag, err
	}
}
