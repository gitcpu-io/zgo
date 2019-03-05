package zgo

import (
	"fmt"
	"git.zhugefang.com/gocore/zgo/config"
	"git.zhugefang.com/gocore/zgo/zgocache"
	"git.zhugefang.com/gocore/zgo/zgoes"
	"git.zhugefang.com/gocore/zgo/zgofile"
	"git.zhugefang.com/gocore/zgo/zgogrpc"
	"git.zhugefang.com/gocore/zgo/zgokafka"
	"git.zhugefang.com/gocore/zgo/zgolog"
	"git.zhugefang.com/gocore/zgo/zgomongo"
	"git.zhugefang.com/gocore/zgo/zgomysql"
	"git.zhugefang.com/gocore/zgo/zgonsq"
	"git.zhugefang.com/gocore/zgo/zgopika"
	"git.zhugefang.com/gocore/zgo/zgoredis"
	"git.zhugefang.com/gocore/zgo/zgoutils"
	"git.zhugefang.com/gocore/zgo/zgozoneinfo"
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
	ladech, err := opt.init() //把zgo_start中用户定义的，映射到zgo的内存变量上
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

	} else {

		go func() {
			for v := range ladech {
				var tmp config.LabelDetail
				var labelDetArr []config.LabelDetail
				mk := string(v.Key)
				smk := strings.Split(mk, "/")
				b := v.Value
				var m []config.ConnDetail
				err := zgoutils.Utils.Unmarshal(b, &m)
				if err != nil {
					fmt.Println("反序列化当前值失败", mk)
				}
				tmp.Key = smk[2]
				tmp.Values = m

				labelDetArr = append(labelDetArr, tmp)

				key := smk[1]

				//fmt.Println(smk[1],"-----",labelDetArr)

				switch key {
				case mysqlT:
					//init mysql again

				case mongoT:
					//init mongo again
					if len(opt.Mongo) > 0 {
						//todo someting
						hsm := engine.getConfigByOption(labelDetArr, opt.Mongo)
						//fmt.Println("--zgo.go--",labelDetArr, opt.Mongo, hsm)
						if len(hsm) > 0 {
							in := <-zgomongo.InitMongo(hsm)
							Mongo = in
						}

					}
				case redisT:
					//init redis again
					if len(opt.Redis) > 0 {
						//todo someting
						hsm := engine.getConfigByOption(labelDetArr, opt.Redis)
						//fmt.Println(hsm)
						if len(hsm) > 0 {
							in := <-zgoredis.InitRedis(hsm)
							Redis = in
						}

					}
				case pikaT:
					//init pika again
					if len(opt.Pika) > 0 {
						//todo someting
						hsm := engine.getConfigByOption(labelDetArr, opt.Pika)
						//fmt.Println(hsm)
						if len(hsm) > 0 {
							in := <-zgopika.InitPika(hsm)
							Pika = in
						}

					}
				case nsqT:
					//init nsq again
					if len(opt.Nsq) > 0 { //>0表示用户要求使用nsq
						hsm := engine.getConfigByOption(labelDetArr, opt.Nsq)
						//fmt.Println("===zgo.go==", hsm)
						//return nil
						if len(hsm) > 0 {
							in := <-zgonsq.InitNsq(hsm)
							Nsq = in
						}

					}
				case kafkaT:
					//init kafka again
					if len(opt.Kafka) > 0 {
						//todo someting
						hsm := engine.getConfigByOption(labelDetArr, opt.Kafka)
						//fmt.Println(hsm)
						//return nil
						if len(hsm) > 0 {
							in := <-zgokafka.InitKafka(hsm)
							Kafka = in
						}

					}
				case esT:
					//init es again
					if len(opt.Es) > 0 {
						hsm := engine.getConfigByOption(labelDetArr, opt.Es)
						if len(hsm) > 0 {
							in := <-zgoes.InitEs(hsm)
							Es = in
						}

					}
				case etcdT:
					//init etcd again

				}
				//fmt.Println(Nsq)
			}
		}()
	}

	//初始化GRPC
	Grpc = zgogrpc.GetGrpc()
	// 初始化缓存模块
	Cache = zgocache.InitCache()

	if opt.Project != "" {
		config.Project = opt.Project
	}
	if opt.Loglevel != "" {
		config.Loglevel = opt.Loglevel
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
	NsqMessage = *nsq.Message
)

var (
	Kafka zgokafka.Kafkaer
	Nsq   zgonsq.Nsqer
	Mongo zgomongo.Mongoer
	Es    zgoes.Eser

	Grpc zgogrpc.Grpcer

	Redis zgoredis.Rediser
	Pika  zgopika.Pikaer

	Mysql    zgomysql.MysqlServiceInterface
	File     = zgofile.NewLocal()
	Utils    = zgoutils.NewUtils()
	Log      = zgolog.Newzgolog()
	ZoneInfo = zgozoneinfo.NewZoneInfo()

	Cache zgocache.CacheServiceInterface
)
