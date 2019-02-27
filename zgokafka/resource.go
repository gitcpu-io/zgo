package zgokafka

import (
	"context"
	"errors"
	"fmt"
	"git.zhugefang.com/gocore/zgo.git/config"
	"github.com/Shopify/sarama"
	"github.com/bsm/sarama-cluster"
	"log"
	"os"
	"os/signal"
	"sync"
)

type (
	KafkaMessage     = *cluster.Consumer
	KafkaHandlerFunc func(consumer KafkaMessage, signals chan os.Signal) error //定义一个对外的handlerfunc
)

//
//type kafkaHandler interface {
//	handleMessage(m *kafka.Message) error
//	HandleMessage(m *kafka.Message) error
//}
//
////实现kafka原生的方法
//func (h KafkaHandlerFunc) handleMessage(m KafkaMessage) error {
//	return h(m)
//}
//
////实现kafka原生的方法
//func (h KafkaHandlerFunc) HandleMessage(m *kafka.Message) error {
//	return h(m)
//}

//KafkaResourcer 给service使用
type KafkaResourcer interface {
	GetConnChan(label string) chan *sarama.AsyncProducer
	Producer(ctx context.Context, topic string, body []byte) (chan uint8, error)
	ProducerMulti(ctx context.Context, topic string, body [][]byte) (chan uint8, error)
	//Consumer(topic, channel string, mode int, fn KafkaHandlerFunc) //自定义的func
}

//内部结构体
type kafkaResource struct {
	label    string
	mu       sync.RWMutex
	connpool ConnPooler
}

func NewKafkaResourcer(label string) KafkaResourcer {
	return &kafkaResource{
		label:    label,
		connpool: NewConnPool(label), //使用connpool
	}
}

func InitKafkaResource(hsm map[string][]*config.ConnDetail) {
	InitConnPool(hsm)
}

//GetConnChan 返回存放连接的chan
func (n *kafkaResource) GetConnChan(label string) chan *sarama.AsyncProducer {
	return n.connpool.GetConnChan(label)
}

//Producer 生产者函数
func (n *kafkaResource) Producer(ctx context.Context, topic string, body []byte) (chan uint8, error) {
	out := make(chan uint8, 1)
	if len(body) == 0 { //不能发布空串，否则会导致error
		out <- 0
		return out, errors.New("message is empty")
	}
	doneChan := <-n.connpool.GetConnChan(n.label)
	err := n.PublishAsync(topic, body, doneChan) // 发布消息
	if err != nil {
		out <- 0
		return out, nil
	}
	//go func() {
	//	r := <-doneChan //一定要消费掉这个，要不然会丢消息
	//	if r == nil {
	//		//fmt.Println(topic,"--发送到NSQ失败--",err)
	//		out <- 0
	//	} else {
	//		//fmt.Println(topic,"==发送到NSQ成功==", string(body),err)
	//		out <- 1
	//	}
	//}()
	return out, nil
}

//Producer 生产者函数
func (n *kafkaResource) ProducerMulti(ctx context.Context, topic string, body [][]byte) (chan uint8, error) {
	out := make(chan uint8, 1)
	if len(body) == 0 { //不能发布空串，否则会导致error
		out <- 0
		return out, errors.New("message is empty")
	}
	doneChan := <-n.connpool.GetConnChan(n.label)
	err := n.PublishAsync(topic, []byte(""), doneChan) // 发布消息
	if err != nil {
		out <- 0
		return out, nil
	}
	//go func() {
	//	r := <-doneChan //一定要消费掉这个，要不然会丢消息
	//	if r == nil {
	//		//fmt.Println(topic,"--发送到NSQ失败--",err)
	//		out <- 0
	//	} else {
	//		//fmt.Println(topic,"==发送到NSQ成功==", string(body),err)
	//		out <- 1
	//	}
	//}()
	return out, nil
}

//初始化消费者
func (n *kafkaResource) Consumer(brokers, topics []string, groupId string, callback KafkaHandlerFunc) {
	config := cluster.NewConfig()
	//config.Consumer.Return.Errors = true
	//config.Group.Return.Notifications = true

	config.Group.Mode = cluster.ConsumerModePartitions

	consumer, err := cluster.NewConsumer(brokers, groupId, topics, config)
	if err != nil {
		log.Println(brokers, err.Error())
		//panic(err)
	}
	defer consumer.Close()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	// consume errors
	//go func() {
	//	for err := range consumer.Errors() {
	//		log.Printf("Error: %s\n", err.Error())
	//	}
	//}()

	// consume notifications
	//go func() {
	//	for ntf := range consumer.Notifications() {
	//		log.Printf("Kafka connection: %+v\n", ntf)
	//	}
	//}()

	callback(consumer, signals)

}

// asyncProducer 异步生产者
// 并发量大时，必须采用这种方式
func (n *kafkaResource) PublishAsync(topics string, value []byte, ptr *sarama.AsyncProducer) error {
	p := *ptr
	//必须有这个匿名函数内容
	go func(p sarama.AsyncProducer) {
		errors := p.Errors()
		success := p.Successes()
		for {
			select {
			case err := <-errors:
				if err != nil {
					fmt.Println("asyn send=", err, topics, value)
				}

			case <-success:
				//fmt.Printf("Partition:%d\nOffset:%d\n%s\n%s",s.Partition,s.Offset,s.Value,s.Timestamp)
				fmt.Fprintln(os.Stdout, "\n", value, "==done success==", topics)
			}
		}
	}(p)

	msg := &sarama.ProducerMessage{
		Topic: topics,
		Value: sarama.ByteEncoder(value),
	}
	p.Input() <- msg
	return nil

}
