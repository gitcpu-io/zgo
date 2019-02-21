package zgo_mongo

import (
	"context"
	"git.zhugefang.com/gocore/zgo.git/dbs/zgo_db_mongo"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

type Mongo struct {
}

func NewMongo() *Mongo {
	return &Mongo{}
}

func (m *Mongo) GetClientChan(name string) chan *mgo.Session {
	//name用来查找对应的库
	return zgo_db_mongo.MongoClientChan()
}

func (m *Mongo) Create(ctx context.Context, ch chan *mgo.Session, args map[string]interface{}) (interface{}, error) {
	return zgo_db_mongo.Create(ch, args["db"].(string), args["collection"].(string), args["args"])
}

func (m *Mongo) Update(ctx context.Context, ch chan *mgo.Session, args map[string]interface{}) (interface{}, error) {

	return nil, zgo_db_mongo.UpdateOne(ch, args)
}

func (m *Mongo) Delete(ctx context.Context, ch chan *mgo.Session, key string, args map[string]interface{}) (interface{}, error) {

	return nil, zgo_db_mongo.DeleteOne(ch, args["db"].(string), args["collection"].(string),
		args["query"].(bson.M))
}

func (m *Mongo) List(ctx context.Context, ch chan *mgo.Session, args map[string]interface{}) ([]interface{}, error) {
	//sort := args["sort"]
	if args["from"] == nil {
		args["from"] = 0
	}
	if args["size"] == nil {
		args["size"] = 10
	}
	if args["limit"] == nil {
		args["limit"] = 10
	}
	if args["sort"] == nil{
		args["sort"] = []string{}
	}
	return zgo_db_mongo.List(ch, args)
	//return zgo_db_mongo.List(ch, args["db"].(string), args["collection"].(string),
	//	args["query"].(bson.M))
}

func (m *Mongo) Get(ctx context.Context, ch chan *mgo.Session, args map[string]interface{}) (interface{}, error) {
	if args["select"] != nil {
		return zgo_db_mongo.GetBySelect(ch, args)
	}

	return zgo_db_mongo.Get(ctx, ch, args)
}


//type Monargs struct {
//	DB string `json:"db"`
//	Collection string `json:"col"`
//	Select string `json:"select"`
//}