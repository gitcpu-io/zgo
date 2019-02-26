package zgomysql

import (
	"fmt"
	"git.zhugefang.com/gocore/zgo.git/config"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"log"
)

// 连接池 存放的map对象
var (
	connPoolMap = make(map[string]*gorm.DB)
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
		for i := 0; i < len(value); i++ {
			fmt.Println("initConnPool")
			db, err := gorm.Open("mysql", value[i].Host)
			if err != nil {
				// 	链接mysql异常时 打印并退出系统
				log.Fatalf(err.Error())
			}
			// 开发模式打开日志
			//if config.ServerConfig.Env == DevelopmentMode {
			db.LogMode(true)
			//}
			// 最大空闲连接 5
			db.DB().SetMaxIdleConns(value[i].MaxIdleSize)
			// 最大打开链接 50
			db.DB().SetMaxOpenConns(value[i].MaxOpenConn)

			// 禁用复数表名
			db.SingularTable(true)
			// todo 如何区分mysql不同库？多从库模式下
			connPoolMap[key] = db
		}

	}
}
