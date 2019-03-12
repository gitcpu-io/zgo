package zgo

import (
	"context"
	"git.zhugefang.com/gocore/zgo/zgonsq"
	"testing"
	"time"
)

func TestEngine(t *testing.T) {

	err := Engine(&Options{
		Env:     "dev",
		Project: "zgo_start",
		Nsq:     []string{ //测试etcd配置动态改库成功
			//"nsq_label_bj",
			//"nsq_label_sh",
		},
		Mongo: []string{ //测试etcd配置动态改库成功
			//"mongo_label_bj",
			//"mongo_label_sh",
		},
		Kafka: []string{ //测试etcd配置动态改库成功
			//"kafka_label_bj",
		},
		Redis: []string{ //测试etcd配置动态改库成功
			//"redis_label_bj",
		},
		Pika: []string{ //测试etcd配置动态改库成功
			//"pika_label_rw",
		},

		Mysql: []string{
			//"mysql_sell_1",
		},
	})

	if err != nil {
		panic(err)
	}

	for {
		select {
		case <-time.Tick(time.Duration(5) * time.Second):
			n := zgonsq.Nsq("nsq_label_bj")
			n.Producer(context.TODO(), "zgo_start", []byte("999999999999"))
			Log.Error("start engine for test")
		}
	}
}
