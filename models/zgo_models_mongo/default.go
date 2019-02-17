package zgo_models_mongo

import (
	"fmt"
	"github.com/globalsign/mgo"
)

var (
	Session *mgo.Session
	Mongo   MongoResourcer
)

// dialect "mongodb://test:123@127.0.0.1:27017"
func Init(addr string) error {
	var err error
	//dialInfo := &mgo.DialInfo{
	//	Addrs:[]string{addr, port},
	//	Username:user,
	//	Password:password,
	//	PoolLimit:4096,
	//}
	//fmt.Println(*dialInfo)
	Session, err = mgo.Dial(addr)
	//Session, err = mgo.DialWithInfo(dialInfo)
	if err != nil {
		return err
	}
	if Mongo == nil {
		Mongo, _ = NewMongoResource(Session)
		fmt.Println("set mgo successful")
	}
	return err
}

func init() {
	Init("mongodb://127.0.0.1:27017/local")
}
