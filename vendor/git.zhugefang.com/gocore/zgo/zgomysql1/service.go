package zgomysql1

import (
	"context"
	"git.zhugefang.com/gocore/zgo/comm"
	"git.zhugefang.com/gocore/zgo/config"
	"github.com/jinzhu/gorm"
	"sync"
	"time"
)

var (
	currentLabels = make(map[string][]*config.ConnDetail)
	muLabel       sync.RWMutex
)

//Mysql 对外
type Mysqler interface {
	NewMysql(label ...string) (*zgomysql, error)
	GetConnChan(label ...string) (chan *gorm.DB, error)
	//Create(ctx context.Context, args map[string]interface{}) (interface{}, error)
	//Update(ctx context.Context, args map[string]interface{}) (interface{}, error)
	//UpdateAll(ctx context.Context, args map[string]interface{}) (interface{}, error)
	//Delete(ctx context.Context, args map[string]interface{}) (interface{}, error)
	//List(ctx context.Context, args map[string]interface{}) ([]interface{}, error)
	Get(ctx context.Context, args map[string]interface{}) error
}

func Mysql(l string) Mysqler {
	return &zgomysql{
		res: NewMysqlResourcer(l),
	}
}

//zgomong实现了Mysql的接口
type zgomysql struct {
	res MysqlResourcer //使用resource另外的一个接口
}

//InitMysql 初始化连接mysql
func InitMysql(hsm map[string][]*config.ConnDetail) {
	muLabel.Lock()
	defer muLabel.Unlock()

	currentLabels = hsm
	InitMysqlResource(hsm)
}

//GetMysql zgo内部获取一个连接mysql
func GetMysql(label ...string) (*zgomysql, error) {
	l, err := comm.GetCurrentLabel(label, muLabel, currentLabels)
	if err != nil {
		return nil, err
	}

	return &zgomysql{
		res: NewMysqlResourcer(l), //interface
	}, nil
}

func (n *zgomysql) NewMysql(label ...string) (*zgomysql, error) {
	return GetMysql(label...)
}

func (m *zgomysql) GetConnChan(label ...string) (chan *gorm.DB, error) {
	//label用来查找对应的库
	l, err := comm.GetCurrentLabel(label, muLabel, currentLabels)
	if err != nil {
		return nil, err
	}
	return m.res.GetConnChan(l), nil
}

func (m *zgomysql) Get(ctx context.Context, args map[string]interface{}) error {
	time.Sleep(10 * time.Second)
	return m.res.Get(ctx, args)
}
