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
	"strings"
)

const (
	//********************************以下是 loglevel 千万不要换顺序********************************
	Debug = iota //0
	Info         //1
	Warn         //2
	Error        //3

	Version       = "1.0.4"       //zgo版本号
	ProjectPrefix = "zgo/project" //读取ETCD配置时prefix
	FileStoreType = "local"       //文件存储类型
	FileStoreHome = "/tmp"        //文件存储目录
	Local         = "local"       //本地开发环境标识
	Dev           = "dev"         //开发联调环境标识
	Qa            = "qa"          //QA测试环境标识
	Pro           = "pro"         //生产环境标识
	Pro2          = "pro2"        //生产环境标识

	//********************************以下是 etcd监听常量********************************
	EtcTKCache    = "cache"
	EtcTKLog      = "log"
	EtcTKMysql    = "mysql"
	EtcTKPostgres = "postgres"
	EtcTKNeo4j    = "neo4j"
	EtcTKMongo    = "mongo"
	EtcTKMgo      = "mgo"
	EtcTKRedis    = "redis"
	EtcTKPia      = "pika"
	EtcTKNsq      = "nsq"
	EtcTKKafka    = "kafka"
	EtcTKEs       = "es"
	EtcTKEtcd     = "etcd"

	//****************************以下是 mongodb bulk write常量**************************
	InsertOne  = "insertOne"
	UpdateOne  = "updateOne"
	ReplaceOne = "replaceOne"
	DeleteOne  = "deleteOne"
	UpdateMany = "updateMany"
	DeleteMany = "deleteMany"
)

var Levels = []string{"debug", "info", "warn", "error"}

var (
	DevEtcHosts = []string{ //开发联调ETCD地
		//"10.45.146.41:2380", //测试时使用内网ip
		"47.95.20.12:2381", //如果本机联调，想用测试机的etcd可以使用公网ip
		//"localhost:2381",
	}
	QaEtcHosts = []string{ //QA环境ETCD地址，同正式
		//"47.95.20.12:2381",
		"10.24.188.182:2381",
	}
	ProEtcHosts = []string{ //生产环境ETCD地址，需要使用内部dns解析，在k8s的worker节点配置/etc/hosts下面的域名和真实的etcd的ip
		"10.25.96.1:2379",
		"10.26.100.217:2379",
		"10.26.162.67:2379",
	}
	cityDbConfig = map[string]map[string]string{
		"sell": {
			"bj":  "1",
			"nj":  "1",
			"sh":  "1",
			"cd":  "1",
			"tj":  "1",
			"cq":  "1",
			"heb": "1",
		},
	}
)

type ConnDetail struct {
	C           string `json:"c,omitempty"`
	Host        string `json:"host,omitempty"`
	Port        int    `json:"port,omitempty"`
	ConnSize    int    `json:"connSize"`
	PoolSize    int    `json:"poolSize"`
	MaxIdleSize int    `json:"maxIdleSize,omitempty"` // mysql 最大空闲连接数
	MaxOpenConn int    `json:"maxOpenConn,omitempty"` // mysql 最大可用连接数
	Username    string `json:"username,omitempty"`
	Password    string `json:"password,omitempty"`
	Db          int    `json:"db,omitempty"`
	T           string `json:"t,omitempty"` // w 写入 r 只读
	Prefix      string `json:"prefix,omitempty"`
	DbName      string `json:"dbName,omitempty"`  // 数据库名称
	LogMode     int    `json:"logMode,omitempty"` // 日志类型
}

type CacheConfig struct {
	//same as LogConfig so 共用一个struct
	LogLevel int    `json:"loglevel,omitempty"`
	C        string `json:"c,omitempty"`
	Rate     int    `json:"rate,omitempty"`   // 缓存失效时间 倍率
	Label    string `json:"label,omitempty"`  // 缓存所需的 pikaLabel
	Start    int    `json:"start,omitempty"`  // 是否开启 1 开启 0关闭
	DbType   string `json:"dbType,omitempty"` // 数据库类型 默认pika
	TcType   int    `json:"tcType,omitempty"` // 降级缓存类型 1正常降级缓存 2转为普通缓存
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
	EtcdHosts    []string                     `json:"etcdHosts,omitempty"`
	Nsq          []LabelDetail                `json:"nsq,omitempty"`
	Mongo        []LabelDetail                `json:"mongo,omitempty"`
	Mgo          []LabelDetail                `json:"mgo,omitempty"`
	Mysql        []LabelDetail                `json:"mysql,omitempty"`
	Postgres     []LabelDetail                `json:"postgres,omitempty"`
	Neo4j        []LabelDetail                `json:"neo4j,omitempty"`
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

func InitConfig(env, project, etcdHosts string) ([]*mvccpb.KeyValue, chan map[string][]*ConnDetail, chan map[string]*CacheConfig, chan map[string][]*ConnDetail, chan map[string]*CacheConfig) {

	LoadConfig(env, project, etcdHosts)

	if env != Local {
		//用etcd的配置
		ec := EtcConfig{
			Key:       fmt.Sprintf("%s/%s", ProjectPrefix, project),
			Endpoints: Conf.EtcdHosts,
		}
		return ec.InitConfigByEtcd()
	}
	return nil, nil, nil, nil, nil
}

func LoadConfig(e, project, etcdHosts string) {
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
			EtcdHosts: DevEtcHosts,
			File: FileStore{
				Type: FileStoreType,
				Home: FileStoreHome,
			},
		}

	case Qa:
		Conf = &allConfig{
			Env:       e,
			Project:   project,
			EtcdHosts: QaEtcHosts,
			File: FileStore{
				Type: FileStoreType,
				Home: FileStoreHome,
			},
		}

	case Pro:
		Conf = &allConfig{
			Env:       e,
			Project:   project,
			EtcdHosts: ProEtcHosts,
			File: FileStore{
				Type: FileStoreType, //以后生产环境可以存到aws s3，在这里直接更改
				Home: FileStoreHome,
			},
		}
	case Pro2:
		Conf = &allConfig{
			Env:       e,
			Project:   project,
			EtcdHosts: ProEtcHosts,
			File: FileStore{
				Type: FileStoreType, //以后生产环境可以存到aws s3，在这里直接更改
				Home: FileStoreHome,
			},
		}
	}

	if etcdHosts != "" {
		Conf.EtcdHosts = strings.Split(etcdHosts, ",")
	}

	//default init city db config
	Conf.CityDbConfig = cityDbConfig

	fmt.Printf("zgo engine %s is started on the ... %s %s\n", Version, Conf.Env, Conf.EtcdHosts)

}

// LoadConfig 暂时不用
func LoadConfigByFile(path string) *allConfig {
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
