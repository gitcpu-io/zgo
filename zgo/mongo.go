package zgo

import (
	"context"
	"github.com/globalsign/mgo"
	"github.com/rubinus/zgo/logic/zgo_mongo"
)

var Mongo mongoer

func init() {
	Mongo = zgo_mongo.NewMongo()
}

type mongoer interface {
	GetClientChan(name string) chan *mgo.Session
	Get(ctx context.Context, session chan *mgo.Session, args map[string]interface{}) (chan interface{}, error) //返回chan

	List(ctx context.Context, session chan *mgo.Session, args map[string]interface{}) ([]interface{}, error)
	Create(ctx context.Context, session chan *mgo.Session, args map[string]interface{}) (interface{}, error)
	Update(ctx context.Context, session chan *mgo.Session, key string, args map[string]interface{}) (interface{}, error)
	Delete(ctx context.Context, session chan *mgo.Session, key string, args map[string]interface{}) (interface{}, error)
}
