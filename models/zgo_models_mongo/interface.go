package zgo_models_mongo

import "github.com/globalsign/mgo"

type MongoResourcer interface {
	MongoClient() *mgo.Session
}
