package zgo

import "git.zhugefang.com/gocore/zgo.git/config"

type Options struct {
	Env      string   `json:"env"`
	Project  string   `json:"project"`
	Loglevel string   `json:"loglevel"`
	Mongo    []string `json:"mongo"`
	Mysql    []string `json:"mysql"`
	Es       []string `json:"es"`
	Redis    []string `json:"redis"`
	Pika     []string `json:"pika"`
	Kafka    []string `json:"kafka"`
	Nsq      []string `json:"nsq"`
}

func (opt *Options) init() {
	if opt.Env == "" {
		opt.Env = "local"
	}
	if opt.Project == "" {
		opt.Project = config.Project
	}
	if opt.Loglevel == "" {
		opt.Loglevel = config.Loglevel
	}

	//init config
	config.InitConfig(opt.Env)
}
