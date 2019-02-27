package zgokafka

import (
	"context"
	"git.zhugefang.com/gocore/zgo.git/comm"
	"git.zhugefang.com/gocore/zgo.git/config"
	"github.com/Shopify/sarama"
	"sync"
)

var (
	currentLabels = make(map[string][]*config.ConnDetail)
	muLabel       sync.RWMutex
)

//Kafka 对外
type Kafkaer interface {
	NewKafka(label ...string) (*zgokafka, error)
	GetConnChan(label ...string) (chan *sarama.AsyncProducer, error)
	Producer(ctx context.Context, topic string, body []byte) (chan uint8, error)
	ProducerMulti(ctx context.Context, topic string, body [][]byte) (chan uint8, error)
	Consumer(topic, channel string, mode int, fn KafkaHandlerFunc)
}

func Kafka(label string) Kafkaer {
	return &zgokafka{
		res: NewKafkaResourcer(label),
	}
}

//zgokafka实现了Kafka的接口
type zgokafka struct {
	res KafkaResourcer //使用resource另外的一个接口
}

//InitKafka 初始化连接kafka
func InitKafka(hsm map[string][]*config.ConnDetail) {
	muLabel.Lock()
	defer muLabel.Unlock()

	currentLabels = hsm
	InitKafkaResource(hsm)
}

//GetKafka zgo内部获取一个连接kafka
func GetKafka(label ...string) (*zgokafka, error) {
	l, err := comm.GetCurrentLabel(label, muLabel, currentLabels)
	if err != nil {
		return nil, err
	}
	return &zgokafka{
		res: NewKafkaResourcer(l),
	}, nil
}

func (n *zgokafka) NewKafka(label ...string) (*zgokafka, error) {
	return GetKafka(label...)
}

//GetConnChan 供用户使用原生连接的chan
func (n *zgokafka) GetConnChan(label ...string) (chan *sarama.AsyncProducer, error) {
	l, err := comm.GetCurrentLabel(label, muLabel, currentLabels)
	if err != nil {
		return nil, err
	}
	return n.res.GetConnChan(l), nil
}

func (n *zgokafka) Producer(ctx context.Context, topic string, body []byte) (chan uint8, error) {
	return n.res.Producer(ctx, topic, body)
}

func (n *zgokafka) ProducerMulti(ctx context.Context, topic string, body [][]byte) (chan uint8, error) {
	return n.res.ProducerMulti(ctx, topic, body)
}

//func (n *zgokafka) Consumer(topic, channel string, mode int, fn KafkaHandlerFunc) {
//	go n.res.Consumer(topic, channel, mode, fn)
//}
