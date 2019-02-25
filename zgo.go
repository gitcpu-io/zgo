package zgo

import (
	"git.zhugefang.com/gocore/zgo.git/config"
	"git.zhugefang.com/gocore/zgo.git/zgomongo"
	"git.zhugefang.com/gocore/zgo.git/zgonsq"
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

	if len(opt.mongo) > 0 {
		//todo someting
		hsm := engine.getConfigByOption(config.Nsq, opt.Nsq)
		//fmt.Println(hsm)
		zgomongo.InitMongo(hsm)
	}

	if len(opt.mysql) > 0 {
		//todo someting

	}
	if len(opt.es) > 0 {
		//todo someting
	}
	if len(opt.redis) > 0 {
		//todo someting
	}
	if len(opt.pika) > 0 {
		//todo someting
	}
	if len(opt.Nsq) > 0 { //>0表示用户要求使用nsq
		hsm := engine.getConfigByOption(config.Nsq, opt.Nsq)
		//fmt.Println(hsm)
		zgonsq.InitNsq(hsm)
	}
	if len(opt.kafka) > 0 {
		//todo someting
	}

	return engine
}

//getConfigByOption 把zgo_start中的[]和config中的map进行match并取到关系
func (e *engine) getConfigByOption(cmap map[string][]string, us []string) map[string][]string {
	m := make(map[string][]string)
	for _, v := range us {
		//v == label_bj 用户传来的label，它并不知道具体的连接地址
		//v == label_sh 用户传来的label，它并不知道具体的连接地址
		if hp, ok := cmap[v]; ok {
			m[v] = hp
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
)
