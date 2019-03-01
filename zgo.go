package zgo

import (
	"fmt"
	"git.zhugefang.com/gocore/zgo.git/config"
	"git.zhugefang.com/gocore/zgo.git/zgoes"
	"git.zhugefang.com/gocore/zgo.git/zgofile"
	"git.zhugefang.com/gocore/zgo.git/zgogrpc"
	"git.zhugefang.com/gocore/zgo.git/zgokafka"
	"git.zhugefang.com/gocore/zgo.git/zgolog"
	"git.zhugefang.com/gocore/zgo.git/zgomongo"
	"git.zhugefang.com/gocore/zgo.git/zgomysql1"
	"git.zhugefang.com/gocore/zgo.git/zgonsq"
	"git.zhugefang.com/gocore/zgo.git/zgoredis"
	"git.zhugefang.com/gocore/zgo.git/zgoutils"
	"git.zhugefang.com/gocore/zgo.git/zgozoneinfo"
	"github.com/nsqio/go-nsq"
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

	if len(opt.Mongo) > 0 {
		//todo someting
		hsm := engine.getConfigByOption(config.Mongo, opt.Mongo)
		//fmt.Println(hsm)
		in := <-zgomongo.InitMongo(hsm)
		Mongo = in
	}
	if len(opt.Mysql) > 0 {
		//todo someting
		hsm := engine.getConfigByOption(config.Mysql, opt.Mongo)
		fmt.Println(hsm)
		zgomysql1.InitMysql(hsm)
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
	}
	if len(opt.Nsq) > 0 { //>0表示用户要求使用nsq
		hsm := engine.getConfigByOption(config.Nsq, opt.Nsq)
		//fmt.Println(hsm)
		//return nil
		in := <-zgonsq.InitNsq(hsm)
		Nsq = in
	}
	if len(opt.Kafka) > 0 {
		//todo someting
		hsm := engine.getConfigByOption(config.Kafka, opt.Kafka)
		//fmt.Println(hsm)
		//return nil
		k := <-zgokafka.InitKafka(hsm)
		Kafka = k
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
	Grpc  = zgogrpc.Grpc()
	Redis zgoredis.Rediser

	Mysql    = zgomysql1.Mysql("")
	File     = zgofile.NewLocal()
	Utils    = zgoutils.NewUtils()
	Log      = zgolog.Newzgolog()
	ZoneInfo = zgozoneinfo.NewZoneInfo()
)
