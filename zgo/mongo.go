package zgo

import "github.com/rubinus/zgo/logic/zgo_mongo"

var Mongo mongoer

func init() {
	Mongo = zgo_mongo.NewMongo()
}

type mongoer interface {
	Get(map[string]interface{}) (interface{}, error)
	List(map[string]interface{}) ([]interface{}, error)
	Create(map[string]interface{}) (interface{}, error)
	Update(string, map[string]interface{}) (interface{}, error)
	Delete(string, map[string]interface{}) (interface{}, error)
}
