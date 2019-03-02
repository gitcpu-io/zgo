package zgomysql

import (
	"fmt"
	"git.zhugefang.com/gocore/zgo/config"
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

func GetPool(label string) *gorm.DB {
	return connPoolMap[label]
}

// 初始化连接池
func InitConnPool(hsm map[string][]*config.ConnDetail) {
	for key, value := range hsm {
		c := &connPool{
			label: key,
		}
		c.setConnPoolToChan(value)
	}
}

func (cp *connPool) setConnPoolToChan(v []*config.ConnDetail) {
	for i := 0; i < len(v); i++ {
		pool, err := cp.createClient(v[i])
		if err != nil {
			fmt.Println("创建mysql链接池失败：", cp.label)
			log.Fatalf(err.Error())
		} else {
			connPoolMap[fmt.Sprintf("%s:%d", cp.label, i)] = pool
		}
	}
}

func (cp *connPool) createClient(v *config.ConnDetail) (*gorm.DB, error) {
	fmt.Println("initConnPool")
	db, err := gorm.Open("mysql", v.Host)
	if err != nil {
		// 	链接mysql异常时 打印并退出系统
		log.Fatalf(err.Error())
		return nil, err
	}
	// 开发模式打开日志
	//if config.ServerConfig.Env == DevelopmentMode {
	db.LogMode(true)
	//}
	// 最大空闲连接 5
	fmt.Println(v.MaxIdleSize)
	db.DB().SetMaxIdleConns(v.MaxIdleSize)
	// 最大打开链接 50
	fmt.Println(v.MaxOpenConn)
	db.DB().SetMaxOpenConns(v.MaxOpenConn)
	// 禁用复数表名
	db.SingularTable(true)
	return db, nil
}
