package zgonsq

import (
	"context"
	"errors"
	"fmt"
	"git.zhugefang.com/gocore/zgo/config"
	"github.com/nsqio/go-nsq"
	"sync"
	"time"
)

type (
	NsqMessage     = *nsq.Message
	NsqHandlerFunc func(message NsqMessage) error //定义一个对外的handlerfunc
)

type nsqHandler interface {
	handleMessage(m *nsq.Message) error
	HandleMessage(m *nsq.Message) error
}

//实现nsq原生的方法
func (h NsqHandlerFunc) handleMessage(m NsqMessage) error {
	return h(m)
}

//实现nsq原生的方法
func (h NsqHandlerFunc) HandleMessage(m *nsq.Message) error {
	return h(m)
}

//NsqResourcer 给service使用
type NsqResourcer interface {
	GetConnChan(label string) chan *nsq.Producer
	Producer(ctx context.Context, topic string, body []byte) (chan uint8, error)
	ProducerMulti(ctx context.Context, topic string, body [][]byte) (chan uint8, error)
	Consumer(topic, channel string, fn NsqHandlerFunc) //自定义的func
}

//内部结构体
type nsqResource struct {
	label    string
	mu       sync.RWMutex
	connpool ConnPooler
}

func NewNsqResourcer(label string) NsqResourcer {
	return &nsqResource{
		label:    label,
		connpool: NewConnPool(label), //使用connpool
	}
}

func InitNsqResource(hsm map[string][]*config.ConnDetail) {
	InitConnPool(hsm)
}

//GetConnChan 返回存放连接的chan
func (n *nsqResource) GetConnChan(label string) chan *nsq.Producer {
	return n.connpool.GetConnChan(label)
}

//Producer 生产者函数
func (n *nsqResource) Producer(ctx context.Context, topic string, body []byte) (chan uint8, error) {

	out := make(chan uint8, 1)
	if len(body) == 0 { //不能发布空串，否则会导致error
		out <- 0
		return out, errors.New("message is empty")
	}

	producer := <-n.connpool.GetConnChan(n.label)

	if producer == nil {
		out <- 0
		return out, errors.New("conn is nil")
	}

	//for {
	//	if err := producer.Ping(); err != nil {
	//		producer = <-n.connpool.GetConnChan(n.label)
	//		fmt.Println("------producer is nil")
	//	}else{
	//		break
	//	}
	//	time.Sleep(10 * time.Millisecond)
	//}

	doneChan := make(chan *nsq.ProducerTransaction)
	err := producer.PublishAsync(topic, body, doneChan) // 发布消息

	fmt.Println(err, "---PublishAsync----", producer)

	if err != nil {
		out <- 0
		return out, nil
	}
	go func() {
		r := <-doneChan //一定要消费掉这个，要不然会丢消息
		if r == nil {
			fmt.Println(topic, "--发送到NSQ失败--", err)
			out <- 0
		} else {
			fmt.Println(topic, "==发送到NSQ成功==", string(body), err)
			out <- 1
		}
	}()
	return out, nil
}

//Producer 生产者函数
func (n *nsqResource) ProducerMulti(ctx context.Context, topic string, body [][]byte) (chan uint8, error) {
	out := make(chan uint8, 1)
	if len(body) == 0 { //不能发布空串，否则会导致error
		out <- 0
		return out, errors.New("message is empty")
	}
	producer := <-n.connpool.GetConnChan(n.label)
	doneChan := make(chan *nsq.ProducerTransaction)
	err := producer.MultiPublishAsync(topic, body, doneChan) // 发布消息
	if err != nil {
		out <- 0
		return out, nil
	}
	go func() {
		r := <-doneChan //一定要消费掉这个，要不然会丢消息
		if r == nil {
			//fmt.Println(topic,"--发送到NSQ失败--",err)
			out <- 0
		} else {
			//fmt.Println(topic,"==发送到NSQ成功==", string(body),err)
			out <- 1
		}
	}()
	return out, nil
}

//初始化消费者
func (n *nsqResource) Consumer(topic, channel string, fn NsqHandlerFunc) {
	cfg := nsq.NewConfig()
	cfg.LookupdPollInterval = time.Second //设置重连时间
	cfg.MaxInFlight = 20
	c, err := nsq.NewConsumer(topic, channel, cfg) // 新建一个消费者
	if err != nil {
		fmt.Println("===NewConsumer===", err)
		return
	}
	//c.SetLogger(nil, 0)        //屏蔽系统日志
	c.SetLogger(nil, 2)
	c.AddHandler(fn) // 添加消费者接口

	n.mu.RLock()
	defer n.mu.RUnlock()
	var address string
	if addr, ok := currentLabels[n.label]; ok {
		for k, v := range addr {
			if k == 0 && v.Host != "" && v.Port != 0 {
				address = fmt.Sprintf("%s:%d", v.Host, v.Port)
				break
			}
		}
	}

	if err := c.ConnectToNSQD(address); err != nil {
		fmt.Println("===ConnectToNSQD===", err)
	}

	//if mode == 1 {
	//	//建立NSQLookupd连接
	//	if err := c.ConnectToNSQLookupd(address); err != nil {
	//		fmt.Println("===ConnectToNSQLookupd===", err)
	//	}
	//} else if mode == 2 {
	//	if err := c.ConnectToNSQD(address); err != nil {
	//		fmt.Println("===ConnectToNSQD===", err)
	//	}
	//}

	//建立多个nsqd连接
	// if err := c.ConnectToNSQDs([]string{"127.0.0.1:4150", "127.0.0.1:4152"}); err != nil {
	//  panic(err)
	// }

	// 建立一个nsqd连接

}
