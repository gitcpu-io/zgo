package zgomysql

import (
	"errors"
	"fmt"
	"git.zhugefang.com/gocore/zgo/config"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"log"
	"math/rand"
)

// 连接池 存放的map对象
var (
	connPoolMap = make(map[string]map[string][]*gorm.DB)
)

// ConnPooler实现类 与 方法
type connPool struct {
	label string
}

func GetPool(label string, T string) (*gorm.DB, error) {
	if pss, ok := connPoolMap[T]; ok {
		if ps, ok1 := pss[label]; ok1 {
			if len(ps) == 0 {
				fmt.Println("错误的label：" + label)
				return nil, errors.New("错误的label：" + label)
			} else if len(ps) > 1 {
				index := rand.Intn(len(ps)) //随机取一个相同label下的连接
				return ps[index], nil
			} else {
				return ps[0], nil
			}
		}
	}
	return nil, errors.New("GetPool param T is wrong")

}

// InitConnPool 初始化连接池
/*
 1.
*/
func InitConnPool(hsm map[string][]*config.ConnDetail) {
	for key, value := range hsm {
		c := &connPool{
			label: key,
		}
		c.setConnPoolToChan(value)
	}
}

// setConnPoolToChan 创建链接池对象并放入全局Map内
/*
 1. v []*config.ConnDetail  配置信息对象集合
*/
func (cp *connPool) setConnPoolToChan(v []*config.ConnDetail) {
	for i := 0; i < len(v); i++ {
		pool, err := cp.createClient(v[i])
		if err != nil {
			fmt.Println("创建mysql链接池失败：", cp.label)
			log.Fatalf(err.Error())
		} else {
			key := fmt.Sprintf("%s", cp.label)
			fmt.Printf("init Mysql to Pool ... [%s] Host:%s, Port:%d, MaxOpenConn:%d, MaxIdleSize:%d, T:%s, LogMode:%d, %s\n",
				cp.label, v[i].Host, v[i].Port, v[i].MaxOpenConn, v[i].MaxIdleSize, v[i].T, v[i].LogMode, v[i].C)
			if value, ok := connPoolMap[v[i].T]; ok { // 是否能获取到2级Map
				value[key] = append(value[key], pool)
			} else { // 创建二级map
				// 创建slice
				pools := []*gorm.DB{pool}
				connPoolMap[v[i].T] = map[string][]*gorm.DB{key: pools}
			}
		}
	}
}

// createClient 创建链接池对象方法
/*
 1. v *config.ConnDetail  配置信息对象
*/
func (cp *connPool) createClient(v *config.ConnDetail) (*gorm.DB, error) {
	host := fmt.Sprintf("%v:%v@(%v:%v)/%v?charset=utf8mb4&parseTime=True&loc=Local", v.Username, v.Password, v.Host, v.Port, v.DbName)
	if v.C == "doris" {
		host = fmt.Sprintf("%v:%v@(%v:%v)/%v?charset=utf8&parseTime=True&loc=Local", v.Username, v.Password, v.Host, v.Port, v.DbName)
	}
	db, err := gorm.Open("mysql", host)
	if err != nil {
		// 	链接mysql异常时 打印并退出系统
		log.Fatalf(err.Error())
		return nil, err
	}
	// 开发模式打开日志
	//if config.ServerConfig.Env == DevelopmentMode {
	db.LogMode(v.LogMode == 1)
	//}
	// 最大空闲连接 5
	db.DB().SetMaxIdleConns(v.MaxIdleSize)
	// 最大打开链接 50
	db.DB().SetMaxOpenConns(v.MaxOpenConn)
	// 禁用复数表名
	db.SingularTable(true)
	return db, nil
}
