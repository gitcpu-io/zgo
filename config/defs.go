package config

import (
	"errors"
	"fmt"
	"github.com/json-iterator/go"
	"io/ioutil"
	"path/filepath"
	"runtime"
)

var jsonIterator = jsoniter.ConfigCompatibleWithStandardLibrary

var Version = "0.0.1"

const (
	mysqlT = "mysql"
	mongoT = "mongo"
	redisT = "redis"
	pikaT  = "pika"
	nsqT   = "nsq"
	kafkaT = "kafka"
	esT    = "es"
	etcdT  = "etcd"
)

type ConnDetail struct {
	C           string `json:"c"`
	Host        string `json:"host,omitempty"`
	Port        int    `json:"port,omitempty"`
	ConnSize    int    `json:"connSize"`
	PoolSize    int    `json:"poolSize"`
	MaxIdleSize int    `json:"maxIdleSize,omitempty"`
	MaxOpenConn int    `json:"maxOpenConn,omitempty"`
	Uri         string `json:"uri,omitempty"`
	Username    string `json:"username,omitempty"`
	Password    string `json:"password,omitempty"`
	Db          int    `json:"db,omitempty"`
	T           string `json:"t,omitempty"` // w 写入 r 只读
	Prefix      string `json:"prefix,omitempty"`
}

type CacheConfig struct {
	Expire int    `json:"expire"`     // 缓存失效时间 单位sec
	Label  string `json:"cacheLabel"` // 缓存所需的 redisLabel
	Start  int    `json:"start"`      // 是否开启
	DbType string `json:"dbType"`     //
}

type LabelDetail struct {
	Key    string       `json:"key"`
	Values []ConnDetail `json:"values"`
}

type FileStore struct {
	Type string `json:"type"`
	Home string `json:"home"`
}

type allConfig struct {
	Env          string                       `json:"env,omitempty"`
	File         FileStore                    `json:"file,omitempty"`
	Project      string                       `json:"project,omitempty"`
	Loglevel     string                       `json:"loglevel,omitempty"`
	EtcdHosts    []string                     `json:"etcdHosts,omitempty"`
	Nsq          []LabelDetail                `json:"nsq,omitempty"`
	Mongo        []LabelDetail                `json:"mongo,omitempty"`
	Mysql        []LabelDetail                `json:"mysql,omitempty"`
	Redis        []LabelDetail                `json:"redis,omitempty"`
	Pika         []LabelDetail                `json:"pika,omitempty"`
	Kafka        []LabelDetail                `json:"kafka,omitempty"`
	Es           []LabelDetail                `json:"es,omitempty"`
	Etcd         []LabelDetail                `json:"etcd,omitempty"`
	Cache        CacheConfig                  `json:"cache"`
	CityDbConfig map[string]map[string]string `json:"cityDbConfig,omitempty"`
}

type Labelconns struct {
	Label string        `json:"label"`
	Hosts []*ConnDetail `json:"hosts"`
}

var (
	Env          string
	File         FileStore
	Project      string
	Loglevel     string
	EtcdHosts    []string
	Es           []LabelDetail
	Etcd         []LabelDetail
	Mongo        []LabelDetail
	Nsq          []LabelDetail
	Redis        []LabelDetail
	Pika         []LabelDetail
	Mysql        []LabelDetail
	Kafka        []LabelDetail
	Cache        CacheConfig
	CityDbConfig map[string]map[string]string
)

func InitConfig(e string) chan map[string][]*ConnDetail {
	ReadFildByConfig(e)

	if e != "local" {
		//用etcd
		return InitConfigByEtcd()
	}
	return nil
}

func ReadFildByConfig(e string) {
	_, f, _, ok := runtime.Caller(1)
	if !ok {
		panic(errors.New("Can not get current file info"))
	}
	cf := fmt.Sprintf("%s/%s.json", filepath.Dir(f), e)

	bf, _ := ioutil.ReadFile(cf)
	acfg := allConfig{}
	err := jsonIterator.Unmarshal(bf, &acfg)
	if err != nil {
		panic(err)
	}

	Env = acfg.Env
	File = acfg.File
	Project = acfg.Project
	Loglevel = acfg.Loglevel
	EtcdHosts = acfg.EtcdHosts
	Nsq = acfg.Nsq
	Es = acfg.Es
	Etcd = acfg.Etcd
	Mongo = acfg.Mongo
	Redis = acfg.Redis
	Pika = acfg.Pika
	Kafka = acfg.Kafka
	Mysql = acfg.Mysql
	CityDbConfig = acfg.CityDbConfig
	Cache = acfg.Cache

	fmt.Printf("zgo engine %s is started on the ... %s\n", Version, Env)

}
