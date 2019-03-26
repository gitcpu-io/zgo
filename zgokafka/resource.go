package zgokafka

import (
	"context"
	"errors"
	"fmt"
	"git.zhugefang.com/gocore/zgo/config"
	"github.com/Shopify/sarama"
	"github.com/bsm/sarama-cluster"
	"sync"
	"time"
)

//KafkaResourcer 给service使用
type KafkaResourcer interface {
	GetConnChan(label string) chan *sarama.AsyncProducer
	Producer(ctx context.Context, topic string, body []byte) (chan uint8, error)
	ProducerMulti(ctx context.Context, topic string, body [][]byte) (chan uint8, error)
	Consumer(topic, groupId string) (*cluster.Consumer, error) //自定义的func
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
	producer := n.connpool.GetConnChan(n.label)
	//if len(producer) == 0 {
	//	out <- 0
	//	return out, errors.New("conn is invalid")
	//}
	n.PublishAsync(topic, body, <-producer, out) // 发布消息

	return out, nil
}

//Producer 生产者函数
func (n *kafkaResource) ProducerMulti(ctx context.Context, topic string, body [][]byte) (chan uint8, error) {
	out := make(chan uint8, 1)
	if len(body) == 0 { //不能发布空串，否则会导致error
		out <- 0
		return out, errors.New("message is empty")
	}
	producer := n.connpool.GetConnChan(n.label)
	//if len(producer) == 0 {
	//	out <- 0
	//	return out, errors.New("conn is invalid")
	//}

	n.PublishMultiAsync(topic, body, <-producer, out) // 发布消息

	return out, nil
}

//初始化消费者
func (n *kafkaResource) Consumer(topic, groupId string) (*cluster.Consumer, error) {
	config := cluster.NewConfig()
	config.Consumer.Return.Errors = true
	config.Group.Return.Notifications = true
	config.Net.KeepAlive = 30 * time.Minute
	config.Net.MaxOpenRequests = 20000

	//config.Group.Mode = cluster.ConsumerModePartitions

	n.mu.RLock()
	defer n.mu.RUnlock()
	var brokers []string
	if addr, ok := currentLabels[n.label]; ok {
		for k, v := range addr {
			if k == 0 && v.Host != "" && v.Port != 0 {
				address := fmt.Sprintf("%s:%d", v.Host, v.Port)
				brokers = append(brokers, address)
				break
			}
		}
	}

	consumer, err := cluster.NewConsumer(brokers, groupId, []string{topic}, config)
	if err != nil {
		panic(err)
	}

	// consume errors
	go func() {
		for err := range consumer.Errors() {
			fmt.Printf("Error: %s\n", err.Error())
		}
	}()

	// consume notifications
	go func() {
		for _ = range consumer.Notifications() {
			//fmt.Printf("Kafka connection: %+v\n", ntf)
		}
	}()

	return consumer, err

	//if err != nil {
	//	log.Println(brokers, err.Error())
	//panic(err)
	//}
	//defer consumer.Close()

	//signals := make(chan os.Signal, 1)
	//signal.Notify(signals, os.Interrupt)

	//callback(consumer, signals)

}

// asyncProducer 异步生产者
// 并发量大时，必须采用这种方式
func (n *kafkaResource) PublishAsync(topics string, value []byte, ptr *sarama.AsyncProducer, in chan uint8) {
	p := *ptr
	//必须有这个匿名函数内容
	//go func(p sarama.AsyncProducer, in chan uint8) {
	//	errors := p.Errors()
	//	success := p.Successes()
	//	for {
	//		select {
	//		case err := <-errors:
	//			if err != nil {
	//				fmt.Println(err)
	//			}
	//			in <- 0
	//		case <-success:
	//			//fmt.Printf("Partition:%d\nOffset:%d\n%s\n%s",s.Partition,s.Offset,s.Value,s.Timestamp)
	//			//fmt.Fprintln(os.Stdout, "\n", string(value), "==done success==", topics)
	//			in <- 1
	//		}
	//	}
	//}(p, in)

	msg := &sarama.ProducerMessage{
		Topic: topics,
		Value: sarama.ByteEncoder(value),
	}
	p.Input() <- msg
	in <- 1
}

func (n *kafkaResource) PublishMultiAsync(topics string, value [][]byte, ptr *sarama.AsyncProducer, in chan uint8) {
	p := *ptr
	//必须有这个匿名函数内容
	go func(p sarama.AsyncProducer, in chan uint8) {
		errors := p.Errors()
		success := p.Successes()
		for {
			select {
			case err := <-errors:
				if err != nil {
					fmt.Println(err, topics)
				}
				in <- 0
			case <-success:
				//fmt.Printf("Partition:%d\nOffset:%d\n%s\n%s",s.Partition,s.Offset,s.Value,s.Timestamp)
				//fmt.Fprintln(os.Stdout, "\n", string(value), "==done success==", topics)
				in <- 1
			}
		}
	}(p, in)

	for _, v := range value {
		msg := &sarama.ProducerMessage{
			Topic: topics,
			Value: sarama.ByteEncoder(v),
		}
		p.Input() <- msg
	}

}
