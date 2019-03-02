package zgomysql

import (
	"errors"
	"fmt"
	"git.zhugefang.com/gocore/zgo.git/config"
	"math/rand"
	"sync"
)

var currentLabels = make(map[string][]*config.ConnDetail)

var cityDbConfig = make(map[string]map[string]string)
var muLabel sync.Once

//Mongo 对外
type MysqlServiceInterface interface {
	NewRs(label string) (MysqlResourcerInterface, error)
	GetLabelByCity(city string, biz string, t string) (string, error)
	GetDbByCityBiz(city string, biz string) (string, error)
}

// 初始化
//InitMongo 初始化连接mongo
func InitMysqlService(hsm map[string][]*config.ConnDetail, cdc map[string]map[string]string) {
	muLabel.Do(
		func() {
			currentLabels = hsm
			cityDbConfig = cdc
			InitMysqlResource(hsm)
		},
	)
}

// 对外接口
func MysqlService() MysqlServiceInterface {
	return &zgoMysqlService{}
}

// 内部就结构体
type zgoMysqlService struct {
	res MysqlResourcerInterface //使用resource另外的一个接口
}

// 实现方法
func (c *zgoMysqlService) NewRs(label string) (MysqlResourcerInterface, error) {
	configs := currentLabels[label]
	if len(configs) == 0 {
		return nil, errors.New("错误的label：" + label)
	} else if len(configs) > 1 {
		index := rand.Intn(len(configs)) //随机取一个相同label下的连接
		return NewMysqlResourcer(label + ":" + string(index)), nil
	} else {
		return NewMysqlResourcer(label + ":0"), nil
	}
}

// 根据城市
func (c *zgoMysqlService) GetLabelByCity(city string, biz string, t string) (string, error) {
	if value, ok := cityDbConfig[biz]; ok {
		if value, ok := value[city]; ok {
			label := "mysql_" + biz + "_" + t + value
			return label, nil
		}
	}
	return "", errors.New(fmt.Sprintf("未知mysql label;city:%s;biz:%s;t:%s", city, biz, t))

}

// 根据城市和业务 获取dbname和实例label
func (c *zgoMysqlService) GetDbByCityBiz(city string, biz string) (string, error) {
	db := ""
	if biz == "sell" {
		if city == "bj" {
			db = "spider"
			return db, nil
		} else {
			db = "spider_" + city
			return db, nil
		}
	} else if biz == "newhouse" || biz == "rent" {
		return biz + "_" + city, nil
	} else if biz == "data" {
		return biz + "_" + city, nil
	}

	return db, errors.New("未知dbname biz:" + biz)
}
