package config

import (
	"context"
	"errors"
	"fmt"
	"git.zhugefang.com/gocore/zgo/zgoutils"
	"go.etcd.io/etcd/clientv3"
	"go.etcd.io/etcd/mvcc/mvccpb"
	"reflect"
	"strings"
	"time"
)

/*
@Time : 2019-03-04 15:09
@Author : rubinus.chu
@File : etcd
@project: zgo
*/

var client *clientv3.Client

func init() {
	cli, err := CreateClient()
	if err != nil {
		return
	}
	client = cli
	return
}

func InitConfigByEtcd() chan map[string][]*ConnDetail {
	prefixKey := "zgo"

	//从etcd中取出key并赋值
	response, err := client.KV.Get(context.TODO(), prefixKey, clientv3.WithPrefix())
	if err != nil {
		panic(errors.New("Etcd can't connected ..."))
	}

	ch := make(chan *mvccpb.KeyValue)
	go func() {
		for v := range ch {
			var tmp LabelDetail
			var lade []LabelDetail
			mk := string(v.Key)
			smk := strings.Split(mk, "/")
			b := v.Value
			var m []ConnDetail
			err := zgoutils.Utils.Unmarshal(b, &m)
			if err != nil {
				fmt.Println("反序列化当前值失败", mk)
			}
			tmp.Key = smk[2]
			tmp.Values = m

			lade = append(lade, tmp)

			key := smk[1]
			switch key {
			case mysqlT:
				//init mysql again
				Mysql = lade
			case mongoT:
				//init mongo again
				Mongo = lade
			case redisT:
				//init redis again
				Redis = lade
			case pikaT:
				//init pika again
				Pika = lade
			case nsqT:
				//init nsq again
				Nsq = lade
			case kafkaT:
				//init kafka again
				Kafka = lade
			case esT:
				//init es again
				Es = lade
			case etcdT:
				//init etcd again
			}
			fmt.Println(Nsq)

		}

	}()
	for _, v := range response.Kvs {
		ch <- v
	}
	close(ch)
	//开始监控
	return Watcher(client, prefixKey)
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
							for _, vv := range m {
								tmp = append(tmp, &vv)
							}
							//k := strings.Split(key, "/")[2]
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
		Endpoints: []string{
			"0.0.0.0:2381",
		},
		DialTimeout: 10 * time.Second,
	})
	return cli, err
}
