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
          "c": "北京主库1-----nsq",
          "host": "localhost",
          "port": 4150,
          "connSize": 50,
          "poolSize": 250
        },
        {
          "c": "北京主库2-----nsq",
          "host": "localhost",
          "port": 4150,
          "connSize": 50,
          "poolSize": 369
        }
      ]`

	key := "zgo/nsq/nsq_label_bj"
	cli.KV.Put(context.TODO(), key, v)

	//gr, err := cli.KV.Get(context.TODO(), key, clientv3.WithPrevKV())
	//
	//vv := gr.Kvs[0].Value
	//cnd := []config.ConnDetail{}
	//json.Unmarshal(vv, &cnd)
	//
	//fmt.Println(cnd)

	//watcher := clientv3.NewWatcher(cli)
	//wch := watcher.Watch(context.TODO(), key, clientv3.WithPrevKV())
	//go func() {
	//	for {
	//		select {
	//		case r := <-wch:
	//			fmt.Println("----watch---")
	//			fmt.Printf("%+v %s", r, "\n")
	//		}
	//	}
	//}()

	//sr, err := cli.Status(context.TODO(), "10.20.80.132:2379")
	//if err != nil {
	//	fmt.Println(err)
	//}
	//fmt.Printf("%+v\n", sr)

	//------------
	//cli.KV.Put(context.TODO(), "abc", "abcdef")
	//gr, err := cli.KV.Get(context.TODO(), "abc", clientv3.WithPrevKV())
	//fmt.Println("第一次get abc:", gr.Kvs, err, gr.Count)
	//
	////------------
	//
	//le, err := cli.Lease.Grant(context.TODO(), 3)
	////fmt.Println(le.TTL)
	//s := "/abc/def"
	//_, err = cli.KV.Put(context.TODO(), s, "傻逼12", clientv3.WithLease(le.ID))
	////fmt.Println(pr, "test etcdDemo ...")
	////gr, err := cli.KV.Get(context.TODO(), "test")
	////fmt.Println(gr.Kvs)
	//
	////------------
	//
	//ctx, cancelFunc := context.WithCancel(context.TODO())
	//defer cancelFunc()
	//defer cli.Lease.Revoke(context.TODO(), le.ID)
	//
	//lech, err := cli.Lease.KeepAlive(ctx, le.ID)
	//if err != nil {
	//	fmt.Println(err, "--keepalive")
	//}
	//
	////------------
	//
	//kvc := clientv3.NewKV(cli)
	//gr, err = kvc.Get(context.TODO(), s, clientv3.WithPrevKV())
	//fmt.Println("第一次get:", gr.Kvs, err, gr.Count)
	//
	//Watcher(cli, s, lech, le)
	////for x := range wch {
	////	for _, v := range x.Events {
	////		fmt.Printf("%+v\n", v)
	////	}
	////}
	//test(kvc, le, err)
	//
	//time.Sleep(3500 * time.Millisecond)

	//fmt.Println("/n======")
	//gr, err = kvc.Get(context.TODO(), "/abc/def", clientv3.WithPrevKV())
	//fmt.Println(gr.Kvs, err, gr.Count)

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
