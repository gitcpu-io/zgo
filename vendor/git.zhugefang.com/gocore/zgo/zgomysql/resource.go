package zgomysql

import (
	"context"
	"fmt"
	"git.zhugefang.com/gocore/zgo/config"
	"github.com/jinzhu/gorm"
)

// 初始化 连接池
func InitMysqlResource(hsm map[string][]*config.ConnDetail) {
	InitConnPool(hsm)
}

// 对外接口
type MysqlResourcerInterface interface {
	GetPool() *gorm.DB
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
	connpool *gorm.DB
}

// 对外函数 -- 创建mysqlResourcer对象
func NewMysqlResourcer(label string) MysqlResourcerInterface {
	return &mysqlResource{
		label:    label,
		connpool: GetPool(label), //使用connpool
	}
}

// mysqlResourcer 实现方法
func (mr *mysqlResource) GetPool() *gorm.DB {
	return mr.connpool
}

func (mr *mysqlResource) Get(ctx context.Context, args map[string]interface{}) error {
	gormpoll := mr.connpool
	fmt.Println("resource--Get")
	err := gormpoll.Table(args["tablename"].(string)).Where(args["query"], args["args"].([]interface{})...).First(args["out"]).Error
	return err
}
