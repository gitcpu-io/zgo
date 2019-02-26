package zgo_db_mysql

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"log"
	"time"
)

// 数据库链接存储对象
var dbs = make(map[string]*gorm.DB)

// 模拟数据库配置
var addr = map[string]string{
	"sell":    "root:123456@/spider?charset=utf8&parseTime=True&loc=Local",
	"complex": "root:123456@/spider?charset=utf8&parseTime=True&loc=Local",
}

// mysql连接池初始化方法
func init() {
	for key, value := range addr {
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
		dbs[key] = db
	}
}

// channel 方式线程池
func MysqlClientChan() chan *gorm.DB {
	return MysqlChan
}

const limitConn = 50

var (
	MysqlChan chan *gorm.DB
	mysqlDb   []*gorm.DB
)

func init1() {
	MysqlChan = make(chan *gorm.DB, 10000)
	mysqlDb = []*gorm.DB{}

	// 数据库配置
	addr := []string{
		"127.0.0.1:27017",
	}

	//每个host:port连接创建50个连接，放入slice中
	ssChanChan := make(chan chan *gorm.DB, limitConn*len(addr))

	go func() {
		for sessionCh := range ssChanChan {
			if session, ok := <-sessionCh; ok {
				mysqlDb = append(mysqlDb, session)
			}
		}
	}()

	for i := 0; i < limitConn; i++ {
		for _, host := range addr {
			ssChanChan <- createConnection(host)
		}
	}

	go func() {
		for {
			if len(MysqlChan) < 10000 {
				for _, s := range mysqlDb {
					if s != nil {
						MysqlChan <- s
					}
				}
			}
			time.Sleep(limitConn * time.Millisecond)
			//fmt.Println(len(MysqlChan), "--MysqlChan--")
		}

	}()
	go func() {
		time.Sleep(3 * time.Second)
		fmt.Println("init Mysql connection to MysqlChan ...", len(MysqlChan))
	}()

}

func createConnection(host string) chan *gorm.DB {
	out := make(chan *gorm.DB)

	go func() {
		db, err := gorm.Open("mysql", "root:123456@127.0.0.1/spider?charset=utf8&parseTime=True&loc=Local")
		if err != nil || db == nil {
			fmt.Println(db, err)
			out <- nil
			return
		}
		out <- db
	}()
	return out

}
