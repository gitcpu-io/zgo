// zgokafka是对消息中间件Kafka的封装，提供新建连接，生产数据，消费数据接口
package zgokafka

import (
  "context"
  "github.com/Shopify/sarama"
  "github.com/bsm/sarama-cluster"
  "github.com/gitcpu-io/zgo/comm"
  "github.com/gitcpu-io/zgo/config"
  "sync"
)

var (
  currentLabels = make(map[string][]*config.ConnDetail)
  muLabel       sync.RWMutex
)

//Kafka 对外
type Kafkaer interface {
  New(label ...string) (*zgokafka, error)
  GetConnChan(label ...string) (chan *sarama.AsyncProducer, error)
  Producer(ctx context.Context, topic string, body []byte) (chan uint8, error)
  ProducerMulti(ctx context.Context, topic string, body [][]byte) (chan uint8, error)
  Consumer(topic, groupId string) (*cluster.Consumer, error)
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
func InitKafka(hsmIn map[string][]*config.ConnDetail, label ...string) chan *zgokafka {
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

  InitKafkaResource(hsm)

  //自动为变量初始化对象
  initLabel := ""
  for k, _ := range hsm {
    if k != "" {
      initLabel = k
      break
    }
  }
  out := make(chan *zgokafka)
  go func() {
    in, err := GetKafka(initLabel)
    if err != nil {
      out <- nil
    }
    out <- in
    close(out)
  }()
  return out

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

func (n *zgokafka) New(label ...string) (*zgokafka, error) {
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

func (n *zgokafka) Consumer(topic, groupId string) (*cluster.Consumer, error) {
  return n.res.Consumer(topic, groupId)
}
