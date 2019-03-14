package zgo

import (
	"context"
	"git.zhugefang.com/gocore/zgo/zgokafka"
	"testing"
	"time"
)

func TestEngine(t *testing.T) {

	err := Engine(&Options{
		Env:     "local",
		Project: "zgo_start",



		//如果是在本地开发可以对下面的组件开启使用(local.json)，如果是线上，不需要填写，走的配置是etcd
		Kafka: []string{
			"kafka_label_bj",
		},
		Nsq: []string{
			"nsq_label_bj",
			"nsq_label_sh",
		},
	})

	if err != nil {
		panic(err)
	}

	for {
		select {
		case <-time.Tick(time.Duration(5) * time.Second):
			Log.Error("start engine for test")
			n := zgokafka.Kafka("kafka_label_bj")
			n.Producer(context.TODO(),"zgo_start", []byte("dsfsdfsdfsfsfsdfsdfss"))
		}
	}
}
