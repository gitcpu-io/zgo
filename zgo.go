package zgo

import (
	"git.zhugefang.com/gocore/zgo.git/config"
	"git.zhugefang.com/gocore/zgo.git/zgoes"
	"git.zhugefang.com/gocore/zgo.git/zgomongo"
	"git.zhugefang.com/gocore/zgo.git/zgonsq"
	"git.zhugefang.com/gocore/zgo.git/zgoredis"
	"github.com/nsqio/go-nsq"
)

var Version = "0.1"

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
		zgomongo.InitMongo(hsm)
	}
	if len(opt.Mysql) > 0 {
		//todo someting
		//hsm := engine.getConfigByOption(config.Mysql, opt.Mongo)
		//fmt.Println(hsm)
		//zgomysql.InitMysqlService(hsm)
	}
	if len(opt.Es) > 0 {
		hsm := engine.getConfigByOption(config.Es, opt.Es)
		zgoes.InitEs(hsm)
	}
	if len(opt.Redis) > 0 {
		//todo someting
		hsm := engine.getConfigByOption(config.Redis, opt.Redis)
		//fmt.Println(hsm)
		zgoredis.InitRedis(hsm)
	}
	if len(opt.Pika) > 0 {
		//todo someting
	}
	if len(opt.Nsq) > 0 { //>0表示用户要求使用nsq
		hsm := engine.getConfigByOption(config.Nsq, opt.Nsq)
		//fmt.Println(hsm)
		//return nil
		zgonsq.InitNsq(hsm)
	}
	if len(opt.Kafka) > 0 {
		//todo someting
		//hsm := engine.getConfigByOption(config.Kafka, opt.Kafka)
		//fmt.Println(hsm)
		//return nil
		//zgokafka.InitNsq(hsm)
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
	Nsq   = zgonsq.Nsq("")
	Mongo = zgomongo.Mongo("")
	Es    = zgoes.Es("")
	Redis = zgoredis.Redis("")
)
