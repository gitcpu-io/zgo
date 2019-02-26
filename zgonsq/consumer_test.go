package zgonsq

import (
	"fmt"
	"git.zhugefang.com/gocore/zgo.git/config"
	"testing"
	"time"
)

func TestConsumer(t *testing.T) {
	hsm := make(map[string][]config.ConnDetail)
	cd_bj := config.ConnDetail{
		C:        "北京主库-----nsq",
		Host:     "localhost",
		Port:     4150,
		ConnSize: 50,
		PoolSize: 20000,
	}
	cd_sh := config.ConnDetail{
		C:        "上海主库-----nsq",
		Host:     "localhost",
		Port:     4150,
		ConnSize: 50,
		PoolSize: 20000,
	}
	var s1 []config.ConnDetail
	var s2 []config.ConnDetail
	s1 = append(s1, cd_bj)
	s2 = append(s2, cd_sh)
	hsm = map[string][]config.ConnDetail{
		label_bj: s1,
		label_sh: s2,
	}
	InitNsq(hsm) //测试时表示使用nsq，在zgo_start中使用一次

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
