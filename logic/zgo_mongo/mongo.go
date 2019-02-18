package zgo_mongo

import (
	"context"
	"git.zhugefang.com/gocore/zgo.git/models/zgo_models_mongo"
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
	return zgo_models_mongo.MongoClientChan()
}

func (m *Mongo) Create(ctx context.Context, ch chan *mgo.Session, args map[string]interface{}) (interface{}, error) {
	return zgo_models_mongo.Create(ch, args["db"].(string), args["collection"].(string), args["args"])
}

func (m *Mongo) Update(ctx context.Context, ch chan *mgo.Session, key string, args map[string]interface{}) (interface{}, error) {

	return nil, zgo_models_mongo.UpdateOne(ch, args["db"].(string), args["collection"].(string),
		args["query"].(bson.M), args["update"].(bson.M))
}

func (m *Mongo) Delete(ctx context.Context, ch chan *mgo.Session, key string, args map[string]interface{}) (interface{}, error) {

	return nil, zgo_models_mongo.DeleteOne(ch, args["db"].(string), args["collection"].(string),
		args["query"].(bson.M))
}

func (m *Mongo) List(ctx context.Context, ch chan *mgo.Session, args map[string]interface{}) ([]interface{}, error) {
	sort := args["sort"]
	if args["from"] == nil {
		args["from"] = 0
	}
	if args["size"] == nil {
		args["size"] = 10
	}
	if args["limit"] != 0 {
		return zgo_models_mongo.ListByLimit(ch, args["db"].(string), args["collection"].(string),
			args["from"].(int), args["size"].(int), args["query"].(bson.M), args["select"].(bson.M), sort.([]string))
	}
	return zgo_models_mongo.List(ch, args["db"].(string), args["collection"].(string),
		args["query"].(bson.M))
}

func (m *Mongo) Get(ctx context.Context, ch chan *mgo.Session, args map[string]interface{}) (chan interface{}, error) {
	//if args["select"] != nil {
	//	return zgo_models_mongo.GetBySelect(ch,args["db"].(string), args["collection"].(string),
	//		args["query"].(bson.M), args["select"].(bson.M))
	//}
	return zgo_models_mongo.Get(ch, args["db"].(string), args["collection"].(string),
		args["query"].(bson.M))
}
