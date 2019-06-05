package zgoetcd

import (
	"git.zhugefang.com/gocore/zgo/config"
	"go.etcd.io/etcd/clientv3"
	"sync"
)

//EtcdResourcer 给service使用
type EtcdResourcer interface {
	GetConnChan(label string) chan *clientv3.Client
}

//内部结构体
type EtcdResource struct {
	label    string
	mu       sync.RWMutex
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
