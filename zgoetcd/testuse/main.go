package main

import (
	"encoding/json"
	"github.com/gitcpu-io/zgo"
	"github.com/gitcpu-io/zgo/config"

	"time"

	"fmt"

	"context"

	"github.com/coreos/etcd/clientv3"
)

type A struct {
	ABC string `json:"abc"`
	DEF string `json:"def"`
}

var cli *clientv3.Client

func main() {
	c, err := CreateClient()
	if err != nil {
		return
	}
	cli = c

	//------------

	err = zgo.Engine(&zgo.Options{
		Env:     "local",
		Project: "origin",
	})
	if err != nil {
		panic(err)
	}

	for _, v := range config.Conf.Nsq {
		k := v.Key
		value := v.Values
		key := "zgo/project/origin/nsq/" + k
		val, _ := json.Marshal(value)
		res, err := cli.KV.Put(context.TODO(), key, string(val), clientv3.WithPrevKV())
		if err != nil {
			panic(err)
		}
		fmt.Println(res)
	}

	for _, v := range config.Conf.Mongo {
		k := v.Key
		value := v.Values
		key := "zgo/project/origin/mongo/" + k
		val, _ := json.Marshal(value)
		cli.KV.Put(context.TODO(), key, string(val))
	}

	for _, v := range config.Conf.Es {
		k := v.Key
		value := v.Values
		key := "zgo/project/origin/es/" + k
		val, _ := json.Marshal(value)
		cli.KV.Put(context.TODO(), key, string(val))
	}
	for _, v := range config.Conf.Mysql {
		k := v.Key
		value := v.Values
		key := "zgo/project/origin/mysql/" + k
		val, _ := json.Marshal(value)
		cli.KV.Put(context.TODO(), key, string(val))
	}
	for _, v := range config.Conf.Etcd {
		k := v.Key
		value := v.Values
		key := "zgo/project/origin/etcd/" + k
		val, _ := json.Marshal(value)
		cli.KV.Put(context.TODO(), key, string(val))
	}

	for _, v := range config.Conf.Kafka {
		k := v.Key
		value := v.Values
		key := "zgo/project/origin/kafka/" + k
		val, _ := json.Marshal(value)
		cli.KV.Put(context.TODO(), key, string(val))
	}

	for _, v := range config.Conf.Redis {
		k := v.Key
		value := v.Values
		key := "zgo/project/origin/redis/" + k
		val, _ := json.Marshal(value)
		cli.KV.Put(context.TODO(), key, string(val))
	}

	for _, v := range config.Conf.Postgres {
		k := v.Key
		value := v.Values
		key := "zgo/project/origin/postgres/" + k
		val, _ := json.Marshal(value)
		cli.KV.Put(context.TODO(), key, string(val))
	}

	for _, v := range config.Conf.Neo4j {
		k := v.Key
		value := v.Values
		key := "zgo/project/origin/neo4j/" + k
		val, _ := json.Marshal(value)
		cli.KV.Put(context.TODO(), key, string(val))
	}

	for _, v := range config.Conf.Etcd {
		k := v.Key
		value := v.Values
		key := "zgo/project/origin/etcd/" + k
		val, _ := json.Marshal(value)
		cli.KV.Put(context.TODO(), key, string(val))
	}

	//for _, v := range config.Conf.Pika {
	//	k := v.Key
	//	value := v.Values
	//	key := "zgo/project/origin/pika/" + k
	//	val, _ := json.Marshal(value)
	//	cli.KV.Put(context.TODO(), key, string(val))
	//}
	//

	key := "zgo/project/origin/cache"
	val, _ := json.Marshal(config.Conf.Cache)
	cli.KV.Put(context.TODO(), key, string(val))

	key_log := "zgo/project/origin/log"
	//val_log, _ := json.Marshal(config.Log)
	val_log := "{\"c\": \"日志存储12222\",\"start\": 1,\"dbType\": \"nsq\",\"label\":\"nsq_label_bj\"}"
	res, err := cli.KV.Put(context.TODO(), key_log, string(val_log))
	fmt.Println(res, err)

	fmt.Println("all config to etcd done")

}

func CreateClient() (*clientv3.Client, error) {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints: []string{
			"127.0.0.1:2381",
			//"123.56.173.28:2380",
		},
		DialTimeout: 10 * time.Second,
	})
	return cli, err
}

func Watcher(cli *clientv3.Client, s string, lech <-chan *clientv3.LeaseKeepAliveResponse, le *clientv3.LeaseGrantResponse) {
	watcher := clientv3.NewWatcher(cli)
	wch := watcher.Watch(context.TODO(), s, clientv3.WithPrevKV())
	go func() {
		for {
			select {
			case r := <-wch:
				fmt.Println("----watch---")
				fmt.Printf("%+v %s", r, "\n")
			case l := <-lech:
				if l == nil {
					fmt.Println("invalid keepalive", le.ID)
				} else {
					fmt.Println("-keepalive-", le.ID, l.ID, l.Revision)
				}
			}
		}
	}()
}

func test(kvc clientv3.KV, le *clientv3.LeaseGrantResponse, err error) {
	txn := kvc.Txn(context.TODO())
	txn.If(clientv3.Compare(clientv3.CreateRevision("/abc/def2"), "=", 0)).
		Then(clientv3.OpPut("/abc/def2", "10000", clientv3.WithLease(le.ID))).
		Else(clientv3.OpGet("/abc/def2"))
	response, err := txn.Commit()
	if response.Succeeded {
		fmt.Println("response", response)
	} else {
		fmt.Printf("锁占用 ")
	}
}
