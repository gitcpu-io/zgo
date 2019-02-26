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

type ConnDetail struct {
	C        string `json:"c"`
	Host     string `json:"host,omitempty"`
	Port     int    `json:"port,omitempty"`
	ConnSize int    `json:"connSize"`
	PoolSize int    `json:"poolSize"`
	Uri      string `json:"uri,omitempty"`
}
type LabelDetail struct {
	Key    string `json:"key"`
	Values []ConnDetail
}

type allConfig struct {
	Env   string        `json:"env"`
	Nsq   []LabelDetail `json:"nsq"`
	Mongo []LabelDetail `json:"mongo"`
	Mysql []LabelDetail `json:"mysql"`
	Redis []LabelDetail `json:"redis"`
	Kafka []LabelDetail `json:"kafka"`
	Es    []LabelDetail `json:"es"`
}

type Labelconns struct {
	Label string        `json:"label"`
	Hosts []*ConnDetail `json:"hosts"`
}

var (
	Env   string
	Es    []LabelDetail
	Mongo []LabelDetail
	Nsq   []LabelDetail
	Redis []LabelDetail
	Mysql []LabelDetail
	Kafka []LabelDetail
)

func InitConfig(e string) {
	initConfig(e)
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
	Nsq = acfg.Nsq
	Es = acfg.Es
	Mongo = acfg.Mongo
	Redis = acfg.Redis
	Kafka = acfg.Kafka
	Mysql = acfg.Mysql

	//fmt.Println(Nsq)

	fmt.Println("zgo engine is started on the ... ", Env)

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
