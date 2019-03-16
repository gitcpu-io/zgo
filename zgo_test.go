package zgo

import (
	"testing"
	"time"
)

func TestEngine(t *testing.T) {

	err := Engine(&Options{
		Env:     "dev",
		Project: "zgo_start",

		//如果是在本地开发可以对下面的组件开启使用(local.json)，如果是线上，不需要填写，走的配置是etcd
		Kafka: []string{
			"kafka_label_bj",
			"kafka_label_sh",
		},
		Nsq: []string{
			"nsq_label_bj",
			"nsq_label_sh",
		},
		Pika: []string{
			//"pika_label_rw",
			//"pika_label_r",
		},
		//Redis: []string{
		//	"redis_label_bj",
		//"redis_label_sh",
		//},
		//Es: []string{
		//	"label_new",
		//	"label_rent",
		//	"label_sell",
		//},
		//Mysql: []string{
		//	"mysql_sell_1",
		//	"mysql_sell_2",
		//},
		//Mongo: []string{
		//	"mongo_label_bj",
		//	"mongo_label_sh",
		//},
	})

	if err != nil {
		panic(err)
	}

	for {
		select {
		case <-time.Tick(time.Duration(3) * time.Second):
			Log.Error("start engine for test")

			//n := zgokafka.Kafka("kafka_label_bj")
			//n.Producer(context.TODO(), "zgo_start", []byte("dsfsdfsdfsfsfsdfsdfss"))
		}
	}
}
