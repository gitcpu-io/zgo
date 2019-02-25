package zgonsq

import (
	"fmt"
	"testing"
	"time"
)

func TestConsumer(t *testing.T) {
	InitNsq(map[string][]string{
		label_bj: []string{
			"localhost:4150",
		},
		label_sh: []string{
			"localhost:4150",
		},
	}) //测试时表示使用nsq，在zgo_start中使用一次
	labelBj, err := GetNsq(label_bj)
	labelSh, err := GetNsq(label_sh)
	if err != nil {
		panic(err)
	}
	c := chat{
		Topic:   label_bj,
		Channel: label_bj,
		Nsq:     labelBj,
	}
	c.Consumer()

	c2 := chat{
		Topic:   label_sh,
		Channel: label_sh,
		Nsq:     labelSh,
	}
	c2.Consumer()

	for {
		select {
		case <-time.Tick(time.Duration(3 * time.Second)):
			fmt.Println("一直在消费着")
		}
	}
}
