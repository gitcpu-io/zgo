package zgomysql

import (
	"context"
	"errors"
	"git.zhugefang.com/gocore/zgo/config"
	"github.com/jinzhu/gorm"
)

// 初始化 连接池
func InitMysqlResource(hsm map[string][]*config.ConnDetail) {
	InitConnPool(hsm)
}

// 对外接口
type MysqlResourcerInterface interface {
	GetPool(t string) (*gorm.DB, error)
	GetRPool() (*gorm.DB, error)
	GetWPool() (*gorm.DB, error)
	List(ctx context.Context, args map[string]interface{}) error
	Count(ctx context.Context, args map[string]interface{}) error
	Get(ctx context.Context, args map[string]interface{}) error
	Create(ctx context.Context, args map[string]interface{}) error
	UpdateOne(ctx context.Context, args map[string]interface{}) (int, error)
	//UpdateAll(ctx context.Context, args map[string]interface{}) error
	DeleteOne(ctx context.Context, args map[string]interface{}) (int, error)
	//DeleteAll(ctx context.Context, args map[string]interface{}) error
}

//内部结构体
type mysqlResource struct {
	label string
	//connpool *gorm.DB
}

// 对外函数 -- 创建mysqlResourcer对象
func NewMysqlResourcer(label string) MysqlResourcerInterface {
	return &mysqlResource{
		label: label,
	}
}

// mysqlResourcer 实现方法
func (mr *mysqlResource) GetPool(t string) (*gorm.DB, error) {
	return GetPool(mr.label, t)
}

// mysqlResourcer 实现方法
func (mr *mysqlResource) GetRPool() (*gorm.DB, error) {
	return GetPool(mr.label, "r")
}

func (mr *mysqlResource) GetWPool() (*gorm.DB, error) {
	return GetPool(mr.label, "w")
}

func (mr *mysqlResource) Get(ctx context.Context, args map[string]interface{}) error {
	gormpoll, err := mr.GetRPool()
	if err != nil {
		return err
	}
	err = gormpoll.Table(args["tablename"].(string)).Where(args["query"], args["args"].([]interface{})...).First(args["obj"]).Error
	return err
}

func (mr *mysqlResource) List(ctx context.Context, args map[string]interface{}) error {
	gormpoll, err := mr.GetRPool()
	if err != nil {
		return err
	}
	gormpoll = gormpoll.Table(args["tablename"].(string)).Where(args["query"], args["args"].([]interface{})...)
	currentLimit := 30
	if limit, ok := args["limit"]; ok {
		gormpoll = gormpoll.Limit(limit)
		currentLimit = limit.(int)
	} else {
		gormpoll = gormpoll.Limit(currentLimit)
	}
	if page, ok := args["page"]; ok {
		gormpoll = gormpoll.Offset((page.(int) - 1) * currentLimit)
	} else if offset, ok := args["offset"]; ok {
		gormpoll = gormpoll.Offset(offset)
	}
	if order, ok := args["order"]; ok {
		gormpoll = gormpoll.Order(order)
	}
	err = gormpoll.Find(args["obj"]).Error
	return err
}

func (mr *mysqlResource) Count(ctx context.Context, args map[string]interface{}) error {
	gormpoll, err := mr.GetRPool()
	if err != nil {
		return err
	}
	gormpoll = gormpoll.Table(args["tablename"].(string))
	err = gormpoll.Count(args["count"]).Error
	return err
}

func (mr *mysqlResource) Create(ctx context.Context, args map[string]interface{}) error {
	gormpoll, err := mr.GetWPool()
	if err != nil {
		return err
	}
	if gormpoll.Table(args["tablename"].(string)).NewRecord(args["obj"]) {
		err = gormpoll.Table(args["tablename"].(string)).Create(args["obj"]).Error
		return err
	} else {
		return errors.New("被创建对象不能有主键")
	}

}

func (mr *mysqlResource) UpdateOne(ctx context.Context, args map[string]interface{}) (int, error) {
	gormpoll, err := mr.GetWPool()
	if err != nil {
		return 0, err
	}
	db := gormpoll.Table(args["tablename"].(string)).Where(" id = ? ", args["id"].(int)).Updates(args["data"])
	count := db.RowsAffected
	err = db.Error
	return int(count), err
}

func (mr *mysqlResource) DeleteOne(ctx context.Context, args map[string]interface{}) (int, error) {
	gormpoll, err := mr.GetWPool()
	if err != nil {
		return 0, err
	}
	db := gormpoll.Table(args["tablename"].(string)).Where(" id = ? ", args["id"].(int)).Delete(args["data"])
	count := db.RowsAffected
	err = db.Error
	return int(count), err
}
