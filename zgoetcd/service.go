// zgoetcd是对消息中间件Etcd的封装，提供新建连接，生产数据，消费数据接口
package zgoetcd

import (
	"git.zhugefang.com/gocore/zgo/comm"
	"git.zhugefang.com/gocore/zgo/config"
	"go.etcd.io/etcd/clientv3"
	"sync"
)

var (
	currentLabels = make(map[string][]*config.ConnDetail) //用于存放label与具体Host:port的map
	muLabel       sync.RWMutex                            //用于并发读写上面的map
)

//Etcd 对外
type Etcder interface {
	/*
	 label: 可选，如果使用者，用了2个或多个label时，需要调用这个函数，传入label
	*/
	// New 生产一条消息到Etcd
	New(label ...string) (*zgoetcd, error)

	/*
	 label: 可选，如果使用者，用了2个或多个label时，需要调用这个函数，传入label
	*/
	// GetConnChan 获取原生的生产者client，返回一个chan，使用者需要接收 <- chan
	GetConnChan(label ...string) (chan *clientv3.Client, error)
}

// Etcd用于对zgo.Etcd这个全局变量赋值
func Etcd(label string) Etcder {
	return &zgoetcd{
		res: NewEtcdResourcer(label),
	}
}

// zgoetcd实现了Etcd的接口
type zgoetcd struct {
	res EtcdResourcer //使用resource另外的一个接口
}

// InitEtcd 初始化连接postgres，用于使用者zgo.engine时，zgo init
func InitEtcd(hsmIn map[string][]*config.ConnDetail, label ...string) chan *zgoetcd {
	muLabel.Lock()
	defer muLabel.Unlock()

	var hsm map[string][]*config.ConnDetail

	if len(label) > 0 && len(currentLabels) > 0 { //此时是destory操作,传入的hsm是nil
		//fmt.Println("--destory--前",currentLabels)
		for _, v := range label {
			delete(currentLabels, v)
		}
		hsm = currentLabels
		//fmt.Println("--destory--后",currentLabels)

	} else { //这是第一次创建操作或etcd中变更时init again操作
		hsm = hsmIn
		//currentLabels = hsm	//this operation is error
		for k, v := range hsm { //so big bug can't set hsm to currentLabels，must be for, may be have old label
			currentLabels[k] = v
		}
	}

	if len(hsm) == 0 {
		return nil
	}

	InitEtcdResource(hsm)

	//自动为变量初始化对象
	initLabel := ""
	for k, _ := range hsm {
		if k != "" {
			initLabel = k
			break
		}
	}
	out := make(chan *zgoetcd)
	go func() {

		in, err := GetEtcd(initLabel)
		if err != nil {
			panic(err)
		}
		out <- in
		close(out)
	}()

	return out

}

// GetEtcd zgo内部获取一个连接postgres
func GetEtcd(label ...string) (*zgoetcd, error) {
	l, err := comm.GetCurrentLabel(label, muLabel, currentLabels)
	if err != nil {
		return nil, err
	}
	return &zgoetcd{
		res: NewEtcdResourcer(l),
	}, nil
}

// NewEtcd获取一个Etcd生产者的client，用于发送数据
func (n *zgoetcd) New(label ...string) (*zgoetcd, error) {
	return GetEtcd(label...)
}

//GetConnChan 供用户使用原生连接的chan
func (n *zgoetcd) GetConnChan(label ...string) (chan *clientv3.Client, error) {
	l, err := comm.GetCurrentLabel(label, muLabel, currentLabels)
	if err != nil {
		return nil, err
	}
	return n.res.GetConnChan(l), nil
}
