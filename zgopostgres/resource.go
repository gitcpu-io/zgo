package zgopostgres

import (
	"git.zhugefang.com/gocore/zgo/config"
	"github.com/go-pg/pg"
	"sync"
)

//PostgresResourcer 给service使用
type PostgresResourcer interface {
	GetConnChan(label string) chan *pg.DB
}

//内部结构体
type PostgresResource struct {
	label    string
	mu       sync.RWMutex
	connpool ConnPooler
}

func NewPostgresResourcer(label string) PostgresResourcer {
	return &PostgresResource{
		label:    label,
		connpool: NewConnPool(label), //使用connpool
	}
}

func InitPostgresResource(hsm map[string][]*config.ConnDetail) {
	InitConnPool(hsm)
}

//GetConnChan 返回存放连接的chan
func (n *PostgresResource) GetConnChan(label string) chan *pg.DB {
	return n.connpool.GetConnChan(label)
}
