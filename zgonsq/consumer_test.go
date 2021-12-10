package zgonsq

import (
	"fmt"
	"github.com/gitcpu-io/zgo/config"
	"testing"
	"time"
)

func TestConsumer(t *testing.T) {
	hsm := make(map[string][]*config.ConnDetail)
	cd_bj := config.ConnDetail{
		C:        "北京从库1-----nsq",
		Host:     "localhost",
		Port:     4150,
		ConnSize: 5,
		PoolSize: 246,
	}
	//cd_bj2 := config.ConnDetail{
	//	C:        "北京从库2-----nsq",
	//	Host:     "localhost",
	//	Port:     4150,
	//	ConnSize: 10,
	//	PoolSize: 135,
	//}
	cd_sh := config.ConnDetail{
		C:        "上海主库-----nsq",
		Host:     "localhost",
		Port:     4152,
		ConnSize: 50,
		PoolSize: 20000,
	}
	var s1 []*config.ConnDetail
	var s2 []*config.ConnDetail
	s1 = append(s1, &cd_bj)
	//s1 = append(s1, &cd_bj, &cd_bj2)
	s2 = append(s2, &cd_sh)
	hsm = map[string][]*config.ConnDetail{
		label_bj: s1,
		label_sh: s2,
	}
	InitNsq(hsm) //测试时表示使用nsq，在origin中使用一次

	time.Sleep(2 * time.Second)

	labelBj, err := GetNsq(label_bj)

	labelSh, err := GetNsq(label_sh)
	if err != nil {
		panic(err)
	}
	c := chat{
		Topic:   "origin",
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
	//time.Sleep(3 * time.Second)

}
