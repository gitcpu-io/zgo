package config

import (
	"context"
	"encoding/json"
	"etcd/clientv3"
	"fmt"
	"time"
)

/*
@Time : 2019-03-04 15:09
@Author : rubinus.chu
@File : etcd
@project: zgo
*/

var client *clientv3.Client

func Init() *clientv3.Client {
	cli, err := CreateClient()
	if err != nil {
		return nil
	}
	client = cli
	return client
}

func InitConfigByEtcd() {
	cli, err := CreateClient()
	if err != nil {
		return
	}
	key := "zgo/nsq/nsq_label_bj"

	gr, err := cli.KV.Get(context.TODO(), key, clientv3.WithPrevKV())

	vv := gr.Kvs[0].Value
	cnd := []ConnDetail{}
	json.Unmarshal(vv, &cnd)

	fmt.Println(cnd)
	fmt.Println(cnd)

	Watcher(cli, key)
}

func Watcher(client *clientv3.Client, key string) {

	watcher := clientv3.NewWatcher(client)
	wch := watcher.Watch(context.TODO(), key, clientv3.WithPrevKV())
	go func() {
		for {
			select {
			case r := <-wch:
				fmt.Println("----watch key do something---", key)
				fmt.Printf("%+v %s", r, "\n")
				for k, v := range r.Events {
					fmt.Println("key==", string(v.Kv.Key))
					b := v.Kv.Value
					m := []ConnDetail{}
					json.Unmarshal(b, &m)
					fmt.Println(m, k)
					fmt.Println("有没有变化", v.IsModify())
				}
			}
		}
	}()
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
