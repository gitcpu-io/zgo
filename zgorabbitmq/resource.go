package zgorabbitmq

import (
	"context"
	"errors"
	"fmt"
	"git.zhugefang.com/gocore/zgo/config"
	"github.com/streadway/amqp"
	"sync"
	"time"
)

//RabbitmqResourcer 给service使用
type RabbitmqResourcer interface {
	GetConnChan(label string) chan *amqp.Connection
	Producer(ctx context.Context, exchangeName, exchangeType, routingKey string, body []byte) (chan uint8, error)
	Consumer(exchangeName, exchangeType, routingKey, queueName string) (<-chan amqp.Delivery, error)
	ProducerByQueue(ctx context.Context, queueName string, body []byte) (chan uint8, error)
	ConsumerByQueue(queueName string) (<-chan amqp.Delivery, error)
}

//内部结构体
type rabbitmqResource struct {
	label    string
	mu       sync.RWMutex
	connpool ConnPooler
}

func NewRabbitmqResourcer(label string) RabbitmqResourcer {
	return &rabbitmqResource{
		label:    label,
		connpool: NewConnPool(label), //使用connpool
	}
}

func InitRabbitmqResource(hsm map[string][]*config.ConnDetail) {
	InitConnPool(hsm)
}

//GetConnChan 返回存放连接的chan
func (n *rabbitmqResource) GetConnChan(label string) chan *amqp.Connection {
	return n.connpool.GetConnChan(label)
}

//Producer 生产者函数
func (n *rabbitmqResource) Producer(ctx context.Context, exchangeName, exchangeType, routingKey string, body []byte) (chan uint8, error) {

	out := make(chan uint8, 1)
	if len(body) == 0 { //不能发布空串，否则会导致error
		out <- 0
		return out, errors.New("message is empty")
	}

	connChan := n.connpool.GetConnChan(n.label)

	if len(connChan) == 0 {
		out <- 0
		return out, errors.New("conn is nil")
	}

	if conn, ok := <-connChan; ok {
		c, err := conn.Channel()
		if err != nil {
			out <- 0
			return out, errors.New(fmt.Sprintf("channel.open: %s", err))
		}

		err = c.ExchangeDeclare(exchangeName, exchangeType, true, false, false, false, nil)
		if err != nil {
			out <- 0
			return out, errors.New(fmt.Sprintf("exchange.declare: %v", err))
		}

		msg := amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			Timestamp:    time.Now(),
			ContentType:  "text/plain",
			Body:         body,
		}

		err = c.Publish(exchangeName, routingKey, false, false, msg)
		if err != nil {
			out <- 0
			return out, errors.New(fmt.Sprintf("basic.publish: %v", err))
		}
		out <- 1
		return out, nil

	} else {
		out <- 0
		fmt.Println("---Publish error----no connection")
		return out, nil
	}

}

//初始化消费者
func (n *rabbitmqResource) Consumer(exchangeName, exchangeType, routingKey, queueName string) (<-chan amqp.Delivery, error) {
	out := make(<-chan amqp.Delivery, 1)
	connChan := n.connpool.GetConnChan(n.label)

	if len(connChan) == 0 {
		return out, errors.New("conn is nil")
	}
	if conn, ok := <-connChan; ok {
		c, err := conn.Channel()
		if err != nil {
			return out, errors.New(fmt.Sprintf("channel.open: %s", err))
		}
		err = c.ExchangeDeclare(exchangeName, exchangeType, true, false, false, false, nil)
		if err != nil {
			return out, errors.New(fmt.Sprintf("exchange.declare: %s", err))
		}

		_, err = c.QueueDeclare(queueName, true, false, false, false, nil)
		if err != nil {
			return out, errors.New(fmt.Sprintf("queue.declare: %v", err))
		}

		err = c.QueueBind(queueName, routingKey, exchangeName, false, nil)
		if err != nil {
			return out, errors.New(fmt.Sprintf("queue.bind: %v", err))
		}

		dev, err := c.Consume(queueName, queueName, true, false, false, false, nil)
		if err != nil {
			return out, errors.New(fmt.Sprintf("basic.consume: %v", err))
		}

		return dev, nil

	} else {
		return out, errors.New("conn is nil")
	}
}

//ProducerByQueue 生产者函数
func (n *rabbitmqResource) ProducerByQueue(ctx context.Context, queueName string, body []byte) (chan uint8, error) {

	out := make(chan uint8, 1)
	if len(body) == 0 { //不能发布空串，否则会导致error
		out <- 0
		return out, errors.New("message is empty")
	}

	connChan := n.connpool.GetConnChan(n.label)

	if len(connChan) == 0 {
		out <- 0
		return out, errors.New("conn is nil")
	}

	if conn, ok := <-connChan; ok {
		c, err := conn.Channel()
		if err != nil {
			out <- 0
			return out, errors.New(fmt.Sprintf("channel.open: %s", err))
		}

		msg := amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			Timestamp:    time.Now(),
			ContentType:  "text/plain",
			Body:         body,
		}

		err = c.Publish("", queueName, false, false, msg)
		if err != nil {
			out <- 0
			return out, errors.New(fmt.Sprintf("basic.publish: %v", err))
		}
		out <- 1
		return out, nil

	} else {
		out <- 0
		fmt.Println("---Publish error queue----no connection")
		return out, nil
	}

}

// ConsumerByQueue 初始化消费者
func (n *rabbitmqResource) ConsumerByQueue(queueName string) (<-chan amqp.Delivery, error) {
	out := make(<-chan amqp.Delivery, 1)
	connChan := n.connpool.GetConnChan(n.label)

	if len(connChan) == 0 {
		return out, errors.New("conn is nil")
	}
	if conn, ok := <-connChan; ok {
		c, err := conn.Channel()
		if err != nil {
			return out, errors.New(fmt.Sprintf("channel.open: %s", err))
		}
		_, err = c.QueueDeclare(queueName, true, false, false, false, nil)
		if err != nil {
			return out, errors.New(fmt.Sprintf("queue.declare: %v", err))
		}

		dev, err := c.Consume(queueName, queueName, true, false, false, false, nil)
		if err != nil {
			return out, errors.New(fmt.Sprintf("basic.consume: %v", err))
		}

		return dev, nil

	} else {
		return out, errors.New("conn is nil")
	}
}
