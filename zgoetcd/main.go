package main

import (
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

	v := `[
        {
          "c": "北京主库1-----etcd nsq",
          "host": "localhost",
          "port": 4150,
          "connSize": 5,
          "poolSize": 25
        },
        {
          "c": "北京主库2-----etcd nsq",
          "host": "localhost",
          "port": 4150,
          "connSize": 5,
          "poolSize": 390
        }
      ]`

	key := "zgo/nsq/nsq_label_bj"
	cli.KV.Put(context.TODO(), key, v)

	v_mongo := `[
        {
          "c": "北京主库1-----etcd nsq",
          "host": "localhost",
          "port": 4150,
          "connSize": 5,
          "poolSize": 25
        },
        {
          "c": "北京主库2-----etcd nsq",
          "host": "localhost",
          "port": 4150,
          "connSize": 5,
          "poolSize": 390
        }
      ]`

	key_mongo := "zgo/mongo/mongo_label_bj"
	cli.KV.Put(context.TODO(), key_mongo, v_mongo)

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
