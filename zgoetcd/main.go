package main

import (
	"encoding/json"
	"git.zhugefang.com/gocore/zgo"
	"git.zhugefang.com/gocore/zgo/config"
	"time"

	"fmt"

	"context"

	"go.etcd.io/etcd/clientv3"
)

type A struct {
	ABC string `json:"abc"`
	DEF string `json:"def"`
}

func main() {
	cli, err := CreateClient()
	if err != nil {
		return
	}
	defer cli.Close()

	//------------

	err = zgo.Engine(&zgo.Options{
		Env: "local",
		Project: "zgo_start",
	})
	if err != nil {
		panic(err)
	}


	for _, v := range config.Nsq {
		go func(v config.LabelDetail) {
			k := v.Key
			value := v.Values
			key := "zgo/project/zgo_start/nsq/" + k
			val, _ := json.Marshal(value)

			res, err := cli.KV.Put(context.TODO(), key, string(val),clientv3.WithPrevKV())
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println(res.PrevKv)
		}(v)
	}
	time.Sleep(5 * time.Second)

	//for _, v := range config.Mongo {
	//	k := v.Key
	//	value := v.Values
	//	key := "zgo/project/zgo_start/mongo/" + k
	//	val, _ := json.Marshal(value)
	//	go cli.KV.Put(context.TODO(), key, string(val))
	//}
	//
	//for _, v := range config.Es {
	//	k := v.Key
	//	value := v.Values
	//	key := "zgo/project/zgo_start/es/" + k
	//	val, _ := json.Marshal(value)
	//	go cli.KV.Put(context.TODO(), key, string(val))
	//}
	//for _, v := range config.Mysql {
	//	k := v.Key
	//	value := v.Values
	//	key := "zgo/project/zgo_start/mysql/" + k
	//	val, _ := json.Marshal(value)
	//	go cli.KV.Put(context.TODO(), key, string(val))
	//}
	//for _, v := range config.Etcd {
	//	k := v.Key
	//	value := v.Values
	//	key := "zgo/project/zgo_start/etcd/" + k
	//	val, _ := json.Marshal(value)
	//	go cli.KV.Put(context.TODO(), key, string(val))
	//}
	//for _, v := range config.Kafka {
	//	k := v.Key
	//	value := v.Values
	//	key := "zgo/project/zgo_start/kafka/" + k
	//	val, _ := json.Marshal(value)
	//	go cli.KV.Put(context.TODO(), key, string(val))
	//}
	//
	//for _, v := range config.Redis {
	//	k := v.Key
	//	value := v.Values
	//	key := "zgo/project/zgo_start/redis/" + k
	//	val, _ := json.Marshal(value)
	//	go cli.KV.Put(context.TODO(), key, string(val))
	//}
	//for _, v := range config.Pika {
	//	k := v.Key
	//	value := v.Values
	//	key := "zgo/project/zgo_start/pika/" + k
	//	val, _ := json.Marshal(value)
	//	go cli.KV.Put(context.TODO(), key, string(val))
	//}
	//
	//key := "zgo/cache/one"
	//val, _ := json.Marshal(config.Cache)
	//cli.KV.Put(context.TODO(), key, string(val))

	fmt.Println("all config to etcd done")

}

func CreateClient() (*clientv3.Client, error) {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints: []string{
			"0.0.0.0:2381",
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
