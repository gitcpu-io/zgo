package zgo

import (
	"git.zhugefang.com/gocore/zgo/config"
	"git.zhugefang.com/gocore/zgo/zgocache"
	"git.zhugefang.com/gocore/zgo/zgoclickhouse"
	"git.zhugefang.com/gocore/zgo/zgocrypto"
	"git.zhugefang.com/gocore/zgo/zgoes"
	"git.zhugefang.com/gocore/zgo/zgoetcd"
	"git.zhugefang.com/gocore/zgo/zgofile"
	"git.zhugefang.com/gocore/zgo/zgogrpc"
	"git.zhugefang.com/gocore/zgo/zgohttp"
	"git.zhugefang.com/gocore/zgo/zgokafka"
	"git.zhugefang.com/gocore/zgo/zgolb"
	"git.zhugefang.com/gocore/zgo/zgolimiter"
	"git.zhugefang.com/gocore/zgo/zgolog"
	"git.zhugefang.com/gocore/zgo/zgomap"
	"git.zhugefang.com/gocore/zgo/zgomgo"
	"git.zhugefang.com/gocore/zgo/zgomongo"
	"git.zhugefang.com/gocore/zgo/zgomysql"
	"git.zhugefang.com/gocore/zgo/zgonsq"
	"git.zhugefang.com/gocore/zgo/zgopika"
	"git.zhugefang.com/gocore/zgo/zgopostgres"
	"git.zhugefang.com/gocore/zgo/zgorabbitmq"
	"git.zhugefang.com/gocore/zgo/zgoredis"
	"git.zhugefang.com/gocore/zgo/zgoutils"
	kafkaCluter "github.com/bsm/sarama-cluster"
	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
	"github.com/nsqio/go-nsq"
	"github.com/streadway/amqp"
	"go.etcd.io/etcd/clientv3"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
		if len(opt.Mgo) > 0 {
			//todo someting
			hsm := engine.getConfigByOption(config.Conf.Mgo, opt.Mgo)
			//fmt.Println("--zgo.go--",config.Conf.Mgo, opt.Mgo, hsm)
			in := <-zgomgo.InitMgo(hsm)
			Mgo = in
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
		if len(opt.ClickHouse) > 0 {
			//todo someting
			hsm := engine.getConfigByOption(config.Conf.ClickHouse, opt.ClickHouse)
			//fmt.Println(hsm)
			in := <-zgoclickhouse.InitClickHouse(hsm)
			CK = in
		}
		if len(opt.Rabbitmq) > 0 {
			//todo someting
			hsm := engine.getConfigByOption(config.Conf.Rabbitmq, opt.Rabbitmq)
			//fmt.Println(hsm)
			in := <-zgorabbitmq.InitRabbitmq(hsm)
			MQ = in
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
	RabbitmqPublishing = amqp.Publishing
	NsqMessage         = *nsq.Message
	PartitionConsumer  = kafkaCluter.PartitionConsumer
	Bucketer           = zgolimiter.SimpleBucketer //zgo 自定义的bucket
	WR2er              = zgolb.WR2er
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

	MgoObjectId           = primitive.ObjectID           //mongo bson id
	MgoBulkWriteOperation = zgomgo.MgoBulkWriteOperation //多个并行计算
	MgoArgs               = zgomgo.MgoArgs               //CRUD->mongodb时的传入参数，具体参数由以下选择，>>>>>请使用前详细阅读>>>>>
	/*	Document     []interface{}    //保存时用到的结构体的指针
		Result interface{}            //接受结构体的指针 比如: r := &User{} 这里的result就是r
		Filter map[string]interface{} //查询条件
		ArrayFilters []map[string]interface{} //子文档的查询条件
		Fields map[string]interface{} //字段筛选，形如SQL中的select选择字段
		Update map[string]interface{} //更新项 或 替换项
		Sort map[string]interface{}   //排序 1是升序，-1是降序
		Limit int64                   //限制数量
		Skip int64                    //查询的offset，开区间，不包括这个skip对应的值
		Upsert bool                   //当查询不到时，true表示插入一条新的
	*/

	/**
	########################其中Filter构造如下########################
	filter = make(map[string]interface{})
	//filter["_id"] = "5d81e00bada5f1088cb1d236"
	filter["username"] = "朱大仙儿"	//可以是某字段或_id
	filter["houses"] = map[string]interface{}{
		"$gte": 130,	//可以是其它$or、$not、$lt
	}

	########################其中ArrayFilters构造如下########################
	var arrayFilters []map[string]interface{}
	af := make(map[string]interface{})
	//af["element"] = map[string]interface{}{
	//	//这里的element对应下面update中的houses.$[element]，意思是数组中的每一项元素
	//	"$gte": 134,
	//}
	af["elem.grade"] = map[string]interface{}{	//elem.grade和element二选一
		"$gte": 70,
	}
	af["elem.mean"] = map[string]interface{}{	//但elem中可以有多个elem.xx或elem.yy
		"$gte": 60,
	}
	arrayFilters = append(arrayFilters, af)

	########################其中Update构造如下########################
	update := make(map[string]interface{})
	update["$inc"] = map[string]interface{}{	******$inc******
		"age": 100,
		"money": -100,
		//可以有多个字段k,v;
	}
	update["$set"] = map[string]interface{}{	******$set******
		"address": "FindOneAndUpdate更新",
		"post": "100002",	//更新某字段
		//"houses.$[element]": 411001, //如果houses是纯数组:[xx,xx,xx]
		//子文档的$[element] 其中这个element可以自定义名字
		"grades.$[elem].mean": 100, //如果grades是对象数组:[{k:v,mean:v},{k:v,mean:v}]
		//子文档$[elem]
		//可以有多个字段k,v;但只能有一个顶级字段，意味着$[element]和$[ele]二选一
	}
	type Score struct {
		Grade int `json:"grade"`
		Mean int `json:"mean"`
	}
	update["$push"] = map[string]interface{}{	******$push******
		"scores": Score{	//已有一个数组，这里是一个个的push object对象进数组中
			Grade: 70,
			Mean:65,
		},
	}

	########################其中Sort构造如下########################
	sort := make(map[string]interface{})
	sort["_id"] = 1 //1升序
	sort["age] = -1	//-1降序

	########################其中Fields构造如下########################
	//如果返回错误：Projection cannot have a mix of inclusion and exclusion; //要么全是1，要么全是0
	fields := make(map[string]interface{})
	fields["age"] = 1 	//显示返回age字段
	fields["address"] = 1
	fields["username"] = 1
	*/

)

var (
	Kafka             zgokafka.Kafkaer
	Nsq               zgonsq.Nsqer
	Mongo             zgomongo.Mongoer
	Mgo               zgomgo.Mgoer
	Es                zgoes.Eser
	Grpc              zgogrpc.Grpcer
	Redis             zgoredis.Rediser
	Pika              zgopika.Pikaer
	Mysql             zgomysql.Mysqler
	Postgres          zgopostgres.Postgreser
	CK                zgoclickhouse.ClickHouseer
	MQ                zgorabbitmq.Rabbitmqer
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

	Limiter = zgolimiter.New()
	LB      = zgolb.NewLB()

	MgoBulkWriteOperation_InsertOne  = config.InsertOne
	MgoBulkWriteOperation_UpdateOne  = config.UpdateOne
	MgoBulkWriteOperation_ReplaceOne = config.ReplaceOne
	MgoBulkWriteOperation_DeleteOne  = config.DeleteOne
	MgoBulkWriteOperation_UpdateMany = config.UpdateMany
	MgoBulkWriteOperation_DeleteMany = config.DeleteMany
)
