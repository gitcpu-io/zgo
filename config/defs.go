package config

import (
  "fmt"
  "github.com/gitcpu-io/zgo/zgoutils"
  "go.etcd.io/etcd/api/v3/mvccpb"
  "io/ioutil"
  "strings"
)

const (
  //********************************以下是 loglevel 千万不要换顺序********************************
  Debug = iota //0
  Info         //1
  Warn         //2
  Error        //3

  Version       = "v1.1.1"      //zgo版本号
  ProjectPrefix = "zgo/project" //读取ETCD配置时prefix
  FileStoreType = "local"       //文件存储类型
  FileStoreHome = "/tmp"        //文件存储目录
  Local         = "local"       //本地开发环境标识
  Dev           = "dev"         //开发联调环境标识
  Qa            = "qa"          //QA测试环境标识
  Pro           = "pro"         //生产环境标识
  Container     = "container"   //k8s生产环境标识
  K8s           = "k8s"         //k8s生产环境标识

  //********************************以下是 etcd监听常量********************************
  EtcTKCache      = "cache"
  EtcTKLog        = "log"
  EtcTKMysql      = "mysql"
  EtcTKPostgres   = "postgres"
  EtcTKClickHouse = "clickhouse"
  EtcTKRabbitmq   = "rabbitmq"
  EtcTKNeo4j      = "neo4j"
  EtcTKMgo        = "mongo"
  EtcTKRedis      = "redis"
  EtcTKPia        = "pika"
  EtcTKNsq        = "nsq"
  EtcTKKafka      = "kafka"
  EtcTKEs         = "es"
  EtcTKEtcd       = "etcd"

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
    "localhost:2379", //如果本机联调，想用测试机的etcd可以使用公网ip
  }
  QaEtcHosts = []string{ //QA环境ETCD地址，同正式
    "localhost:2379",
  }
  ProEtcHosts = []string{ //生产环境ETCD地址，需要使用内部dns解析，在k8s的worker节点配置/etc/hosts下面的域名和真实的etcd的ip
    "localhost:2379",
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
  Cluster     int    `json:"cluster,omitempty"` // 是否是集群 用于redis
  Vhost       string `json:"vhost,omitempty"`   // 虚拟目录用于rabbitmq
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
  ClickHouse   []LabelDetail                `json:"clickhouse,omitempty"`
  Rabbitmq     []LabelDetail                `json:"rabbitmq,omitempty"`
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

func InitConfig(cpath, env, project, etcdHosts string) ([]*mvccpb.KeyValue, chan map[string][]*ConnDetail, chan map[string]*CacheConfig, chan map[string][]*ConnDetail, chan map[string]*CacheConfig) {

  LoadConfig(cpath, env, project, etcdHosts)

  if env != Local && env != Container {
    //用etcd的配置
    ec := EtcConfig{
      Key:       fmt.Sprintf("%s/%s", ProjectPrefix, project),
      Endpoints: Conf.EtcdHosts,
    }
    return ec.InitConfigByEtcd()
  }
  return nil, nil, nil, nil, nil
}

func LoadConfig(cpath, env, project, etcdHosts string) {
  var cf string
  switch env {
  case Container:
    cf = fmt.Sprintf("%s/%s.json", cpath, env)
    bf, err := ioutil.ReadFile(cf)
    if err != nil {
      panic(err)
    }

    Conf = &allConfig{}
    err = zgoutils.Utils.Unmarshal(bf, Conf)
    if err != nil {
      panic(err)
    }
  case Local:
    cf = fmt.Sprintf("%s/%s.json", cpath, env)
    bf, err := ioutil.ReadFile(cf)
    if err != nil {
      panic(err)
    }

    Conf = &allConfig{}
    err = zgoutils.Utils.Unmarshal(bf, Conf)
    if err != nil {
      panic(err)
    }

  case Dev:
    Conf = &allConfig{
      Env:       env,
      Project:   project,
      EtcdHosts: DevEtcHosts,
      File: FileStore{
        Type: FileStoreType,
        Home: FileStoreHome,
      },
    }

  case Qa:
    Conf = &allConfig{
      Env:       env,
      Project:   project,
      EtcdHosts: QaEtcHosts,
      File: FileStore{
        Type: FileStoreType,
        Home: FileStoreHome,
      },
    }

  case Pro:
    Conf = &allConfig{
      Env:       env,
      Project:   project,
      EtcdHosts: ProEtcHosts,
      File: FileStore{
        Type: FileStoreType, //以后生产环境可以存到aws s3，在这里直接更改
        Home: FileStoreHome,
      },
    }
  case K8s:
    Conf = &allConfig{
      Env:       env,
      Project:   project,
      EtcdHosts: ProEtcHosts,
      File: FileStore{
        Type: FileStoreType, //以后生产环境可以存到aws s3，在这里直接更改
        Home: FileStoreHome,
      },
    }
  }

  if etcdHosts != "" {
    if strings.Contains(etcdHosts,",") {
      Conf.EtcdHosts = strings.Split(etcdHosts, ",")
    }else{
      Conf.EtcdHosts = []string{etcdHosts}
    }
  }

  //default init city db config
  Conf.CityDbConfig = cityDbConfig

  info := fmt.Sprintf("zgo engine %s is started on the ... %s %s\n", Version, Conf.Env, Conf.EtcdHosts)

  welcome(info)
}

func welcome(info string) {
  welcome :=
    `______________________________
\                             \           _         ______ |
 \                             \        /   \___-=O'/|O'/__|
  \  zgo engine, Here we go !!  \_______\          / | /    )
  /                             /        '/-==__ _/__|/__=-|  -GM
 /          @2019~2022              /         *             \ | |
/                             /                        (o)
------------------------------
`
  startMsg := "if you have any problem, please visit https://github.com/gitcpu-io/zgo"
  fmt.Printf("\n%s\n%s\n%s\n", welcome, info, startMsg)
}
