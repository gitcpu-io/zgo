package zgo

import "git.zhugefang.com/gocore/zgo.git/config"

type Options struct {
	Env   string
	mongo []string
	mysql []string
	es    []string
	redis []string
	pika  []string
	kafka []string
	Nsq   []string
}

func (opt *Options) init() {
	if opt.Env == "" {
		opt.Env = "local"
	}

	//init config
	config.InitConfig(opt.Env)
}
