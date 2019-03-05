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
	T           string `json:"db,omitempty"` // w 写入 r 只读
	Prefix      string `json:"prefix,omitempty"`
	Expire      int    `json:"prefix,omitempty"`     // 缓存失效时间 单位sec
	CacheLabel  string `json:"cacheLabel,omitempty"` // 缓存所需的 redisLabel
}
type LabelDetail struct {
	Key    string `json:"key"`
	Values []ConnDetail
}

type FileStore struct {
	Type string `json:"type"`
	Home string `json:"home"`
}

type allConfig struct {
	Env          string                       `json:"env"`
	File         FileStore                    `json:"file"`
	Project      string                       `json:"project"`
	Loglevel     string                       `json:"loglevel"`
	Nsq          []LabelDetail                `json:"nsq"`
	Mongo        []LabelDetail                `json:"mongo"`
	Mysql        []LabelDetail                `json:"mysql"`
	Redis        []LabelDetail                `json:"redis"`
	Pika         []LabelDetail                `json:"pika"`
	Kafka        []LabelDetail                `json:"kafka"`
	Es           []LabelDetail                `json:"es"`
	Etcd         []LabelDetail                `json:"etcd"`
	Cache        LabelDetail                  `json:"cache"`
	CityDbConfig map[string]map[string]string `json:"cityDbConfig"`
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
	Es           []LabelDetail
	Etcd         []LabelDetail
	Mongo        []LabelDetail
	Nsq          []LabelDetail
	Redis        []LabelDetail
	Pika         []LabelDetail
	Mysql        []LabelDetail
	Kafka        []LabelDetail
	Cache        LabelDetail
	CityDbConfig map[string]map[string]string
)

func InitConfig(e string) chan map[string][]*ConnDetail {
	if e == "local" {
		initConfig(e)
		return nil
	} else {
		//用etcd
		return InitConfigByEtcd()
	}
}

func initConfig(e string) {
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
	Nsq = acfg.Nsq
	Es = acfg.Es
	Etcd = acfg.Etcd
	Mongo = acfg.Mongo
	Redis = acfg.Redis
	Pika = acfg.Pika
	Kafka = acfg.Kafka
	Mysql = acfg.Mysql
	CityDbConfig = acfg.CityDbConfig

	fmt.Printf("zgo engine %s is started on the ... %s\n", Version, Env)

}
