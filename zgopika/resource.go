package zgopika

import (
	"context"
	"errors"
	"git.zhugefang.com/gocore/zgo.git/config"
	"github.com/mediocregopher/radix"
	"sync"
)

//NsqResourcer 给service使用
type PikaResourcer interface {
	GetConnChan(label string) chan *radix.Pool
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
	//hget
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

func (p *pikaResource) Set(ctx context.Context, key string, value string, time int) (interface{}, error) {
	s := <-p.connpool.GetConnChan(p.label)
	prefix := p.connpool.GetPrefix(p.label)
	key = prefix + key
	pipline := radix.Pipeline(
		radix.FlatCmd(nil, "SET", key, value),
		radix.FlatCmd(nil, "Expire", key, time),
	)
	return nil, s.Do(pipline)
}

func (p *pikaResource) Expire(ctx context.Context, key string, time int) (interface{}, error) {
	s := <-p.connpool.GetConnChan(p.label)
	prefix := p.connpool.GetPrefix(p.label)
	key = prefix + key
	return nil, s.Do(radix.FlatCmd(nil, "Expire", key, time))
}

func (p *pikaResource) Hset(ctx context.Context, key string, name string, value string) (interface{}, error) {
	s := <-p.connpool.GetConnChan(p.label)
	prefix := p.connpool.GetPrefix(p.label)
	key = prefix + key
	return nil, s.Do(radix.FlatCmd(nil, "Hset", key, name, value))
}

func (p *pikaResource) Lpush(ctx context.Context, key string, value string) (interface{}, error) {
	s := <-p.connpool.GetConnChan(p.label)
	prefix := p.connpool.GetPrefix(p.label)
	key = prefix + key
	return nil, s.Do(radix.FlatCmd(nil, "Lpush", key, value))
}

func (p *pikaResource) Rpush(ctx context.Context, key string, value string) (interface{}, error) {
	s := <-p.connpool.GetConnChan(p.label)
	prefix := p.connpool.GetPrefix(p.label)
	key = prefix + key
	return nil, s.Do(radix.FlatCmd(nil, "Rpush", key, value))
}

func (p *pikaResource) Sadd(ctx context.Context, key string, value string) (interface{}, error) {
	s := <-p.connpool.GetConnChan(p.label)
	prefix := p.connpool.GetPrefix(p.label)
	key = prefix + key
	return nil, s.Do(radix.FlatCmd(nil, "Sadd", key, value))
}

func (p *pikaResource) Srem(ctx context.Context, key string, value string) (interface{}, error) {
	s := <-p.connpool.GetConnChan(p.label)
	prefix := p.connpool.GetPrefix(p.label)
	key = prefix + key
	return nil, s.Do(radix.FlatCmd(nil, "Srem", key, value))
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

func (p *pikaResource) Hlen(ctx context.Context, key string) (interface{}, error) {
	s := <-p.connpool.GetConnChan(p.label)
	prefix := p.connpool.GetPrefix(p.label)
	key = prefix + key
	var dataLen int
	if err := s.Do(radix.FlatCmd(&dataLen, "Hlen", key)); err != nil {
		return nil, err
	} else {
		return dataLen, err
	}
}

func (p *pikaResource) Hdel(ctx context.Context, key string, name string) (interface{}, error) {
	s := <-p.connpool.GetConnChan(p.label)
	prefix := p.connpool.GetPrefix(p.label)
	key = prefix + key
	var flag int
	if err := s.Do(radix.FlatCmd(&flag, "Hlen", key)); err != nil {
		return nil, err
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

func (p *pikaResource) Llen(ctx context.Context, key string) (interface{}, error) {
	s := <-p.connpool.GetConnChan(p.label)
	prefix := p.connpool.GetPrefix(p.label)
	key = prefix + key
	var dataLen int
	if err := s.Do(radix.FlatCmd(&dataLen, "Llen", key)); err != nil {
		return nil, err
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

func (p *pikaResource) Lpop(ctx context.Context, key string) (interface{}, error) {
	s := <-p.connpool.GetConnChan(p.label)
	prefix := p.connpool.GetPrefix(p.label)
	key = prefix + key
	var listContent int
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
	var listContent int
	if err := s.Do(radix.FlatCmd(&listContent, "Lpop", key)); err != nil {
		return nil, err
	} else {
		return listContent, err
	}
}

func (p *pikaResource) Scard(ctx context.Context, key string) (interface{}, error) {
	s := <-p.connpool.GetConnChan(p.label)
	prefix := p.connpool.GetPrefix(p.label)
	key = prefix + key
	var setLen int
	if err := s.Do(radix.FlatCmd(&setLen, "Scard", key)); err != nil {
		return nil, err
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

func (p *pikaResource) Sismember(ctx context.Context, key string, value string) (interface{}, error) {
	s := <-p.connpool.GetConnChan(p.label)
	prefix := p.connpool.GetPrefix(p.label)
	key = prefix + key
	var flag int
	if err := s.Do(radix.FlatCmd(&flag, "Sismember", key)); err != nil {
		return nil, err
	} else {
		return flag, err
	}
}
