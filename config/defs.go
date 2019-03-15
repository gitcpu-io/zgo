package config

import (
	"errors"
	"fmt"
	"git.zhugefang.com/gocore/zgo/zgoutils"
	"go.etcd.io/etcd/mvcc/mvccpb"
	"io/ioutil"
	"path/filepath"
	"runtime"
)

type ConnDetail struct {
	C           string `json:"c"`
	Host        string `json:"host,omitempty"`
	Port        int    `json:"port,omitempty"`
	ConnSize    int    `json:"connSize"`
	PoolSize    int    `json:"poolSize"`
	MaxIdleSize int    `json:"maxIdleSize,omitempty"` // mysql 最大空闲连接数
	MaxOpenConn int    `json:"maxOpenConn,omitempty"` // mysql 最大可用连接数
	Uri         string `json:"uri,omitempty"`
	Username    string `json:"username,omitempty"`
	Password    string `json:"password,omitempty"`
	Db          int    `json:"db,omitempty"`
	T           string `json:"t,omitempty"` // w 写入 r 只读
	Prefix      string `json:"prefix,omitempty"`
	DbName      string `json:"dbName,omitempty"` // 数据库名称
}

type CacheConfig struct {
	//same as LogConfig so 共用一个struct
	C      string `json:"c"`
	Rate   int    `json:"rate,omitempty"`   // 缓存失效时间 倍率
	Label  string `json:"label"`            // 缓存所需的 pikaLabel
	Start  int    `json:"start"`            // 是否开启 1 开启 0关闭
	DbType string `json:"dbType"`           // 数据库类型 默认pika
	TcType int    `json:"tcType,omitempty"` // 降级缓存类型 1正常降级缓存 2转为普通缓存
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
	Version      string                       `json:"version"`
	Env          string                       `json:"env"`
	File         FileStore                    `json:"file,omitempty"`
	Project      string                       `json:"project"`
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
	Log          CacheConfig                  `json:"log"`
	CityDbConfig map[string]map[string]string `json:"cityDbConfig,omitempty"`
}

type Labelconns struct {
	Label string        `json:"label"`
	Hosts []*ConnDetail `json:"hosts"`
}

var Conf allConfig

func InitConfig(e, project string) ([]*mvccpb.KeyValue, chan map[string][]*ConnDetail, chan *CacheConfig, chan *CacheConfig) {
	ReadFileByConfig(e)

	if e != "local" {
		//用etcd
		return InitConfigByEtcd(project)
	}
	return nil, nil, nil, nil
}

func ReadFileByConfig(e string) {
	_, f, _, ok := runtime.Caller(1)
	if !ok {
		panic(errors.New("Can not get current file info"))
	}
	cf := fmt.Sprintf("%s/%s.json", filepath.Dir(f), e)

	bf, _ := ioutil.ReadFile(cf)
	Conf = allConfig{}
	err := zgoutils.Utils.Unmarshal(bf, &Conf)
	if err != nil {
		panic(err)
	}

	fmt.Printf("zgo engine %s is started on the ... %s\n", Conf.Version, Conf.Env)

}
