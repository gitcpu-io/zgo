package zgomysql

import (
	"git.zhugefang.com/gocore/zgo.git/config"
	"github.com/jinzhu/gorm"
	"log"
)

// 连接池 存放的map对象
var (
	connPoolMap map[string]*gorm.DB
)

//对外的接口
type ConnPooler interface {
	GetConn(label string) *gorm.DB
}

// ConnPooler实现类 与 方法
type connPool struct {
	label string
}

func (cp *connPool) GetConn(label string) *gorm.DB {
	return connPoolMap[label]
}

func NewConnPool(label string) *connPool {
	return &connPool{
		label: label,
	}
}

func GetPool(label string) *gorm.DB {
	return connPoolMap[label]
}

// 初始化连接池
func InitConnPool(hsm map[string][]config.ConnDetail) {
	for key, value := range hsm {
		db, err := gorm.Open("mysql", value)
		if err != nil {
			// 	链接mysql异常时 打印并退出系统
			log.Fatalf(err.Error())
		}
		// 开发模式打开日志
		//if config.ServerConfig.Env == DevelopmentMode {
		db.LogMode(true)
		//}
		// 最大空闲连接 5
		db.DB().SetMaxIdleConns(2)
		// 最大打开链接 50
		db.DB().SetMaxOpenConns(2)

		// 禁用复数表名
		db.SingularTable(true)
		connPoolMap[key] = db
	}
}
