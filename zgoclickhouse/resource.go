package zgoclickhouse

import (
  "database/sql"
  "github.com/gitcpu-io/zgo/config"
  "sync"
)

//ClickHouseResourcer 给service使用
type ClickHouseResourcer interface {
  GetConnChan(label string) chan *sql.DB
}

//内部结构体
type ClickHouseResource struct {
  label    string
  mu       sync.RWMutex
  connpool ConnPooler
}

func NewClickHouseResourcer(label string) ClickHouseResourcer {
  return &ClickHouseResource{
    label:    label,
    connpool: NewConnPool(label), //使用connpool
  }
}

func InitClickHouseResource(hsm map[string][]*config.ConnDetail) {
  InitConnPool(hsm)
}

//GetConnChan 返回存放连接的chan
func (n *ClickHouseResource) GetConnChan(label string) chan *sql.DB {
  return n.connpool.GetConnChan(label)
}
