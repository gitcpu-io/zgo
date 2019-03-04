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
	Mongo        []LabelDetail
	Nsq          []LabelDetail
	Redis        []LabelDetail
	Pika         []LabelDetail
	Mysql        []LabelDetail
	Kafka        []LabelDetail
	CityDbConfig map[string]map[string]string
)

func InitConfig(e string) {
	if e == "local" {
		initConfig(e)
	} else {
		//用etcd
		InitConfigByEtcd()
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
	Mongo = acfg.Mongo
	Redis = acfg.Redis
	Pika = acfg.Pika
	Kafka = acfg.Kafka
	Mysql = acfg.Mysql
	CityDbConfig = acfg.CityDbConfig

	fmt.Printf("zgo engine %s is started on the ... %s\n", Version, Env)

	//fmt.Println(cf)

	//viper.SetConfigFile(cf)
	//err := viper.ReadInConfig()
	//if err != nil {
	//	panic(err)
	//}

	//for _, v := range acfg.Nsq  {
	//	if "label_bj" == v.Key {
	//		Nsq[v.Key] = v.Values
	//		for _, vv := range v.Values  {
	//			fmt.Println(vv.Host,vv.PoolSize)
	//		}
	//	}
	//}

	//Env = sjson.Get("env").MustString()
	//FormateJsonToMap(sjson, "nsq", Nsq)
	//FormateJsonToMap(sjson, "es", Es)
	//FormateJsonToMap(sjson, "mongo", Mongo)
}

//func FormateJsonToMap(sjson *simplejson.Json, name string, m map[string][]comm.Clabels) {
//	nsqJson, _ := sjson.Get(name).Map()
//
//	var cs []string
//	for k, _ := range nsqJson {
//		cs = append(cs, k)
//	}
//	for _, label := range cs {
//		//v == label_bj 用户传来的label，它并不知道具体的连接地址
//		//v == label_sh 用户传来的label，它并不知道具体的连接地址
//
//		tmpClabels := []comm.Clabels{}
//
//		if km, ok := nsqJson[label]; ok {
//			b, _ := json.Marshal(km)
//			json.Unmarshal(b, &tmpClabels)
//
//			m[label] = tmpClabels
//
//		}
//
//		//for k, v := range tmpClabels {
//		//	fmt.Println(k,v)
//		//}
//
//	}
//}
