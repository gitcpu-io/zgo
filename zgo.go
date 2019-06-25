package zgo

import (
	"git.zhugefang.com/gocore/zgo/config"
	"git.zhugefang.com/gocore/zgo/limiter"
	"git.zhugefang.com/gocore/zgo/zgocache"
	"git.zhugefang.com/gocore/zgo/zgocrypto"
	"git.zhugefang.com/gocore/zgo/zgoes"
	"git.zhugefang.com/gocore/zgo/zgoetcd"
	"git.zhugefang.com/gocore/zgo/zgofile"
	"git.zhugefang.com/gocore/zgo/zgogrpc"
	"git.zhugefang.com/gocore/zgo/zgohttp"
	"git.zhugefang.com/gocore/zgo/zgokafka"
	"git.zhugefang.com/gocore/zgo/zgolog"
	"git.zhugefang.com/gocore/zgo/zgomap"
	"git.zhugefang.com/gocore/zgo/zgomongo"
	"git.zhugefang.com/gocore/zgo/zgomysql"
	"git.zhugefang.com/gocore/zgo/zgonsq"
	"git.zhugefang.com/gocore/zgo/zgopika"
	"git.zhugefang.com/gocore/zgo/zgopostgres"
	"git.zhugefang.com/gocore/zgo/zgoredis"
	"git.zhugefang.com/gocore/zgo/zgoutils"
	kafkaCluter "github.com/bsm/sarama-cluster"
	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
	"github.com/nsqio/go-nsq"
	"go.etcd.io/etcd/clientv3"
)

type engine struct {
	opt *Options
}

//New init zgo engine
func Engine(opt *Options) error {
	engine := &engine{
		opt: opt,
	}

	err := opt.Init() //把zgo_start中用户定义的，映射到zgo的内存变量上

	if err != nil {
		return err
	}

	Crypto = zgocrypto.New()
	File = zgofile.New()

	//初始化GRPC
	Grpc = zgogrpc.GetGrpc()

	Log = zgolog.InitLog(config.Conf.Project)
	//start 日志watch
	zgolog.StartLogStoreWatcher()

	//异步start 日志消费存储协程
	zgolog.LogStore = zgolog.NewLogStore()

	// 从local初始化缓存模块
	in := <-zgocache.InitCache()
	Cache = in

	go zgolog.LogStore.StartQueue()

	if opt.Env == config.Local {
		if len(opt.Mongo) > 0 {
			//todo someting
			hsm := engine.getConfigByOption(config.Conf.Mongo, opt.Mongo)
			//fmt.Println("--zgo.go--",config.Mongo, opt.Mongo, hsm)
			in := <-zgomongo.InitMongo(hsm)
			Mongo = in
		}

		if len(opt.Mysql) > 0 {
			//todo someting
			hsm := engine.getConfigByOption(config.Conf.Mysql, opt.Mysql)
			//fmt.Println(hsm)
			in := <-zgomysql.InitMysql(hsm)
			Mysql = in
		}
		if len(opt.Postgres) > 0 {
			//todo someting
			hsm := engine.getConfigByOption(config.Conf.Postgres, opt.Postgres)
			//fmt.Println(hsm)
			in := <-zgopostgres.InitPostgres(hsm)
			Postgres = in
		}
		//if len(opt.Neo4j) > 0 {
		//	//todo someting
		//	hsm := engine.getConfigByOption(config.Conf.Neo4j, opt.Neo4j)
		//	//fmt.Println(hsm)
		//	in := <-zgoneo4j.InitNeo4j(hsm)
		//	Neo4j = in
		//}
		if len(opt.Etcd) > 0 {
			//todo someting
			hsm := engine.getConfigByOption(config.Conf.Etcd, opt.Etcd)
			//fmt.Println(hsm)
			in := <-zgoetcd.InitEtcd(hsm)
			Etcd = in
		}
		if len(opt.Es) > 0 {
			hsm := engine.getConfigByOption(config.Conf.Es, opt.Es)
			in := <-zgoes.InitEs(hsm)
			Es = in
		}
		if len(opt.Redis) > 0 {
			//todo someting
			hsm := engine.getConfigByOption(config.Conf.Redis, opt.Redis)
			//fmt.Println(hsm)
			in := <-zgoredis.InitRedis(hsm)
			Redis = in
		}
		if len(opt.Pika) > 0 {
			//todo someting
			hsm := engine.getConfigByOption(config.Conf.Pika, opt.Pika)
			//fmt.Println(hsm)
			in := <-zgopika.InitPika(hsm)
			Pika = in
		}
		if len(opt.Nsq) > 0 { //>0表示用户要求使用nsq
			hsm := engine.getConfigByOption(config.Conf.Nsq, opt.Nsq)
			//fmt.Println("===zgo.go==", hsm)
			//return nil
			in := <-zgonsq.InitNsq(hsm)
			Nsq = in
		}
		if len(opt.Kafka) > 0 {
			//todo someting
			hsm := engine.getConfigByOption(config.Conf.Kafka, opt.Kafka)
			//fmt.Println(hsm)
			//return nil
			in := <-zgokafka.InitKafka(hsm)
			Kafka = in
		}

		cc := &config.CacheConfig{
			DbType: config.Conf.Log.DbType,
			Label:  config.Conf.Log.Label,
			Start:  config.Conf.Log.Start,
		}
		zgolog.LogWatch <- cc

		if opt.Loglevel != "" {
			ll := 0
			for k, v := range config.Levels {
				if v == opt.Loglevel {
					ll = k
					break
				}
			}
			config.Conf.Log.LogLevel = ll
		} else {
			config.Conf.Log.LogLevel = config.Debug
		}

	}

	return nil
}

//getConfigByOption 把zgo_start中的[]和config中的map进行match并取到关系
func (e *engine) getConfigByOption(lds []config.LabelDetail, us []string) map[string][]*config.ConnDetail {
	m := make(map[string][]*config.ConnDetail)
	for _, label := range us {
		//v == label_bj 用户传来的label，它并不知道具体的连接地址
		//v == label_sh 用户传来的label，它并不知道具体的连接地址
		for _, v := range lds {
			if label == v.Key {
				var tmp []*config.ConnDetail
				for k, _ := range v.Values {
					tmp = append(tmp, &v.Values[k])
				}
				m[v.Key] = tmp
			}
		}
	}
	return m
}

//定义外部使用的类型
type (
	NsqMessage        = *nsq.Message
	PartitionConsumer = kafkaCluter.PartitionConsumer
	Bucketer          = limiter.SimpleBucketer //zgo 自定义的bucket

	//postgres声明给使用者
	PostgresDB                 = pg.DB
	PostgresCreateTableOptions = orm.CreateTableOptions

	//neo4j声明给使用者
	//Neo4jSession     = neo4j.Session
	//Neo4jTransaction = neo4j.Transaction
	//Neo4jResult      = neo4j.Result

	//etcd声明给使用者
	EtcdClientV3    = clientv3.Client
	EtcdGetResponse = clientv3.GetResponse
)

var (
	Kafka             zgokafka.Kafkaer
	Nsq               zgonsq.Nsqer
	Mongo             zgomongo.Mongoer
	Es                zgoes.Eser
	Grpc              zgogrpc.Grpcer
	Redis             zgoredis.Rediser
	Pika              zgopika.Pikaer
	Mysql             zgomysql.Mysqler
	Postgres          zgopostgres.Postgreser
	PostgresErrNoRows = pg.ErrNoRows

	//Neo4j    zgoneo4j.Neo4jer
	Etcd  zgoetcd.Etcder
	Cache zgocache.Cacher

	Http = zgohttp.New()

	Log    zgolog.Logger
	Utils  = zgoutils.New()
	Crypto zgocrypto.Cryptoer
	Map    = zgomap.GetMap()
	File   zgofile.Filer

	Limiter = limiter.New()
)
