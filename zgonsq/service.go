// zgonsq是对消息中间件NSQ的封装，提供新建连接，生产数据，消费数据接口
package zgonsq

import (
	"context"
	"github.com/nsqio/go-nsq"
	"github.com/gitcpu-io/zgo/comm"
	"github.com/gitcpu-io/zgo/config"
	"sync"
)

var (
	currentLabels = make(map[string][]*config.ConnDetail) //用于存放label与具体Host:port的map
	muLabel       sync.RWMutex                            //用于并发读写上面的map
)

//Nsq 对外
type Nsqer interface {
	/*
	 label: 可选，如果使用者，用了2个或多个label时，需要调用这个函数，传入label
	*/
	// New 生产一条消息到Nsq
	New(label ...string) (*zgonsq, error)

	/*
	 label: 可选，如果使用者，用了2个或多个label时，需要调用这个函数，传入label
	*/
	// GetConnChan 获取原生的生产者client，返回一个chan，使用者需要接收 <- chan
	GetConnChan(label ...string) (chan *nsq.Producer, error)

	/*
	 ctx:是上下文参数，由使用者传入，用于控制这个函数是否超时
	 topic:string
	 body: 是一个[]byte
	*/
	// Producer 生产一条消息到Nsq
	Producer(ctx context.Context, topic string, body []byte) (chan uint8, error)

	/*
	 ctx:是上下文参数，由使用者传入，用于控制这个函数是否超时
	 topic:string
	 body: 是一个slice of []byte
	*/
	// ProducerMulti 生产多条消息到Nsq
	ProducerMulti(ctx context.Context, topic string, body [][]byte) (chan uint8, error)

	/*
	 topic:string
	 channel:string
	 fn:是自定义的回调函数，由使用者传入，经过处理后，使用者可以使用它返回的数据 func(message NsqMessage) error
	*/
	// Consumer 消费者使用
	Consumer(topic, channel string, fn NsqHandlerFunc)
}

// Nsq用于对zgo.Nsq这个全局变量赋值
func Nsq(label string) Nsqer {
	return &zgonsq{
		res: NewNsqResourcer(label),
	}
}

// zgonsq实现了Nsq的接口
type zgonsq struct {
	res NsqResourcer //使用resource另外的一个接口
}

// InitNsq 初始化连接nsq，用于使用者zgo.engine时，zgo init
func InitNsq(hsmIn map[string][]*config.ConnDetail, label ...string) chan *zgonsq {
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

	InitNsqResource(hsm)

	//自动为变量初始化对象
	initLabel := ""
	for k, _ := range hsm {
		if k != "" {
			initLabel = k
			break
		}
	}
	out := make(chan *zgonsq)
	go func() {

		in, err := GetNsq(initLabel)
		if err != nil {
			panic(err)
		}
		out <- in
		close(out)
	}()

	return out

}

// GetNsq zgo内部获取一个连接nsq
func GetNsq(label ...string) (*zgonsq, error) {
	l, err := comm.GetCurrentLabel(label, muLabel, currentLabels)
	if err != nil {
		return nil, err
	}
	return &zgonsq{
		res: NewNsqResourcer(l),
	}, nil
}

// NewNsq获取一个Nsq生产者的client，用于发送数据
func (n *zgonsq) New(label ...string) (*zgonsq, error) {
	return GetNsq(label...)
}

//GetConnChan 供用户使用原生连接的chan
func (n *zgonsq) GetConnChan(label ...string) (chan *nsq.Producer, error) {
	l, err := comm.GetCurrentLabel(label, muLabel, currentLabels)
	if err != nil {
		return nil, err
	}
	return n.res.GetConnChan(l), nil
}

// Producer 生产一条消息
func (n *zgonsq) Producer(ctx context.Context, topic string, body []byte) (chan uint8, error) {
	return n.res.Producer(ctx, topic, body)
}

// ProducerMulti 生产多条消息
func (n *zgonsq) ProducerMulti(ctx context.Context, topic string, body [][]byte) (chan uint8, error) {
	return n.res.ProducerMulti(ctx, topic, body)
}

// Consumer 消费者
func (n *zgonsq) Consumer(topic, channel string, fn NsqHandlerFunc) {
	go n.res.Consumer(topic, channel, fn)
}
