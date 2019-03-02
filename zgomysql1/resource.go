package zgomysql1

import (
	"context"
	"errors"
	"git.zhugefang.com/gocore/zgo.git/config"
	"github.com/jinzhu/gorm"
	"sync"
)

//NsqResourcer 给service使用
type MysqlResourcer interface {
	GetConnChan(label string) chan *gorm.DB
	//List(ctx context.Context, args map[string]interface{}) ([]interface{}, error)
	Get(ctx context.Context, args map[string]interface{}) error
	//Create(ctx context.Context, args map[string]interface{}) (interface{}, error)
	//UpdateOne(ctx context.Context, args map[string]interface{}) error
	//UpdateAll(ctx context.Context, args map[string]interface{}) error
	//DeleteOne(ctx context.Context, args map[string]interface{}) error
	//DeleteAll(ctx context.Context, args map[string]interface{}) error
}

//内部结构体
type mysqlResource struct {
	label    string
	mu       sync.RWMutex
	connpool ConnPooler
}

func NewMysqlResourcer(label string) MysqlResourcer {
	return &mysqlResource{
		label:    label,
		connpool: NewConnPool(label), //使用connpool
	}
}

func InitMysqlResource(hsm map[string][]*config.ConnDetail) {
	InitConnPool(hsm)
}

// type bson.M map[string]interface{}
// test := bson.M{"name": "Jim"}
// test := bson.M{"name": bson.M{"first": "Jim"}}

//GetConnChan 返回存放连接的chan
func (m *mysqlResource) GetConnChan(label string) chan *gorm.DB {
	return m.connpool.GetConnChan(label)
}

//获取一条数据，期望参数db collection update， query参数不传则为全部查询
func (m *mysqlResource) Get(ctx context.Context, args map[string]interface{}) error {
	s := <-m.connpool.GetConnChan(m.label)
	err := s.Table(args["tablename"].(string)).Where(args["query"], args["args"].([]interface{})...).First(args["out"]).Error
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
