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
	"strings"
	"time"
)

var client *clientv3.Client

func InitConfigByEtcd(project string) (chan *mvccpb.KeyValue, chan map[string][]*ConnDetail, chan *CacheConfig, chan *CacheConfig) {
	c, err := CreateClient() //创建etcd client
	if err != nil {
		return nil, nil, nil, nil
	}
	client = c

	prefixKey := fmt.Sprintf("zgo/project/%s", project)
	//从etcd中取出key并赋值
	response, err := client.KV.Get(context.TODO(), prefixKey, clientv3.WithPrefix())
	if err != nil {
		panic(errors.New("Etcd can't connected ..."))
	}
	if len(response.Kvs) == 0 {
		panic(errors.New("Etcd have not u config pls checkout it ..."))
	}

	ch := make(chan *mvccpb.KeyValue, 1000)

	for _, v := range response.Kvs {
		ch <- v //返回到其它channel中
	}
	//从这个version开始监控
	watchStartRev := response.Header.Revision + 1
	ch1, ch2, ch3 := Watcher(prefixKey, watchStartRev)

	return ch, ch1, ch2, ch3
}

func Watcher(prefixKey string, watchStartRev int64) (chan map[string][]*ConnDetail, chan *CacheConfig, chan *CacheConfig) {

	outConnCh := make(chan map[string][]*ConnDetail)
	outCacheCh := make(chan *CacheConfig)
	outLogCh := make(chan *CacheConfig)

	go func() {
		watcher := clientv3.NewWatcher(client)
		wch := watcher.Watch(context.TODO(), prefixKey, clientv3.WithPrefix(), clientv3.WithPrevKV(), clientv3.WithRev(watchStartRev))

		for r := range wch {
			for _, v := range r.Events {
				switch v.Type {
				case clientv3.EventTypePut:
					key := string(v.Kv.Key)
					b := v.Kv.Value
					preb := v.PrevKv.Value //上一次的值
					keyType := strings.Split(key, "/")[3]
					fmt.Println(keyType, "-----------------")
					if keyType == "cache" || keyType == "log" { //如果监听到cache有变化
						cm := CacheConfig{}
						precm := CacheConfig{}
						err := zgoutils.Utils.Unmarshal(b, &cm)
						if err != nil {
							fmt.Println("反序列化当前值失败", key)
							continue
						}
						err = zgoutils.Utils.Unmarshal(preb, &precm)
						if err != nil {
							fmt.Println("反序列上一个值失败", key)
							continue
						}
						if reflect.DeepEqual(cm, precm) != true { //如果有变化
							switch keyType {
							case "cache":
								outCacheCh <- &cm

							case "log":
								fmt.Println("---etcd watching log change ..")
								outLogCh <- &cm
							}
						}

					} else {

						m := []ConnDetail{}
						err := zgoutils.Utils.Unmarshal(b, &m)
						if err != nil {
							fmt.Println("反序列化当前值失败", key)

							continue
						}
						prem := []ConnDetail{}
						err = zgoutils.Utils.Unmarshal(preb, &prem)
						if err != nil {
							fmt.Println("反序列上一个值失败", key)
							continue
						}
						if reflect.DeepEqual(m, prem) != true { //如果有变化
							var tmp []*ConnDetail
							for _, vv := range m {
								pvv := vv
								tmp = append(tmp, &pvv)
							}
							hsm := make(map[string][]*ConnDetail)
							hsm[key] = tmp

							outConnCh <- hsm
						}

					}
				}

			}
		}
	}()
	return outConnCh, outCacheCh, outLogCh
}

func CreateClient() (*clientv3.Client, error) {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   EtcdHosts,
		DialTimeout: 20 * time.Second,
	})
	return cli, err
}
