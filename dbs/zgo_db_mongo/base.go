package zgo_db_mongo

import (
	"context"
	"errors"
	"fmt"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

type mongoResource struct {
	MongoChan chan *mgo.Session
}

func NewMongoResource() *mongoResource {
	return &mongoResource{
		MongoChan: MongoClientChan(),
	}
}

func Login(ch chan *mgo.Session, db, user, pass string) (*mgo.Session, error) {
	s := <-ch
	err := s.DB(db).Login(user, pass)
	return s, err
}

func Create(ch chan *mgo.Session,args map[string]interface{}) (interface{}, error) {
	s := <-ch
	return nil, s.DB(args["db"].(string)).C(args["collection"].(string)).Insert(args["items"])
}

// type bson.M map[string]interface{}
// test := bson.M{"name": "Jim"}
// test := bson.M{"name": bson.M{"first": "Jim"}}

//获取一条数据，期望参数db collection update， query参数不传则为全部查询
func Get(ctx context.Context, ch chan *mgo.Session, args map[string]interface{} )(interface{}, error) {
	//out := make(chan interface{})
	//go func(ctx context.Context) {
	//	//defer close(out)
	//	s := <-ch
	//	var res interface{}
	//
	//	s.DB(db).C(collection).Find(&query).One(&res)
	//	out <- res
	//}(ctx)
	s := <-ch
	var res interface{}
	//judge if not exist
	err := judgeMongo(args)
	if err != nil{
		return nil, err
	}
	preWorkForMongo(args)
	s.DB(args["db"].(string)).C(args["collection"].(string)).Find(args["query"].(bson.M)).Select(args["select"].(bson.M)).One(&res)
	//out <- res
	return res, nil

}

//修改多条数据，期望参数db collection update query参数不传则为全部查询 select 为期望返回字段，若不填则自动为空
func List(ch chan *mgo.Session, args map[string]interface{}) ([]interface{}, error) {
	judgeMongo(args)
	s := <-ch
	preWorkForMongo(args)
	ress := []interface{}{}
	err := s.DB(args["db"].(string)).C(args["collection"].(string)).Find(args["query"].(bson.M)).
		Select(args["select"].(bson.M)).Skip(args["from"].(int)).Limit(args["limit"].(int)).Sort(args["sort"].([]string)...).All(&ress)
	return ress, err
}

//修改多条数据，期望参数db collection update query参数不传则为最后一条更新
func UpdateOne(ch chan *mgo.Session, args map[string]interface{}) error {
	s := <-ch
	preWorkForMongo(args)
	return s.DB(args["db"].(string)).C(args["collection"].(string)).Update(args["query"].(bson.M), args["update"].(bson.M))
}

//修改多条数据，期望参数db collection update query参数不传则为全部查询
func UpdateAll(ch chan *mgo.Session, args map[string]interface{}) error {
	s := <-ch
	preWorkForMongo(args)
	_, err := s.DB(args["db"].(string)).C(args["collection"].(string)).UpdateAll(args["query"].(bson.M), args["update"].(bson.M))
	return err
}

func DeleteOne(ch chan *mgo.Session, args map[string]interface{}) error {
	err := judgeMongo(args)
	if err != nil{
		return err
	}
	s := <-ch
	if args["query"] == nil{
		return errors.New("query param can not empty when operate delete  ")
	}
	preWorkForMongo(args)
	return s.DB(args["db"].(string)).C(args["collection"].(string)).Remove(args["query"])
}

func DeleteAll(ch chan *mgo.Session, args map[string]interface{}) error {
	err := judgeMongo(args)
	if err != nil{
		return err
	}
	s := <-ch
	if args["query"] == nil{
		return errors.New("query param can not empty when operate delete  ")
	}
	preWorkForMongo(args)
	_, err = s.DB(args["db"].(string)).C(args["collection"].(string)).RemoveAll(args["query"])
	return err
}

func judgeMongo(args map[string]interface{})error{
	//switch args {
	//case args["db"] == nil:
	//	return errors.New("db is nil")
	//case args["collection"] == nil:
	//	return errors.New("collection is nil")
	//case args["query"] == nil:
	//	return errors.New("query is nil")
	//}
	if args["db"] == nil{
		return errors.New("db is nil")
	}else if  args["collection"] == nil{
		return errors.New("collection is nil")
	}else if  args["query"] == nil{
		return errors.New("query is nil")
	}
	return nil
}

//map转化成bson 内置函数
func map2Bson(arg map[string]interface{})bson.M{
	data, err := bson.Marshal(&arg)
	if err != nil{
		fmt.Println(err)
	}
	mmap := bson.M{}
	bson.Unmarshal(data, &mmap)
	return mmap
}

//map转换成bson 并且添加$set  避免修改导致对删除其他字段  内置函数
func updateMap2Bson(arg map[string]interface{})bson.M{
	set := map[string]map[string]interface{}{"$set": arg}
	data, err := bson.Marshal(&set)
	if err != nil{
		fmt.Println(err)
	}
	mmap := bson.M{}
	bson.Unmarshal(data, &mmap)
	return mmap
}

//对mongo操作对前置工作 讲map类型转化为bson 内置函数
func preWorkForMongo(args map[string]interface{})  {
	if args["query"] != nil{
		query := map2Bson(args["query"].(map[string]interface{}))
		args["query"] = query
	}else {
		args["query"] = map[string]interface{}{}
		query := map2Bson(args["query"].(map[string]interface{}))
		args["query"] = query
	}

	if args["update"] != nil{
		update := updateMap2Bson(args["update"].(map[string]interface{}))
		args["update"] = update
	}

	if args["select"] != nil{
		bySelect := map2Bson(args["select"].(map[string]interface{}))
		args["select"] = bySelect
	}else {
		args["select"] = map[string]interface{}{}
		bySelect := map2Bson(args["select"].(map[string]interface{}))
		args["select"] = bySelect
	}


}