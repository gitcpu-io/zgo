package zgoneo4j

import (
  "github.com/gitcpu-io/zgo/config"
  "github.com/neo4j/neo4j-go-driver/neo4j"
)

//Neo4jResourcer 给service使用
type Neo4jResourcer interface {
  GetConnChan(label string) chan neo4j.Session
}

//内部结构体
type Neo4jResource struct {
  label    string
  //mu       sync.RWMutex
  connpool ConnPooler
}

func NewNeo4jResourcer(label string) Neo4jResourcer {
  return &Neo4jResource{
    label:    label,
    connpool: NewConnPool(label), //使用connpool
  }
}

func InitNeo4jResource(hsm map[string][]*config.ConnDetail) {
  InitConnPool(hsm)
}

//GetConnChan 返回存放连接的chan
func (n *Neo4jResource) GetConnChan(label string) chan neo4j.Session {
  return n.connpool.GetConnChan(label)
}
