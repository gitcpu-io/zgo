package zgo

import (
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
		Env:     "local",
		Project: "zgo_start",
		//Project: "1552641690",

		//如果是在本地开发可以对下面的组件开启使用(local.json)，如果是线上，不需要填写，走的配置是etcd
		Kafka: []string{
			//"kafka_label_bj",
			//"kafka_label_sh",
		},
		Nsq: []string{
			//"nsq_label_bj",
			//"nsq_label_sh",
		},
		Pika: []string{
			//"pika_label_rw",
			//"pika_label_r",
		},
		//Postgres: []string{
		//	"postgres_label_sh",
		//},
		Neo4j: []string{
			"neo4j_label",
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
			Log.Error("start engine for test")

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

			//****************************************test postgres
			//pgch, err := PG.GetDBChan("postgres_label_sh")
			//if err != nil {
			//	fmt.Println("---error", err)
			//}
			//db := <- pgch
			//fmt.Println("zgo engine is niubility from postgres",db)

			//****************************************test neo4j
			neo4jch, err := Neo4j.GetDBChan("neo4j_label")
			if err != nil {
				fmt.Println("---error", err)
			}
			db := <-neo4jch
			fmt.Println("zgo engine is niubility from neo4j", db)

			//****************************************test nsq
			//nq, err := Nsq.New()
			//if err != nil {
			//	fmt.Println("---error", err)
			//}
			//nq.Producer(context.TODO(), "zgo_start", []byte("zgo engine is niubility from nsq"))

			//****************************************test kafka
			//kq, err := Kafka.New("kafka_label_bj")
			//if err != nil {
			//	fmt.Println("---error", err)
			//}
			//kq.Producer(context.TODO(), "zgo_start", []byte("zgo engine is niubility from kafka"))

		}
	}
}
