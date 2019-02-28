package zgonsq

import (
	"fmt"
	"golang-bootstrap/lib/log"
)

type chat struct {
	Topic   string
	Channel string
	Nsq     Nsqer
}

func (c *chat) Consumer() {
	go c.Nsq.Consumer(c.Topic, c.Channel, 2, c.Deal)
}

//处理消息
func (c *chat) Deal(msg NsqMessage) error {

	fmt.Println("接收到NSQ", msg.NSQDAddress, ",message:", string(msg.Body))
	log.Info("接收到NSQ", msg.NSQDAddress, ",message:", string(msg.Body))

	return nil
}
