package zgopg

import (
	"git.zhugefang.com/gocore/zgo/config"
	"github.com/go-pg/pg"
	"sync"
)

//PgResourcer 给service使用
type PgResourcer interface {
	GetDBChan(label string) chan *pg.DB
}

//内部结构体
type PgResource struct {
	label    string
	mu       sync.RWMutex
	connpool ConnPooler
}

func NewPgResourcer(label string) PgResourcer {
	return &PgResource{
		label:    label,
		connpool: NewConnPool(label), //使用connpool
	}
}

func InitPgResource(hsm map[string][]*config.ConnDetail) {
	InitConnPool(hsm)
}

//GetDBChan 返回存放连接的chan
func (n *PgResource) GetDBChan(label string) chan *pg.DB {
	return n.connpool.GetDBChan(label)
}
