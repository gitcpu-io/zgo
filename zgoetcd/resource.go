package zgoetcd

import (
  "github.com/gitcpu-io/zgo/config"
  "go.etcd.io/etcd/client/v3"
)

//EtcdResourcer 给service使用
type EtcdResourcer interface {
  GetConnChan(label string) chan *clientv3.Client
}

//内部结构体
type EtcdResource struct {
  label string
  //mu       sync.RWMutex
  connpool ConnPooler
}

func NewEtcdResourcer(label string) EtcdResourcer {
  return &EtcdResource{
    label:    label,
    connpool: NewConnPool(label), //使用connpool
  }
}

func InitEtcdResource(hsm map[string][]*config.ConnDetail) {
  InitConnPool(hsm)
}

//GetConnChan 返回存放连接的chan
func (n *EtcdResource) GetConnChan(label string) chan *clientv3.Client {
  return n.connpool.GetConnChan(label)
}
