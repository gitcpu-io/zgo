package zgomysql

import (
	"context"
	"git.zhugefang.com/gocore/zgo.git/config"
	"github.com/jinzhu/gorm"
	"sync"
)

var muLabel sync.Once

//Mongo 对外
type MysqlServiceInterface interface {
	//NewMysql(label string) (*zgoMysqlService, error)
	GetConnPool(label string) (*gorm.DB, error)
	//Create(ctx context.Context, args map[string]interface{}) (interface{}, error)
	//Update(ctx context.Context, args map[string]interface{}) (interface{}, error)
	//UpdateAll(ctx context.Context, args map[string]interface{}) (interface{}, error)
	//Delete(ctx context.Context, args map[string]interface{}) (interface{}, error)
	//List(ctx context.Context, args map[string]interface{}) ([]interface{}, error)
	Get(ctx context.Context, args map[string]interface{}) error
}

// 初始化
//InitMongo 初始化连接mongo
func InitMysqlService(hsm map[string][]*config.ConnDetail) {
	muLabel.Do(
		func() {
			InitMysqlResource(hsm)
		},
	)
}

// 对外接口
func MysqlService(l string) (MysqlServiceInterface, error) {
	return &zgoMysqlService{
		res: NewMysqlResourcer(l),
	}, nil
}

// 内部就结构体
type zgoMysqlService struct {
	res MysqlResourcerInterface //使用resource另外的一个接口
}

// 实现方法
func (c *zgoMysqlService) GetConnPool(label string) (*gorm.DB, error) {
	return c.res.GetPool(), nil
}
func (c *zgoMysqlService) Get(ctx context.Context, args map[string]interface{}) error {
	return c.res.Get(ctx, args)
}
