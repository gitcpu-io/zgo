package zgo

import (
	"fmt"
	"git.zhugefang.com/gocore/zgo/config"
	"git.zhugefang.com/gocore/zgo/zgocache"
	"git.zhugefang.com/gocore/zgo/zgoes"
	"git.zhugefang.com/gocore/zgo/zgofile"
	"git.zhugefang.com/gocore/zgo/zgogrpc"
	"git.zhugefang.com/gocore/zgo/zgohttp"
	"git.zhugefang.com/gocore/zgo/zgokafka"
	"git.zhugefang.com/gocore/zgo/zgolog"
	"git.zhugefang.com/gocore/zgo/zgomongo"
	"git.zhugefang.com/gocore/zgo/zgomysql"
	"git.zhugefang.com/gocore/zgo/zgonsq"
	"git.zhugefang.com/gocore/zgo/zgopika"
	"git.zhugefang.com/gocore/zgo/zgoredis"
	"git.zhugefang.com/gocore/zgo/zgoutils"
	kafkaCluter "github.com/bsm/sarama-cluster"
	"github.com/nsqio/go-nsq"
	"go.etcd.io/etcd/mvcc/mvccpb"
	"strings"
)

type engine struct {
	opt *Options
}

//New init zgo engine
func Engine(opt *Options) error {
	engine := &engine{
		opt: opt,
	}

	resKvs, cacheCh, err := opt.init() //把zgo_start中用户定义的，映射到zgo的内存变量上
	if err != nil {
		return err
	}

	if opt.Project != "" {
		config.Conf.Project = opt.Project
	}
	if opt.Loglevel != "" {
		config.Conf.Loglevel = opt.Loglevel
	}

	//初始化GRPC
	Grpc = zgogrpc.GetGrpc()

	//init 日志存储
	LogStore = InitLogStore()

	//start 日志watch
	StartLogStoreWatcher()

	//异步start 日志消费存储协程
	go LogStore.StartQueue()

	if opt.Env == "local" {
		if len(opt.Mongo) > 0 {
			//todo someting
			hsm := engine.getConfigByOption(config.Conf.Mongo, opt.Mongo)
			//fmt.Println("--zgo.go--",config.Mongo, opt.Mongo, hsm)
			in := <-zgomongo.InitMongo(hsm)
			Mongo = in
		}

		if len(opt.Mysql) > 0 {
			//todo someting
			hsm := engine.getConfigByOption(config.Conf.Mysql, opt.Mysql)
			fmt.Println(hsm)
			// 配置信息： 城市和数据库的关系
			cdc := config.Conf.CityDbConfig
			zgomysql.InitMysqlService(hsm, cdc)
			var err error
			Mysql, err = zgomysql.MysqlService(opt.Mysql[0])
			if err != nil {
				fmt.Println(err)
			}

		}
		if len(opt.Es) > 0 {
			hsm := engine.getConfigByOption(config.Conf.Es, opt.Es)
			in := <-zgoes.InitEs(hsm)
			Es = in
		}
		if len(opt.Redis) > 0 {
			//todo someting
			hsm := engine.getConfigByOption(config.Conf.Redis, opt.Redis)
			//fmt.Println(hsm)
			in := <-zgoredis.InitRedis(hsm)
			Redis = in
		}
		if len(opt.Pika) > 0 {
			//todo someting
			hsm := engine.getConfigByOption(config.Conf.Pika, opt.Pika)
			//fmt.Println(hsm)
			in := <-zgopika.InitPika(hsm)
			Pika = in
		}
		if len(opt.Nsq) > 0 { //>0表示用户要求使用nsq
			hsm := engine.getConfigByOption(config.Conf.Nsq, opt.Nsq)
			//fmt.Println("===zgo.go==", hsm)
			//return nil
			in := <-zgonsq.InitNsq(hsm)
			Nsq = in
		}
		if len(opt.Kafka) > 0 {
			//todo someting
			hsm := engine.getConfigByOption(config.Conf.Kafka, opt.Kafka)
			//fmt.Println(hsm)
			//return nil
			in := <-zgokafka.InitKafka(hsm)
			Kafka = in
		}
		// 从local初始化缓存模块
		in := <-zgocache.InitCache(cacheCh)
		Cache = in

		Log = zgolog.InitLog(config.Conf.Project)
		LogWatch <- &config.CacheConfig{
			DbType: config.Conf.Log.DbType,
			Label:  config.Conf.Log.Label,
			Start:  config.Conf.Log.Start,
		}

	} else {

		var cc *config.CacheConfig
		for _, v := range resKvs {
			go func(v *mvccpb.KeyValue) {
				mk := string(v.Key)
				smk := strings.Split(mk, "/")
				b := v.Value
				if smk[3] == "cache" { //如果cache配置
					cm := config.CacheConfig{}
					err := zgoutils.Utils.Unmarshal(b, &cm)
					if err != nil {
						fmt.Println("反序列化当前值失败", mk)
					}
					config.Conf.Cache.Label = cm.Label
					config.Conf.Cache.Rate = cm.Rate
					config.Conf.Cache.Start = cm.Start
					config.Conf.Cache.TcType = cm.TcType
					config.Conf.Cache.DbType = cm.DbType

					// 从etcd初始化缓存模块
					out := zgocache.InitCache(cacheCh)
					go func() {
						for v := range out {
							fmt.Println("InitCache Success")
							Cache = v
						}
					}()
				} else if smk[3] == "log" { //init log存储配置 by etcd
					cm := config.CacheConfig{}
					err := zgoutils.Utils.Unmarshal(b, &cm)
					if err != nil {
						fmt.Println("反序列化当前值失败", mk)
					}

					Log = zgolog.InitLog(config.Conf.Project)
					config.Conf.Log.DbType = cm.DbType
					config.Conf.Log.Label = cm.Label
					config.Conf.Log.Start = cm.Start

					cc = &config.CacheConfig{
						DbType: config.Conf.Log.DbType,
						Label:  config.Conf.Log.Label,
						Start:  config.Conf.Log.Start,
					}
					LogWatch <- cc
					//fmt.Println("====log init by etcd config====", smk)

				} else if smk[1] == "project" && smk[2] == opt.Project { //init conn config by etcd

					var m []config.ConnDetail
					err := zgoutils.Utils.Unmarshal(b, &m)
					if err != nil {
						fmt.Println("反序列化当前值失败", mk)
					}
					//tmp.Key = smk[4]
					//tmp.Values = m

					var hsm = make(map[string][]*config.ConnDetail)
					for _, vv := range m {
						pvv := vv
						hsm[smk[4]] = append(hsm[smk[4]], &pvv)
						fmt.Printf("\n**********************资源项: %s **************************\n", smk[3])
						fmt.Printf("描述: %s\n", pvv.C)
						fmt.Printf("Label: %s\n", smk[4])
						fmt.Printf("Host: %s\n", pvv.Host)
						fmt.Printf("Port: %d\n", pvv.Port)
						fmt.Printf("DbName: %s\n", pvv.DbName)
					}
					initComponent(hsm, smk[3], smk[4])

				}
			}(v)
		}

	}

	return nil
}

//getConfigByOption 把zgo_start中的[]和config中的map进行match并取到关系
func (e *engine) getConfigByOption(lds []config.LabelDetail, us []string) map[string][]*config.ConnDetail {
	m := make(map[string][]*config.ConnDetail)
	for _, label := range us {
		//v == label_bj 用户传来的label，它并不知道具体的连接地址
		//v == label_sh 用户传来的label，它并不知道具体的连接地址
		for _, v := range lds {
			if label == v.Key {
				var tmp []*config.ConnDetail
				for k, _ := range v.Values {
					tmp = append(tmp, &v.Values[k])
				}
				m[v.Key] = tmp
			}
		}
	}
	return m
}

//定义外部使用的类型
type (
	NsqMessage        = *nsq.Message
	PartitionConsumer = kafkaCluter.PartitionConsumer
)

var (
	Kafka zgokafka.Kafkaer
	Nsq   zgonsq.Nsqer
	Mongo zgomongo.Mongoer
	Es    zgoes.Eser
	Grpc  zgogrpc.Grpcer
	Redis zgoredis.Rediser
	Pika  zgopika.Pikaer
	Mysql zgomysql.Mysqler
	Log   zgolog.Logger
	Cache zgocache.Cacher
	Http  = zgohttp.New()

	Utils = zgoutils.New()
	File  = zgofile.New()
)
