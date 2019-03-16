package zgo

import (
	"errors"
	"fmt"
	"git.zhugefang.com/gocore/zgo/config"
	"git.zhugefang.com/gocore/zgo/zgoes"
	"git.zhugefang.com/gocore/zgo/zgokafka"
	"git.zhugefang.com/gocore/zgo/zgolog"
	"git.zhugefang.com/gocore/zgo/zgomongo"
	"git.zhugefang.com/gocore/zgo/zgomysql"
	"git.zhugefang.com/gocore/zgo/zgonsq"
	"git.zhugefang.com/gocore/zgo/zgopika"
	"git.zhugefang.com/gocore/zgo/zgoredis"
	"go.etcd.io/etcd/mvcc/mvccpb"
	"strings"
)

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

func (opt *Options) init() ([]*mvccpb.KeyValue, chan *config.CacheConfig, error) {
	//init config
	if opt.Env == "" {
		opt.Env = config.Local
	} else {
		if opt.Env != config.Local && opt.Env != config.Dev && opt.Env != config.Qa && opt.Env != config.Pro {
			return nil, nil, errors.New("error env,must be local/dev/qa/pro !")
		}
		if opt.Project == "" {
			return nil, nil, errors.New("u msut input your Project name to zgo.Engine func .")
		}
	}

	//如果inch有值表示启用了etcd为配置中心，并watch了key，等待变更ing...
	resKvs, inch, cacheCh, logCh, delConnCh, delCacheAndLogCh := config.InitConfig(opt.Env, opt.Project)

	opt.watchPutConn(inch)
	opt.watchDeleteConn(delConnCh)

	opt.watchDeleteCacheAndLog(delCacheAndLogCh)

	go func() {
		if logCh != nil {
			for cm := range logCh {
				//KEY: zgo/project/项目名/mysql/label名字
				var keyType string

				fmt.Println(keyType, "log,有变化开始init again", cm)

				Log = zgolog.InitLog(opt.Project)
				config.Conf.Log.DbType = cm.DbType
				config.Conf.Log.Label = cm.Label
				config.Conf.Log.Start = cm.Start

				cc := &config.CacheConfig{
					DbType: config.Conf.Log.DbType,
					Label:  config.Conf.Log.Label,
					Start:  config.Conf.Log.Start,
				}
				zgolog.LogWatch <- cc
			}
		}

	}()

	if opt.Project == "" {
		opt.Project = config.Conf.Project
	}
	if opt.Loglevel == "" {
		opt.Loglevel = config.Conf.Loglevel
	}

	return resKvs, cacheCh, nil
}

// watchDeleteConn 监听从etcd中删除的资源key，连接类型
func (opt *Options) watchDeleteConn(ch chan map[string][]*config.ConnDetail) {
	go func() {
		if ch != nil {
			//KEY: zgo/project/项目名/mysql/label名字
			for h := range ch {
				for k, v := range h {
					smk := strings.Split(k, "/")
					keyType := smk[3]
					sKey := smk[4]
					fmt.Println("[destroy conn]watchDeleteConn", keyType, sKey, v)
					//[destroy conn]watchDeleteConn nsq nsq_label_bj [0xc0004e6840 0xc0004e68f0]

					opt.destroyConn(keyType, sKey, v)
				}
			}
		}
	}()
}

// destroyConn 具体删除操作
func (opt *Options) destroyConn(keyType, sKey string, details []*config.ConnDetail) {
	switch keyType {
	case config.EtcTKMysql:
	case config.EtcTKMongo:
	case config.EtcTKRedis:
	case config.EtcTKPia:
	case config.EtcTKNsq:
	case config.EtcTKKafka:
	case config.EtcTKEs:
	case config.EtcTKEtc:
	}
}

// watchDeleteCacheAndLog 监听删除etcd中的 cache和log类型的key
func (opt *Options) watchDeleteCacheAndLog(ch chan map[string]*config.CacheConfig) {
	go func() {
		if ch != nil {
			//KEY: zgo/project/项目名/mysql/label名字
			for h := range ch {
				for k, v := range h {
					keyType := strings.Split(k, "/")[3]
					fmt.Println("[destroy]watchDeleteCacheAndLog:", keyType, v)
					//[destroy]watchDeleteCacheAndLog: log &{日志存储 0 /tmp 1 file 0}

					opt.destroyCacheAndLog(keyType, v)
				}
			}
		}
	}()
}

// destroyCacheAndLog 具体删除操作
func (opt *Options) destroyCacheAndLog(keyType string, cf *config.CacheConfig) {

	switch keyType {
	case config.EtcTKCache:
		//如果delete是cache todo something
	case config.EtcTKLog:
		//如果delete是log todo something
		config.Conf.Log.DbType = cf.DbType
		config.Conf.Log.Label = cf.Label
		config.Conf.Log.Start = 0

		cc := &config.CacheConfig{
			DbType: cf.DbType,
			Label:  cf.Label,
			Start:  0,
		}
		zgolog.LogWatch <- cc
	}
}

// watchPutConn 监听保存到etcd中的资源key，连接类型
func (opt *Options) watchPutConn(inch chan map[string][]*config.ConnDetail) {
	go func() {
		if inch != nil {
			for h := range inch {
				//KEY: zgo/project/项目名/mysql/label名字
				for k, _ := range h {
					keyType := strings.Split(k, "/")[3]
					mysqlLabel := strings.Split(k, "/")[4]
					hsm := make(map[string][]*config.ConnDetail)
					for mkey, v := range h {
						key := strings.Split(mkey, "/")[4] //改变label，去掉前缀
						hsm[key] = v
					}
					fmt.Println("[init again]watchPutConn:", keyType, hsm)
					//[init again]watchPutConn: nsq map[nsq_label_bj:[0xc0004e62c0 0xc0004e6420]]

					initConn(keyType, hsm, mysqlLabel)
				}
			}
		}

	}()
}

// initConn具体的连接操作
func initConn(keyType string, hsm map[string][]*config.ConnDetail, mysqlLabel string) {
	switch keyType {
	case config.EtcTKMysql:
		//init mysql again
		// 配置信息： 城市和数据库的关系
		cdc := config.Conf.CityDbConfig
		if len(hsm) > 0 {
			zgomysql.InitMysqlService(hsm, cdc)
			var err error
			Mysql, err = zgomysql.MysqlService(mysqlLabel)

			if err != nil {
				fmt.Println(err)
			}
		}

	case config.EtcTKMongo:
		//init mongo again
		if len(hsm) > 0 {
			in := <-zgomongo.InitMongo(hsm)
			Mongo = in
		}

	case config.EtcTKRedis:
		//init redis again
		if len(hsm) > 0 {
			in := <-zgoredis.InitRedis(hsm)
			Redis = in
		}

	case config.EtcTKPia:
		//init pika again
		if len(hsm) > 0 {
			in := <-zgopika.InitPika(hsm)
			Pika = in
		}

	case config.EtcTKNsq:
		//init nsq again
		if len(hsm) > 0 {
			in := <-zgonsq.InitNsq(hsm)
			Nsq = in

		}

	case config.EtcTKKafka:
		//init kafka again
		if len(hsm) > 0 {
			in := <-zgokafka.InitKafka(hsm)
			Kafka = in
		}

	case config.EtcTKEs:
		//init es again
		if len(hsm) > 0 {
			in := <-zgoes.InitEs(hsm)
			Es = in
		}

	case config.EtcTKEtc:
		//init etcd again
	}
}

func getMatchConfig(lds map[string][]*config.ConnDetail, us []string) map[string][]*config.ConnDetail {
	m := make(map[string][]*config.ConnDetail)
	for _, label := range us {
		//v == label_bj 用户传来的label，它并不知道具体的连接地址
		//v == label_sh 用户传来的label，它并不知道具体的连接地址
		for k, v := range lds {
			if label == k {
				m[label] = v
			}
		}
	}
	return m
}
