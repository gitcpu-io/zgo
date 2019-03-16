package config

import (
	"bytes"
	"errors"
	"fmt"
	"git.zhugefang.com/gocore/zgo/zgoutils"
	"go.etcd.io/etcd/mvcc/mvccpb"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
)

const (
	Version       = "0.5.0"        //zgo版本号
	ProjectPrefix = "zgo/project/" //读取ETCD配置时prefix
	FileStoreType = "local"        //文件存储类型
	FileStoreHome = "/tmp"         //文件存储目录
	Local         = "local"        //本地开发环境标识
	Dev           = "dev"          //开发联调环境标识
	Qa            = "qa"           //QA测试环境标识
	Pro           = "pro"          //生产环境标识
)

var (
	DevEtcdHosts = []string{ //开发联调ETCD地址
		//"123.56.173.28:2380",
		"localhost:2381",
	}
	QaEtcdHosts = []string{ //QA环境ETCD地址
		"123.56.173.28:2380",
	}
	ProEtcdHosts = []string{ //生产环境ETCD地址
		"123.56.173.28:2380",
	}
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

var Conf *allConfig

func InitConfig(e, project string) ([]*mvccpb.KeyValue, chan map[string][]*ConnDetail, chan *CacheConfig, chan *CacheConfig) {

	ReadFileByConfig(e, project)

	if e != Local {
		//用etcd
		return InitConfigByEtcd(project)
	}
	return nil, nil, nil, nil
}

func ReadFileByConfig(e, project string) {
	var cf string
	switch e {
	case Local:
		_, f, _, ok := runtime.Caller(1)
		if !ok {
			panic(errors.New("Can not get current file info"))
		}
		cf = fmt.Sprintf("%s/%s.json", filepath.Dir(f), e)

		bf, err := ioutil.ReadFile(cf)
		if err != nil {
			panic(err)
		}

		Conf = &allConfig{}
		err = zgoutils.Utils.Unmarshal(bf, Conf)
		if err != nil {
			panic(err)
		}

		//Conf = LoadConfig(cf)

	case Dev:
		Conf = &allConfig{
			Env:       e,
			Project:   project,
			EtcdHosts: DevEtcdHosts,
			File: FileStore{
				Type: FileStoreType,
				Home: FileStoreHome,
			},
		}
	case Qa:
		Conf = &allConfig{
			Env:       e,
			Project:   project,
			EtcdHosts: QaEtcdHosts,
			File: FileStore{
				Type: FileStoreType,
				Home: FileStoreHome,
			},
		}
	case Pro:
		Conf = &allConfig{
			Env:       e,
			Project:   project,
			EtcdHosts: ProEtcdHosts,
			File: FileStore{
				Type: FileStoreType, //以后生产环境可以存到aws s3，在这里直接更改
				Home: FileStoreHome,
			},
		}
	}

	fmt.Printf("zgo engine %s is started on the ... %s\n", Version, Conf.Env)

}

// LoadConfig 暂时不用
func LoadConfig(path string) *allConfig {
	var config allConfig
	config_file, err := os.Open(path)
	if err != nil {
		emit("Failed to open config file '%s': %s\n", path, err)
		return &config
	}

	fi, _ := config_file.Stat()
	if size := fi.Size(); size > (10 << 20) {
		emit("config file (%q) size exceeds reasonable limit (%d) - aborting", path, size)
		return &config // REVU: shouldn't this return an error, then?
	}

	if fi.Size() == 0 {
		emit("config file (%q) is empty, skipping", path)
		return &config
	}

	buffer := make([]byte, fi.Size())
	_, err = config_file.Read(buffer)
	//emit("\n %s\n", buffer)

	buffer, err = StripComments(buffer) //去掉注释
	if err != nil {
		emit("Failed to strip comments from json: %s\n", err)
		return &config
	}

	buffer = []byte(os.ExpandEnv(string(buffer))) //特殊

	err = zgoutils.Utils.Unmarshal(buffer, &config) //解析json格式数据
	if err != nil {
		emit("Failed unmarshalling json: %s\n", err)
		return &config
	}
	return &config
}

func StripComments(data []byte) ([]byte, error) {
	data = bytes.Replace(data, []byte("\r"), []byte(""), 0) // Windows
	lines := bytes.Split(data, []byte("\n"))                //split to muli lines
	filtered := make([][]byte, 0)

	for _, line := range lines {
		match, err := regexp.Match(`^\s*#`, line)
		if err != nil {
			return nil, err
		}
		if !match {
			filtered = append(filtered, line)
		}
	}

	return bytes.Join(filtered, []byte("\n")), nil
}

func emit(msgfmt string, args ...interface{}) {
	fmt.Printf(msgfmt, args...)
}
