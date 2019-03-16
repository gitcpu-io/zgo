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

// 基类 所有
type Base struct {
	Id int `json:"id"`
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
	//FindById(ctx context.Context, obj interface{}, id int) (int, error)
	//FindById(ctx context.Context, obj interface{}, id int) (int, error)
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
	c, err := GetPool(mr.label, "r")
	if err != nil {
		//zgo.Log.Error(err.Error())
		return GetPool(mr.label, "w")
	}
	return c, err
}

func (mr *mysqlResource) GetWPool() (*gorm.DB, error) {
	return GetPool(mr.label, "w")
}

func (mr *mysqlResource) Get(ctx context.Context, args map[string]interface{}) error {
	var (
		gormPool *gorm.DB
		err      error
	)
	if T, ok := args["T"]; ok {
		gormPool, err = mr.GetPool(T.(string))
	} else {
		gormPool, err = mr.GetRPool()
	}

	if err != nil {
		return err
	}
	err = gormPool.Table(args["table"].(string)).Where(args["query"], args["args"].([]interface{})...).First(args["obj"]).Error
	return err
}

func (mr *mysqlResource) List(ctx context.Context, args map[string]interface{}) error {
	var (
		gormPool *gorm.DB
		err      error
	)
	if T, ok := args["T"]; ok {
		gormPool, err = mr.GetPool(T.(string))
	} else {
		gormPool, err = mr.GetRPool()
	}

	if err != nil {
		return err
	}
	gormPool = gormPool.Table(args["table"].(string)).Where(args["query"], args["args"].([]interface{})...)
	currentLimit := 30
	if limit, ok := args["limit"]; ok {
		gormPool = gormPool.Limit(limit)
		currentLimit = limit.(int)
	} else {
		gormPool = gormPool.Limit(currentLimit)
	}
	if page, ok := args["page"]; ok {
		gormPool = gormPool.Offset((page.(int) - 1) * currentLimit)
	} else if offset, ok := args["offset"]; ok {
		gormPool = gormPool.Offset(offset)
	}
	if order, ok := args["order"]; ok {
		gormPool = gormPool.Order(order)
	}
	err = gormPool.Find(args["obj"]).Error
	return err
}

func (mr *mysqlResource) Count(ctx context.Context, args map[string]interface{}) error {
	var (
		gormPool *gorm.DB
		err      error
	)
	if T, ok := args["T"]; ok {
		gormPool, err = mr.GetPool(T.(string))
	} else {
		gormPool, err = mr.GetRPool()
	}

	if err != nil {
		return err
	}
	gormPool = gormPool.Table(args["table"].(string)).Where(args["query"], args["args"].([]interface{})...)
	err = gormPool.Count(args["count"]).Error
	return err
}

func (mr *mysqlResource) Create(ctx context.Context, args map[string]interface{}) error {
	gormPool, err := mr.GetWPool()
	if err != nil {
		return err
	}
	if gormPool.Table(args["table"].(string)).NewRecord(args["obj"]) {
		err = gormPool.Table(args["table"].(string)).Create(args["obj"]).Error
		return err
	} else {
		return errors.New("被创建对象不能有主键")
	}

}

func (mr *mysqlResource) UpdateOne(ctx context.Context, args map[string]interface{}) (int, error) {
	gormPool, err := mr.GetWPool()
	if err != nil {
		return 0, err
	}
	if v, ok := args["id"]; ok {
		if v.(int) > 0 {
			// args["data"] = map[string]interface{}{"name": "hello", "age": 18}
			db := gormPool.Table(args["table"].(string)).Where(" id = ? ", args["id"]).Updates(args["data"])
			count := db.RowsAffected
			err = db.Error
			return int(count), err
		}
	}
	return 0, errors.New("mysql updateOne method : id not allow null or 0")
}

func (mr *mysqlResource) DeleteOne(ctx context.Context, args map[string]interface{}) (int, error) {
	gormPool, err := mr.GetWPool()
	if err != nil {
		return 0, err
	}
	// 根据id删除
	if v, ok := args["id"]; ok {
		if v.(int) > 0 {
			if !gormPool.NewRecord(args["obj"]) {
				db := gormPool.Table(args["table"].(string)).Delete(args["obj"])
				count := db.RowsAffected
				err = db.Error
				return int(count), err
			}
		}
	}
	return 0, errors.New("mysql deleteOne method : id not allow null or 0")
}
