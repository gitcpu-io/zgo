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

type EtcConfig struct {
	Key       string
	Endpoints []string
}

func (ec *EtcConfig) InitConfigByEtcd() ([]*mvccpb.KeyValue, chan map[string][]*ConnDetail, chan *CacheConfig, chan *CacheConfig, chan map[string][]*ConnDetail, chan map[string]*CacheConfig) {
	c, err := ec.CreateClient() //创建etcd client
	if err != nil || c == nil {
		panic(errors.New("连接ETCD失败:" + err.Error()))
		return nil, nil, nil, nil, nil, nil
	}
	client = c

	//从etcd中取出key并赋值
	response, err := client.KV.Get(context.TODO(), ec.Key, clientv3.WithPrefix())
	if err != nil {
		panic(errors.New("Etcd can't connected ..."))
	}

	if len(response.Kvs) == 0 {
		fmt.Println("Etcd配置中心没有该项目信息,资源组件不可用,请联系zgo engine Admin管理平台添加...")
	}

	watchStartRev := response.Header.Revision + 1

	ch1, ch2, ch3, ch4, ch5 := ec.Watcher(ec.Key, watchStartRev)

	return response.Kvs, ch1, ch2, ch3, ch4, ch5
}

func (ec *EtcConfig) Watcher(prefixKey string, watchStartRev int64) (chan map[string][]*ConnDetail, chan *CacheConfig, chan *CacheConfig, chan map[string][]*ConnDetail, chan map[string]*CacheConfig) {

	outConnCh := make(chan map[string][]*ConnDetail) //put 资源chan
	outCacheCh := make(chan *CacheConfig)            //put cache chan
	outLogCh := make(chan *CacheConfig)              //put log chan

	outDelConnCh := make(chan map[string][]*ConnDetail) //delete 资源chan
	outDelCacheCh := make(chan map[string]*CacheConfig) //delete cache or log chan这里共用一个

	go func() {
		watcher := clientv3.NewWatcher(client)

		wch := watcher.Watch(context.TODO(), prefixKey, clientv3.WithPrefix(), clientv3.WithPrevKV(), clientv3.WithRev(watchStartRev))

		for r := range wch {
			for _, v := range r.Events {
				key := string(v.Kv.Key)
				labelType := strings.Split(key, "/")[3]

				switch v.Type {
				case clientv3.EventTypeDelete: //监听到删除操作
					val := v.PrevKv.Value
					err := ec.watchDelete(labelType, key, val, outDelCacheCh, outDelConnCh)
					if err != nil {
						fmt.Println("反序列化当前值失败", key)
						break
					}
				case clientv3.EventTypePut: //监听到put操作

					val := v.Kv.Value

					if v.IsCreate() { //如果监听到是第一次创建资源组件

						err := ec.watchFirstPut(labelType, key, val, outCacheCh, outLogCh, outConnCh)
						if err != nil {
							fmt.Println("反序列化当前值失败", key)
							break
						}

					} else { //如果监听到是第二次以上更新资源组件

						preVal := v.PrevKv.Value //上一次的值

						err := ec.watchSecondPut(labelType, key, val, preVal, outCacheCh, outLogCh, outConnCh)
						if err != nil {
							fmt.Println("反序列化当前值失败", key)
							break
						}
					}

				}

			}
		}
	}()
	return outConnCh, outCacheCh, outLogCh, outDelConnCh, outDelCacheCh
}

// watchDelete 监听到删除操作时
func (ec *EtcConfig) watchDelete(labelType string, key string, b []byte, outDelCacheCh chan map[string]*CacheConfig, outDelConnCh chan map[string][]*ConnDetail) error {
	var cm CacheConfig
	var m []ConnDetail

	if labelType == EtcTKCache || labelType == EtcTKLog {
		//删除cache或log
		err := zgoutils.Utils.Unmarshal(b, &cm)
		if err != nil {
			return err
		}

		hsm := make(map[string]*CacheConfig)

		hsm[key] = &cm

		outDelCacheCh <- hsm

	} else {
		//删除中间件redis/mysql/nsq/pika/mongo/kafka等
		err := zgoutils.Utils.Unmarshal(b, &m)
		if err != nil {
			return err
		}

		hsm := ec.changeStructToPtr(m, key)

		outDelConnCh <- hsm
	}
	return nil
}

// watchFirstPut 第一次监听到put操作，应用于资源组件第一次创建时
func (ec *EtcConfig) watchFirstPut(labelType string, key string, b []byte, outCacheCh chan *CacheConfig, outLogCh chan *CacheConfig, outConnCh chan map[string][]*ConnDetail) error {
	var cm CacheConfig
	var m []ConnDetail

	if labelType == EtcTKCache || labelType == EtcTKLog {
		err := zgoutils.Utils.Unmarshal(b, &cm)
		if err != nil {
			return err
		}
		switch labelType {
		case EtcTKCache:
			outCacheCh <- &cm //直接把结构体指针放入chan

		case EtcTKLog:
			outLogCh <- &cm //直接把结构体指针放入chan

		}
	} else {

		err := zgoutils.Utils.Unmarshal(b, &m)
		if err != nil {
			return err
		}
		hsm := ec.changeStructToPtr(m, key)

		outConnCh <- hsm

	}
	return nil

}

// watchSecondPut 第二次监听到key的put变化，用上一次的value到当前的比较，不同时就用当前的值
func (ec *EtcConfig) watchSecondPut(labelType string, key string, val []byte, preVal []byte, outCacheCh chan *CacheConfig, outLogCh chan *CacheConfig, outConnCh chan map[string][]*ConnDetail) error {
	var cm CacheConfig
	var m []ConnDetail
	var preCm CacheConfig

	if labelType == EtcTKCache || labelType == EtcTKLog { //如果监听到cache有变化
		err := zgoutils.Utils.Unmarshal(val, &cm)
		if err != nil {
			return err
		}
		err = zgoutils.Utils.Unmarshal(preVal, &preCm)
		if err != nil {
			return err
		}

		if reflect.DeepEqual(cm, preCm) != true { //如果有变化

			switch labelType {

			case EtcTKCache:
				outCacheCh <- &cm

			case EtcTKLog:
				outLogCh <- &cm
			}
		}

	} else {

		err := zgoutils.Utils.Unmarshal(val, &m)
		if err != nil {
			return err
		}
		err = zgoutils.Utils.Unmarshal(preVal, &preCm)
		if err != nil {
			return err
		}
		if reflect.DeepEqual(m, preCm) != true { //如果有变化使用当前的m

			hsm := ec.changeStructToPtr(m, key)

			outConnCh <- hsm
		}

	}
	return nil
}

// changeStructToPtr 转化[]为map，且struct为ptr
func (ec *EtcConfig) changeStructToPtr(m []ConnDetail, key string) map[string][]*ConnDetail {
	var tmp []*ConnDetail
	for _, vv := range m {
		pvv := vv
		tmp = append(tmp, &pvv)
	}
	hsm := make(map[string][]*ConnDetail)
	hsm[key] = tmp

	return hsm
}

func (ec *EtcConfig) CreateClient() (*clientv3.Client, error) {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   ec.Endpoints,
		DialTimeout: 20 * time.Second,
	})
	return cli, err
}

//nsq
//"[{\"c\":\"北京主库2-----etcd nsq\",\"host\":\"localhost\",\"port\":4150,\"connSize\":5,\"poolSize\":550},{\"c\":\"北京主库1-----etcd nsq\",\"host\":\"localhost\",\"port\":4150,\"connSize\":5,\"poolSize\":500}]"

//log
//"{\"c\": \"日志存储3gg000443333g34433333444445599\",\"start\": 1,\"dbType\": \"file\",\"label\":\"/tmp\"}"
