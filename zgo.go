package zgo

import (
	"fmt"
	"git.zhugefang.com/gocore/zgo/config"
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
	"log"
)

type engine struct {
	opt *Options
}

//New init zgo engine
func Engine(opt *Options) *engine {
	engine := &engine{
		opt: opt,
	}
	opt.init() //把zgo_start中用户定义的，映射到zgo的内存变量上

	//初始化GRPC
	Grpc = zgogrpc.GetGrpc()

	if len(opt.Mongo) > 0 {
		//todo someting
		hsm := engine.getConfigByOption(config.Mongo, opt.Mongo)
		//fmt.Println(hsm)
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
			log.Fatalf(err.Error())
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
		fmt.Println(config.Nsq)
		fmt.Println("=====", opt.Nsq)
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

	if opt.Project != "" {
		config.Project = opt.Project
	}
	if opt.Loglevel != "" {
		config.Loglevel = opt.Loglevel
	}

	return engine
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
				for _, vv := range v.Values {
					tmp = append(tmp, &vv)
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
)
