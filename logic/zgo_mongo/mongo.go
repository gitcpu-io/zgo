package zgo_mongo

import (
	"github.com/globalsign/mgo/bson"
	"github.com/rubinus/zgo/models/zgo_models_mongo"
)

type Mongo struct {
	MongoResourcer zgo_models_mongo.MongoResourcer
}

func NewMongo() *Mongo {
	//mongoMgr := new(Mongo)
	//mongoMgr.MongoResourcer = zgo_models_mongo.Mongo
	//
	//return &mongoMgr
	return &Mongo{
		MongoResourcer: zgo_models_mongo.Mongo,
	}
}

func (mongoMgr Mongo) Create(args map[string]interface{}) (interface{}, error) {

	return zgo_models_mongo.Create(mongoMgr.MongoResourcer, args["db"].(string), args["collection"].(string), args["args"])
}

func (mongoMgr Mongo) Update(key string, args map[string]interface{}) (interface{}, error) {

	return nil, zgo_models_mongo.UpdateOne(mongoMgr.MongoResourcer, args["db"].(string), args["collection"].(string),
		args["query"].(bson.M), args["update"].(bson.M))
}

func (mongoMgr Mongo) Delete(key string, args map[string]interface{}) (interface{}, error) {

	return nil, zgo_models_mongo.DeleteOne(mongoMgr.MongoResourcer, args["db"].(string), args["collection"].(string),
		args["query"].(bson.M))
}

func (mongoMgr Mongo) List(args map[string]interface{}) ([]interface{}, error) {
	sort := args["sort"]
	if args["from"] == nil {
		args["from"] = 0
	}
	if args["size"] == nil {
		args["size"] = 10
	}
	if args["limit"] != 0 {
		return zgo_models_mongo.ListByLimit(mongoMgr.MongoResourcer, args["db"].(string), args["collection"].(string),
			args["from"].(int), args["size"].(int), args["query"].(bson.M), args["select"].(bson.M), sort.([]string))
	}
	return zgo_models_mongo.List(mongoMgr.MongoResourcer, args["db"].(string), args["collection"].(string),
		args["query"].(bson.M))
}

func (mongoMgr Mongo) Get(args map[string]interface{}) (interface{}, error) {
	if args["select"] != nil {
		return zgo_models_mongo.GetBySelect(mongoMgr.MongoResourcer, args["db"].(string), args["collection"].(string),
			args["query"].(bson.M), args["select"].(bson.M))
	}
	return zgo_models_mongo.Get(mongoMgr.MongoResourcer, args["db"].(string), args["collection"].(string),
		args["query"].(bson.M))
}
