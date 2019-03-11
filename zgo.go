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
	"github.com/nsqio/go-nsq"
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
	ladech, cacheCh, err := opt.init() //把zgo_start中用户定义的，映射到zgo的内存变量上
	if err != nil {
		return err
	}

	if opt.Env == "local" {
		if len(opt.Mongo) > 0 {
			//todo someting
			hsm := engine.getConfigByOption(config.Mongo, opt.Mongo)
			//fmt.Println("--zgo.go--",config.Mongo, opt.Mongo, hsm)
			in := <-zgomongo.InitMongo(hsm)
			Mongo = in
		}

		if len(opt.Mysql) > 0 {
			//todo someting
			hsm := engine.getConfigByOption(config.Mysql, opt.Mysql)
			fmt.Println(hsm)
			// 配置信息： 城市和数据库的关系
			cdc := config.CityDbConfig
			zgomysql.InitMysqlService(hsm, cdc)
			var err error
			Mysql, err = zgomysql.MysqlService(opt.Mysql[0])
			if err != nil {
				fmt.Println(err)
			}

		}
		if len(opt.Es) > 0 {
			hsm := engine.getConfigByOption(config.Es, opt.Es)
			in := <-zgoes.InitEs(hsm)
			Es = in
		}
		if len(opt.Redis) > 0 {
			//todo someting
			hsm := engine.getConfigByOption(config.Redis, opt.Redis)
			//fmt.Println(hsm)
			in := <-zgoredis.InitRedis(hsm)
			Redis = in
		}
		if len(opt.Pika) > 0 {
			//todo someting
			hsm := engine.getConfigByOption(config.Pika, opt.Pika)
			//fmt.Println(hsm)
			in := <-zgopika.InitPika(hsm)
			Pika = in
		}
		if len(opt.Nsq) > 0 { //>0表示用户要求使用nsq
			hsm := engine.getConfigByOption(config.Nsq, opt.Nsq)
			//fmt.Println("===zgo.go==", hsm)
			//return nil
			in := <-zgonsq.InitNsq(hsm)
			Nsq = in
		}
		if len(opt.Kafka) > 0 {
			//todo someting
			hsm := engine.getConfigByOption(config.Kafka, opt.Kafka)
			//fmt.Println(hsm)
			//return nil
			in := <-zgokafka.InitKafka(hsm)
			Kafka = in
		}
		// 从local初始化缓存模块
		in := <-zgocache.InitCache(cacheCh)
		Cache = in

	} else {

		go func() { //初始化时从etcd配置中读取
			for v := range ladech {
				//var tmp config.LabelDetail
				mk := string(v.Key)
				smk := strings.Split(mk, "/")
				b := v.Value

				if smk[3] == "cache" { //如果cache配置
					cm := config.CacheConfig{}
					err := zgoutils.Utils.Unmarshal(b, &cm)
					if err != nil {
						fmt.Println("反序列化当前值失败", mk)
					}
					config.Cache.Label = cm.Label
					config.Cache.Rate = cm.Rate
					config.Cache.Start = cm.Start
					config.Cache.TcType = cm.TcType
					config.Cache.DbType = cm.DbType

					// 从etcd初始化缓存模块
					out := zgocache.InitCache(cacheCh)
					go func() {
						for v := range out {
							fmt.Println("InitCache Success")
							Cache = v
						}
					}()

				} else if smk[3] == "log" { //log存储配置

					fmt.Println("====log init by etcd config====", smk)

				} else if smk[1] == "project" && smk[2] == opt.Project {

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
						fmt.Printf("\n**********************资源ID: %s **************************\n", smk[4])
						fmt.Printf("描述: %s\n", pvv.C)
						fmt.Printf("Host: %s\n", pvv.Host)
						fmt.Printf("Port: %d\n", pvv.Port)
						fmt.Printf("DbName: %s\n", pvv.DbName)
					}
					initComponent(hsm, smk[3], smk[4])

				}

			}
		}()
	}

	//初始化GRPC
	Grpc = zgogrpc.GetGrpc()

	if opt.Project != "" {
		config.Project = opt.Project
	}
	if opt.Loglevel != "" {
		config.Loglevel = opt.Loglevel
	}

	Log = zgolog.Newzgolog()

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
	NsqMessage = *nsq.Message
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
	Http  = zgohttp.NewHttp()

	Utils = zgoutils.NewUtils()
	File  = zgofile.NewLocal()
)
