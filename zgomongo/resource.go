package zgomongo

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

//NsqResourcer 给service使用
type MongoResourcer interface {
	GetConnChan(label string) chan *mgo.Session
	List(ctx context.Context, args map[string]interface{}) ([]interface{}, error)
	Get(ctx context.Context, args map[string]interface{}) (interface{}, error)
	Create(ctx context.Context, args map[string]interface{}) (interface{}, error)
	UpdateOne(ctx context.Context, args map[string]interface{}) error
	UpdateAll(ctx context.Context, args map[string]interface{}) error
	DeleteOne(ctx context.Context, args map[string]interface{}) error
	DeleteAll(ctx context.Context, args map[string]interface{}) error
}

//内部结构体
type mongoResource struct {
	label    string
	mu       sync.RWMutex
	connpool ConnPooler
}

func NewMongoResourcer(label string) MongoResourcer {
	return &mongoResource{
		label:    label,
		connpool: NewConnPool(label), //使用connpool
	}
}

func InitMongoResource(hsm map[string][]string) {
	InitConnPool(hsm)
}

//GetConnChan 返回存放连接的chan
func (m *mongoResource) GetConnChan(label string) chan *mgo.Session {
	return m.connpool.GetConnChan(label)
}

func (m *mongoResource) Login(ctx context.Context, db, user, pass string) (*mgo.Session, error) {
	s := <-m.connpool.GetConnChan(m.label)

	err := s.DB(db).Login(user, pass)
	return s, err
}

func (m *mongoResource) Create(ctx context.Context, args map[string]interface{}) (interface{}, error) {
	s := <-m.connpool.GetConnChan(m.label)
	return nil, s.DB(args["db"].(string)).C(args["collection"].(string)).Insert(args["items"])
}

// type bson.M map[string]interface{}
// test := bson.M{"name": "Jim"}
// test := bson.M{"name": bson.M{"first": "Jim"}}

//获取一条数据，期望参数db collection update， query参数不传则为全部查询
func (m *mongoResource) Get(ctx context.Context, args map[string]interface{}) (interface{}, error) {
	s := <-m.connpool.GetConnChan(m.label)
	var res interface{}
	//judge if not exist
	err := judgeMongo(args)
	if err != nil {
		return nil, err
	}
	preWorkForMongo(args)
	s.DB(args["db"].(string)).C(args["collection"].(string)).Find(args["query"].(bson.M)).Select(args["select"].(bson.M)).One(&res)
	return res, nil

}

//修改多条数据，期望参数db collection update query参数不传则为全部查询 select 为期望返回字段，若不填则自动为空
func (m *mongoResource) List(ctx context.Context, args map[string]interface{}) ([]interface{}, error) {
	judgeMongo(args)
	s := <-m.connpool.GetConnChan(m.label)
	preWorkForMongo(args)
	ress := []interface{}{}
	err := s.DB(args["db"].(string)).C(args["collection"].(string)).Find(args["query"].(bson.M)).
		Select(args["select"].(bson.M)).Skip(args["from"].(int)).Limit(args["limit"].(int)).Sort(args["sort"].([]string)...).All(&ress)
	return ress, err
}

//修改多条数据，期望参数db collection update query参数不传则为最后一条更新
func (m *mongoResource) UpdateOne(ctx context.Context, args map[string]interface{}) error {
	s := <-m.connpool.GetConnChan(m.label)
	preWorkForMongo(args)
	return s.DB(args["db"].(string)).C(args["collection"].(string)).Update(args["query"].(bson.M), args["update"].(bson.M))
}

//修改多条数据，期望参数db collection update query参数不传则为全部查询
func (m *mongoResource) UpdateAll(ctx context.Context, args map[string]interface{}) error {
	s := <-m.connpool.GetConnChan(m.label)
	preWorkForMongo(args)
	_, err := s.DB(args["db"].(string)).C(args["collection"].(string)).UpdateAll(args["query"].(bson.M), args["update"].(bson.M))
	return err
}

func (m *mongoResource) DeleteOne(ctx context.Context, args map[string]interface{}) error {
	err := judgeMongo(args)
	if err != nil {
		return err
	}
	s := <-m.connpool.GetConnChan(m.label)
	if args["query"] == nil {
		return errors.New("query param can not empty when operate delete  ")
	}
	preWorkForMongo(args)
	return s.DB(args["db"].(string)).C(args["collection"].(string)).Remove(args["query"])
}

func (m *mongoResource) DeleteAll(ctx context.Context, args map[string]interface{}) error {
	err := judgeMongo(args)
	if err != nil {
		return err
	}
	s := <-m.connpool.GetConnChan(m.label)
	if args["query"] == nil {
		return errors.New("query param can not empty when operate delete  ")
	}
	preWorkForMongo(args)
	_, err = s.DB(args["db"].(string)).C(args["collection"].(string)).RemoveAll(args["query"])
	return err
}

func judgeMongo(args map[string]interface{}) error {
	//switch args {
	//case args["db"] == nil:
	//	return errors.New("db is nil")
	//case args["collection"] == nil:
	//	return errors.New("collection is nil")
	//case args["query"] == nil:
	//	return errors.New("query is nil")
	//}
	if args["db"] == nil {
		return errors.New("db is nil")
	} else if args["collection"] == nil {
		return errors.New("collection is nil")
	} else if args["query"] == nil {
		return errors.New("query is nil")
	}
	return nil
}

//map转化成bson 内置函数
func map2Bson(arg map[string]interface{}) bson.M {
	data, err := bson.Marshal(&arg)
	if err != nil {
		fmt.Println(err)
	}
	mmap := bson.M{}
	bson.Unmarshal(data, &mmap)
	return mmap
}

//map转换成bson 并且添加$set  避免修改导致对删除其他字段  内置函数
func updateMap2Bson(arg map[string]interface{}) bson.M {
	set := map[string]map[string]interface{}{"$set": arg}
	data, err := bson.Marshal(&set)
	if err != nil {
		fmt.Println(err)
	}
	mmap := bson.M{}
	bson.Unmarshal(data, &mmap)
	return mmap
}

//对mongo操作对前置工作 讲map类型转化为bson 内置函数
func preWorkForMongo(args map[string]interface{}) {
	if args["query"] != nil {
		//query := map2Bson(args["query"].(map[string]interface{}))
		//args["query"] = query
	} else {
		args["query"] = map[string]interface{}{}
		query := map2Bson(args["query"].(map[string]interface{}))
		args["query"] = query
	}

	if args["update"] != nil {
		update := updateMap2Bson(args["update"].(map[string]interface{}))
		args["update"] = update
	}

	if args["select"] != nil {
		bySelect := map2Bson(args["select"].(map[string]interface{}))
		args["select"] = bySelect
	} else {
		args["select"] = map[string]interface{}{}
		bySelect := map2Bson(args["select"].(map[string]interface{}))
		args["select"] = bySelect
	}

}
