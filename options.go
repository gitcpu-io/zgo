package zgo

import (
  "errors"
  "fmt"
  "github.com/gitcpu-io/zgo/config"
  "github.com/gitcpu-io/zgo/zgocache"
  "github.com/gitcpu-io/zgo/zgoclickhouse"
  "github.com/gitcpu-io/zgo/zgoes"
  "github.com/gitcpu-io/zgo/zgoetcd"
  "github.com/gitcpu-io/zgo/zgokafka"
  "github.com/gitcpu-io/zgo/zgolog"
  "github.com/gitcpu-io/zgo/zgomgo"
  "github.com/gitcpu-io/zgo/zgomysql"
  "github.com/gitcpu-io/zgo/zgonsq"
  "github.com/gitcpu-io/zgo/zgopostgres"
  "github.com/gitcpu-io/zgo/zgorabbitmq"
  "github.com/gitcpu-io/zgo/zgoredis"
  "github.com/gitcpu-io/zgo/zgoutils"
  "go.etcd.io/etcd/api/v3/mvccpb"
  "strings"
)

type Options struct {
  CPath      string   `json:"cpath"`
  Env        string   `json:"env"`
  Project    string   `json:"project"`
  EtcdHosts  string   `json:"etcdHosts"`
  Loglevel   string   `json:"loglevel"`
  Mongo      []string `json:"mongo"`
  Mgo        []string `json:"mgo"`
  Mysql      []string `json:"mysql"`
  Postgres   []string `json:"postgres"`
  ClickHouse []string `json:"clickhouse"`
  Rabbitmq   []string `json:"rabbitmq"`
  Neo4j      []string `json:"neo4j"`
  Etcd       []string `json:"etcd"`
  Es         []string `json:"es"`
  Redis      []string `json:"redis"`
  Pika       []string `json:"pika"`
  Kafka      []string `json:"kafka"`
  Nsq        []string `json:"nsq"`
}

func (opt *Options) Init() error {
  //init config
  if opt.Env == "" {
    opt.Env = config.Local
  } else {
    if opt.Env != config.Local && opt.Env != config.Container && opt.Env != config.Dev && opt.Env != config.Qa && opt.Env != config.Pro && opt.Env != config.K8s {
      return errors.New("error env,must be local/dev/qa/pro/k8s/container !")
    }
    if opt.Project == "" {
      return errors.New("u must input your Project name to zgo.Engine func .")
    }
  }

  //如果connCh有值表示启用了etcd为配置中心，并watch了key，等待变更ing...
  resKvs, connCh, cacheLogCh, delConnCh, delCacheLogCh := config.InitConfig(opt.CPath, opt.Env, opt.Project, opt.EtcdHosts)

  //监听put资源组件
  opt.watchPutConn(connCh)
  //监听delete资源组件
  opt.watchDeleteConn(delConnCh)

  //监听put的cache和log操作
  opt.watchPutCacheOrLog(cacheLogCh)
  //监听删除cache和log操作
  opt.watchDeleteCacheOrLog(delCacheLogCh)

  //解析etcd中的配置
  opt.parseConfig(resKvs, connCh, cacheLogCh)

  return nil
}

// parseConfig 解析etcd中的配置
func (opt *Options) parseConfig(resKvs []*mvccpb.KeyValue, connCh chan map[string][]*config.ConnDetail, cacheLogCh chan map[string]*config.CacheConfig) {
  for _, v := range resKvs {
    go func(v *mvccpb.KeyValue) {
      key := string(v.Key)
      smk := strings.Split(key, "/")
      labelType := smk[3]
      b := v.Value
      if labelType == config.EtcTKCache || labelType == config.EtcTKLog { //如果cache or log配置
        var cm config.CacheConfig

        err := zgoutils.Utils.Unmarshal(b, &cm)
        if err != nil {
          fmt.Println("反序列化当前值失败", key)
        }
        var hsm = make(map[string]*config.CacheConfig)

        hsm[key] = &cm

        cacheLogCh <- hsm

      } else if smk[1] == "project" && smk[2] == opt.Project { //init conn config by etcd

        var m []config.ConnDetail
        err := zgoutils.Utils.Unmarshal(b, &m)
        if err != nil {
          fmt.Println("反序列化当前值失败", key)
        }

        label := smk[4]
        var hsm = make(map[string][]*config.ConnDetail)
        var tmp []*config.ConnDetail

        for _, vv := range m {
          pvv := vv
          tmp = append(tmp, &pvv)

          sb := strings.Builder{}
          sb.WriteString(fmt.Sprintf("\n********************************资源项: %s ********************************\n", labelType))
          sb.WriteString(fmt.Sprintf("描述: %s\n", pvv.C))
          sb.WriteString(fmt.Sprintf("Label: %s\n", label))
          sb.WriteString(fmt.Sprintf("Host: %s\n", pvv.Host))
          sb.WriteString(fmt.Sprintf("Port: %d\n", pvv.Port))
          if labelType == config.EtcTKMysql || labelType == config.EtcTKPostgres || labelType == config.EtcTKClickHouse || labelType == config.EtcTKMgo {
            sb.WriteString(fmt.Sprintf("DbName: %s\n", pvv.DbName))
          }
          if labelType == config.EtcTKRabbitmq {
            sb.WriteString(fmt.Sprintf("Vhost: %s", pvv.Vhost))
          }
          if labelType == config.EtcTKRedis {
            cluster := "单机"
            if pvv.Cluster == 1 {
              cluster = "集群"
            }
            sb.WriteString(fmt.Sprintf("模式: %s\n", cluster))
            sb.WriteString(fmt.Sprintf("Db: %d\n", pvv.Db))
          }
          fmt.Println(sb.String())
          //fmt.Printf("\n**********************资源项: %s **************************\n", labelType)
          //fmt.Printf("描述: %s\n", pvv.C)
          //fmt.Printf("Label: %s\n", label)
          //fmt.Printf("Host: %s\n", pvv.Host)
          //fmt.Printf("Port: %d\n", pvv.Port)
          //fmt.Printf("DbName: %s\n", pvv.DbName)
          //fmt.Printf("Db: %d\n", pvv.Db)
        }
        hsm[key] = tmp
        connCh <- hsm
      }
    }(v)
  }
}

// watchPutConn 监听保存到etcd中的资源key，连接类型
func (opt *Options) watchPutConn(inch chan map[string][]*config.ConnDetail) {
  go func() {
    if inch != nil {
      for h := range inch {
        //KEY: zgo/project/项目名/mysql/label名字
        for k := range h {
          smk := strings.Split(k, "/")
          labelType := smk[3]
          hsm := make(map[string][]*config.ConnDetail)
          var label string
          for k, v := range h {
            label = strings.Split(k, "/")[4] //改变label，去掉前缀
            hsm[label] = v
          }
          fmt.Printf("[init %s conn]watchPutConn: %s\n", labelType, label)
          //[init mongo conn]watchPutConn: 1607450184770

          go opt.initConn(labelType, hsm)
        }
      }
    }

  }()
}

// watchDeleteConn 监听从etcd中删除的资源key，连接类型
func (opt *Options) watchDeleteConn(ch chan map[string][]*config.ConnDetail) {
  go func() {
    if ch != nil {
      //KEY: zgo/project/项目名/mysql/label名字
      for h := range ch {
        for k := range h {
          smk := strings.Split(k, "/")
          labelType := smk[3]
          label := smk[4]
          fmt.Printf("[destroy %s conn]watchDeleteConn %s\n", labelType, label)
          //[destroy nsq conn]watchDeleteConn 1068052762090

          opt.destroyConn(labelType, label)
        }
      }
    }
  }()
}

/**
1.从当前的currentLabels这个map中删除掉key
2.call connpool的map，删除掉对应的key，让gc释放掉连接
*/
// destroyConn 具体删除操作
func (opt *Options) destroyConn(labelType, label string) {
  switch labelType {
  case config.EtcTKMysql:
    in := <-zgomysql.InitMysql(nil, label)
    Mysql = in
  case config.EtcTKPostgres:
    in := <-zgopostgres.InitPostgres(nil, label)
    Postgres = in
  case config.EtcTKClickHouse:
    in := <-zgoclickhouse.InitClickHouse(nil, label)
    CK = in
  case config.EtcTKRabbitmq:
    in := <-zgorabbitmq.InitRabbitmq(nil, label)
    MQ = in
  //case config.EtcTKNeo4j:
  //	in := <-zgoneo4j.InitNeo4j(nil, label)
  //	Neo4j = in
  case config.EtcTKMgo:
    in := <-zgomgo.InitMgo(nil, label)
    Mongo = in
  case config.EtcTKRedis:
    in := <-zgoredis.InitRedis(nil, label)
    Redis = in
  case config.EtcTKPia:
    in := <-zgoredis.InitRedis(nil, label)
    Pika = in
  case config.EtcTKNsq:
    in := <-zgonsq.InitNsq(nil, label)
    Nsq = in
  case config.EtcTKKafka:
    in := <-zgokafka.InitKafka(nil, label)
    Kafka = in
  case config.EtcTKEs:
    in := <-zgoes.InitEs(nil, label)
    Es = in
  case config.EtcTKEtcd:
    in := <-zgoetcd.InitEtcd(nil, label)
    Etcd = in
  }
}

// watchPutCacheOrLog 监听put cache和log的操作
func (opt *Options) watchPutCacheOrLog(cacheLogCh chan map[string]*config.CacheConfig) {
  go func() {
    if cacheLogCh != nil {
      for cm := range cacheLogCh {
        //KEY: zgo/project/项目名/log
        //KEY: zgo/project/项目名/cache

        for k, v := range cm {
          smk := strings.Split(k, "/")
          labelType := smk[3]

          config.Conf.Cache.Label = v.Label
          config.Conf.Cache.Rate = v.Rate
          config.Conf.Cache.Start = v.Start
          config.Conf.Cache.TcType = v.TcType
          config.Conf.Cache.DbType = v.DbType

          switch labelType {
          case config.EtcTKCache:

            // 从etcd初始化缓存模块
            in := zgocache.InitCacheByEtcd(v)
            Cache = <-in

            fmt.Println("[init Cache]watchPutCacheOrLog", labelType, v)
          case config.EtcTKLog:

            Log = zgolog.InitLog(opt.Project)

            config.Conf.Log.DbType = v.DbType
            config.Conf.Log.Label = v.Label
            config.Conf.Log.LogLevel = v.LogLevel
            config.Conf.Log.Start = v.Start

            cc := config.CacheConfig{
              DbType:   v.DbType,
              Label:    v.Label,
              Start:    v.Start,
              LogLevel: v.LogLevel,
            }

            fmt.Println("[init Log]watchPutCacheOrLog", labelType, cc)

            zgolog.LogWatch <- &cc
          }
        }

      }
    }

  }()
}

// watchDeleteCacheAndLog 监听删除etcd中的 cache和log类型的key
func (opt *Options) watchDeleteCacheOrLog(ch chan map[string]*config.CacheConfig) {
  go func() {
    if ch != nil {
      //KEY: zgo/project/项目名/mysql/label名字
      for h := range ch {
        for k, v := range h {
          labelType := strings.Split(k, "/")[3]
          fmt.Printf("[destroy %s]watchDeleteCacheAndLog: %v\n", labelType, v)
          //[destroy]watchDeleteCacheAndLog: log &{日志存储 0 /tmp 1 file 0}

          opt.destroyCacheAndLog(labelType, v)
        }
      }
    }
  }()
}

// destroyCacheAndLog 具体删除操作
func (opt *Options) destroyCacheAndLog(labelType string, cf *config.CacheConfig) {

  switch labelType {
  case config.EtcTKCache:
    //如果delete是cache todo something
    config.Conf.Cache.Label = cf.Label
    config.Conf.Cache.Rate = cf.Rate
    config.Conf.Cache.Start = 0 //停止
    config.Conf.Cache.TcType = cf.TcType
    config.Conf.Cache.DbType = cf.DbType

    in := <-zgocache.InitCache()
    Cache = in

  case config.EtcTKLog:
    //如果delete是log todo something
    config.Conf.Log.DbType = cf.DbType
    config.Conf.Log.Label = cf.Label
    config.Conf.Log.LogLevel = cf.LogLevel
    config.Conf.Log.Start = 0

    cc := &config.CacheConfig{
      DbType:   cf.DbType,
      Label:    cf.Label,
      LogLevel: cf.LogLevel,
      Start:    0,
    }
    zgolog.LogWatch <- cc
  }
}

// initConn具体的连接操作
func (opt *Options) initConn(labelType string, hsm map[string][]*config.ConnDetail) {
  switch labelType {
  case config.EtcTKMysql:
    //init mysql again
    if len(hsm) > 0 {
      in := <-zgomysql.InitMysql(hsm)
      Mysql = in
    }
  case config.EtcTKPostgres:
    //init postgres again
    if len(hsm) > 0 {
      in := <-zgopostgres.InitPostgres(hsm)
      Postgres = in
    }
  case config.EtcTKClickHouse:
    //init clickhouse again
    if len(hsm) > 0 {
      in := <-zgoclickhouse.InitClickHouse(hsm)
      CK = in
    }
  case config.EtcTKRabbitmq:
    //init rabbitmq again
    if len(hsm) > 0 {
      in := <-zgorabbitmq.InitRabbitmq(hsm)
      MQ = in
    }
  case config.EtcTKNeo4j:
    //init neo4j again
    //if len(hsm) > 0 {
    //	in := <-zgoneo4j.InitNeo4j(hsm)
    //	Neo4j = in
    //}

  case config.EtcTKMgo:
    //init mgo again
    if len(hsm) > 0 {
      in := <-zgomgo.InitMgo(hsm)
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
      in := <-zgoredis.InitRedis(hsm)
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

  case config.EtcTKEtcd:
    //init etcd again
    if len(hsm) > 0 {
      in := <-zgoetcd.InitEtcd(hsm)
      Etcd = in
    }
  }
}

//func getMatchConfig(lds map[string][]*config.ConnDetail, us []string) map[string][]*config.ConnDetail {
//  m := make(map[string][]*config.ConnDetail)
//  for _, label := range us {
//    //v == label_bj 用户传来的label，它并不知道具体的连接地址
//    //v == label_sh 用户传来的label，它并不知道具体的连接地址
//    for k, v := range lds {
//      if label == k {
//        m[label] = v
//      }
//    }
//  }
//  return m
//}
