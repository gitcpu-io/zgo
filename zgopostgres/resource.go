package zgopostgres

import (
  "github.com/gitcpu-io/zgo/config"
  "github.com/go-pg/pg"
  "github.com/go-pg/pg/orm"
  "sync"
)

//PostgresResourcer 给service使用
type PostgresResourcer interface {
  GetConnChan(label string) chan *pg.DB
  Scan(values ...interface{}) orm.ColumnScanner
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

func (n *PostgresResource) Scan(values ...interface{}) orm.ColumnScanner {
  return pg.Scan(values...)
}
