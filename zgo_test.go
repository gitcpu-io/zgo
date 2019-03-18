package zgo

import (
	"context"
	"fmt"
	"testing"
	"time"
)

//mysql struct
type MysqlUser struct {
	Host int    `json:"host"`
	User string `json:"user"`
}

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
			//****************************************test log
			//Log.Error("start engine for test")

			//****************************************test mysql default user table
			//n, err := Mysql.New("mysql_sell_1")
			//if err != nil {
			//	fmt.Println("======error=====",err)
			//}
			//args := make(map[string]interface{})
			//args["table"] = "user"
			//args["query"] = " user = ? "
			//args["args"] = []interface{}{string("root")}
			//args["limit"] = 30
			//args["offset"] = 0
			//args["order"] = " host desc "
			//obj := make([]MysqlUser,100)
			//args["obj"] = &obj
			//n.List(context.TODO(), args)
			//fmt.Println(obj)

			//****************************************test nsq
			nq, err := Nsq.New()
			if err != nil {
				fmt.Println("---error", err)
			}
			nq.Producer(context.TODO(), "zgo_start", []byte("zgo engine is niubility from nsq"))

			//****************************************test kafka
			//kq, err := Kafka.New("kafka_label_bj")
			//if err != nil {
			//	fmt.Println("---error", err)
			//}
			//kq.Producer(context.TODO(), "zgo_start", []byte("zgo engine is niubility from kafka"))

		}
	}
}
