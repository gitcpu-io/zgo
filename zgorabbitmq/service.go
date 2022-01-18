// zgorabbitmq是对消息中间件Rabbitmq的封装，提供新建连接，生产数据，消费数据接口
package zgorabbitmq

import (
  "context"
  "github.com/gitcpu-io/zgo/comm"
  "github.com/gitcpu-io/zgo/config"
  "github.com/streadway/amqp"
  "sync"
)

var (
  currentLabels = make(map[string][]*config.ConnDetail) //用于存放label与具体Host:port的map
  muLabel       sync.RWMutex                            //用于并发读写上面的map
)

//Rabbitmq 对外
type Rabbitmqer interface {
  /*
   label: 可选，如果使用者，用了2个或多个label时，需要调用这个函数，传入label
  */
  // New 生产一条消息到Rabbitmq
  New(label ...string) (*zgorabbitmq, error)

  /*
   label: 可选，如果使用者，用了2个或多个label时，需要调用这个函数，传入label
  */
  // GetConnChan 获取原生的生产者client，返回一个chan，使用者需要接收 <- chan
  GetConnChan(label ...string) (chan *amqp.Connection, error)

  /*
   ctx:是上下文参数，由使用者传入，用于控制这个函数是否超时
   exchangeName: 交换机名字
   exchangeType: 交换机类型，topic / direct
   routingKey: 路由key
   body: 发送的消息体 []byte
  */
  // Producer 生产一条消息到Rabbitmq
  Producer(ctx context.Context, exchangeName, exchangeType, routingKey string, body []byte) (chan uint8, error)

  /*
   exchangeName: 交换机名字
   exchangeType: 交换机类型，topic / direct
   routingKey: 路由key
   queueName: 队列的名字
  */
  // Consumer 消费者使用
  Consumer(exchangeName, exchangeType, routingKey, queueName string) (<-chan amqp.Delivery, error)

  /*
   ctx:是上下文参数，由使用者传入，用于控制这个函数是否超时
   queueName: 队列的名字
   body: 发送的消息体 []byte
  */
  // ProducerByQueue 生产一条消息到Rabbitmq，使用队列模式
  ProducerByQueue(ctx context.Context, queueName string, body []byte) (chan uint8, error)
  /*
   queueName: 队列的名字
  */
  // ConsumerByQueue 消费者使用队列模式
  ConsumerByQueue(queueName string) (<-chan amqp.Delivery, error)
}

// Rabbitmq用于对zgo.Rabbitmq这个全局变量赋值
func Rabbitmq(label string) Rabbitmqer {
  return &zgorabbitmq{
    res: NewRabbitmqResourcer(label),
  }
}

// zgorabbitmq实现了Rabbitmq的接口
type zgorabbitmq struct {
  res RabbitmqResourcer //使用resource另外的一个接口
}

// InitRabbitmq 初始化连接nsq，用于使用者zgo.engine时，zgo init
func InitRabbitmq(hsmIn map[string][]*config.ConnDetail, label ...string) chan *zgorabbitmq {
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

  InitRabbitmqResource(hsm)

  //自动为变量初始化对象
  initLabel := ""
  for k, _ := range hsm {
    if k != "" {
      initLabel = k
      break
    }
  }
  out := make(chan *zgorabbitmq)
  go func() {

    in, err := GetRabbitmq(initLabel)
    if err != nil {
      panic(err)
    }
    out <- in
    close(out)
  }()

  return out

}

// GetRabbitmq zgo内部获取一个连接nsq
func GetRabbitmq(label ...string) (*zgorabbitmq, error) {
  l, err := comm.GetCurrentLabel(label, muLabel, currentLabels)
  if err != nil {
    return nil, err
  }
  return &zgorabbitmq{
    res: NewRabbitmqResourcer(l),
  }, nil
}

// NewRabbitmq获取一个Rabbitmq生产者的client，用于发送数据
func (n *zgorabbitmq) New(label ...string) (*zgorabbitmq, error) {
  return GetRabbitmq(label...)
}

//GetConnChan 供用户使用原生连接的chan
func (n *zgorabbitmq) GetConnChan(label ...string) (chan *amqp.Connection, error) {
  l, err := comm.GetCurrentLabel(label, muLabel, currentLabels)
  if err != nil {
    return nil, err
  }
  return n.res.GetConnChan(l), nil
}

// Producer 生产一条消息
func (n *zgorabbitmq) Producer(ctx context.Context, exchangeName, exchangeType, routingKey string, body []byte) (chan uint8, error) {
  return n.res.Producer(ctx, exchangeName, exchangeType, routingKey, body)
}

// Consumer 消费者
func (n *zgorabbitmq) Consumer(exchangeName, exchangeType, routingKey, queueName string) (<-chan amqp.Delivery, error) {
  return n.res.Consumer(exchangeName, exchangeType, routingKey, queueName)
}

// ProducerByQueue 生产一条消息
func (n *zgorabbitmq) ProducerByQueue(ctx context.Context, queueName string, body []byte) (chan uint8, error) {
  return n.res.ProducerByQueue(ctx, queueName, body)
}

// ConsumerByQueue 消费者
func (n *zgorabbitmq) ConsumerByQueue(queueName string) (<-chan amqp.Delivery, error) {
  return n.res.ConsumerByQueue(queueName)
}
