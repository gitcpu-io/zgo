package zgomysql

import (
	"context"
	"database/sql"
	"git.zhugefang.com/gocore/zgo.git/config"
	"github.com/jinzhu/gorm"
)

// 初始化 连接池
func InitMysqlResource(hsm map[string][]config.ConnDetail) {
	InitConnPool(hsm)
}

// 对外接口
type MysqlResourcerInterface interface {
	GetConn() *sql.DB
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
func (mr *mysqlResource) GetConn() *sql.DB {
	return mr.connpool.DB()
}

// mysqlResourcer 实现方法
func (mr *mysqlResource) GetPool() *gorm.DB {
	return mr.connpool
}

func (mr *mysqlResource) Get(ctx context.Context, args map[string]interface{}) error {
	gormpoll := mr.connpool
	err := gormpoll.Table(args["tablename"].(string)).Where(args["query"], args["args"].([]interface{})...).First(args["object"]).Error
	return err
}
