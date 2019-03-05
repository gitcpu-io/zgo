package config

/*
@Time : 2019-03-04 15:09
@Author : rubinus.chu
@File : etcd
@project: zgo
*/

import (
	"context"
	"errors"
	"fmt"
	"git.zhugefang.com/gocore/zgo/zgoutils"
	"go.etcd.io/etcd/clientv3"
	"go.etcd.io/etcd/mvcc/mvccpb"
	"reflect"
	"time"
)

func InitConfigByEtcd() (chan *mvccpb.KeyValue, chan map[string][]*ConnDetail) {
	client, err := CreateClient()
	if err != nil {
		return nil, nil
	}

	prefixKey := "zgo"
	//从etcd中取出key并赋值
	response, err := client.KV.Get(context.TODO(), prefixKey, clientv3.WithPrefix())
	if err != nil {
		panic(errors.New("Etcd can't connected ..."))
	}

	ch := make(chan *mvccpb.KeyValue, 1000)

	for _, v := range response.Kvs {
		ch <- v
	}
	//开始监控
	return ch, Watcher(client, prefixKey)
}

func Watcher(client *clientv3.Client, prefixKey string) chan map[string][]*ConnDetail {
	hsm := make(map[string][]*ConnDetail)

	out := make(chan map[string][]*ConnDetail)
	watcher := clientv3.NewWatcher(client)
	wch := watcher.Watch(context.TODO(), prefixKey, clientv3.WithPrevKV(), clientv3.WithPrefix())
	go func() {
		for {
			select {
			case r := <-wch:
				for _, v := range r.Events {
					if v.Type == clientv3.EventTypePut {
						key := string(v.Kv.Key)
						b := v.Kv.Value
						m := []ConnDetail{}
						err := zgoutils.Utils.Unmarshal(b, &m)
						if err != nil {
							fmt.Println("反序列化当前值失败", key)

							return
						}
						preb := v.PrevKv.Value
						prem := []ConnDetail{}
						err = zgoutils.Utils.Unmarshal(preb, &prem)
						if err != nil {
							fmt.Println("反序列上一个值失败", key)
							return
						}

						if reflect.DeepEqual(m, prem) != true { //如果有变化
							var tmp []*ConnDetail
							for k, _ := range m {
								tmp = append(tmp, &m[k])
							}
							hsm[key] = tmp
							out <- hsm
						}
					}

				}
			}
		}
	}()
	return out
}

func CreateClient() (*clientv3.Client, error) {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   EtcdHosts,
		DialTimeout: 20 * time.Second,
	})
	return cli, err
}
