package zgokafka

import (
	"fmt"
	"github.com/bsm/sarama-cluster"
)

type chat struct {
	Topic   string
	GroupId string
	Kafka   Kafkaer
}

func (c *chat) Consumer() {
	consumer, _ := c.Kafka.Consumer(c.Topic, c.GroupId)
	go func() {
		for {
			select {
			case part, ok := <-consumer.Partitions():

				if !ok {
					return
				}
				// start a separate goroutine to consume messages
				go func(pc cluster.PartitionConsumer) {
					for msg := range pc.Messages() {

						fmt.Printf("==message===%d %s\n", msg.Offset, msg.Value)

					}
				}(part)
			//case <-signals:
			//	fmt.Println("activity no signals ...")
			//	return

			case msg, ok := <-consumer.Messages():
				if ok {
					fmt.Printf("==message===%d %s\n", msg.Offset, msg.Value)

				}

			}
		}
	}()

}
