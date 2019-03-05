package zgomysql

import (
	"context"
	"errors"
	"fmt"
	"git.zhugefang.com/gocore/zgo/config"
	"github.com/jinzhu/gorm"
	"sync"
)

var currentLabels = make(map[string][]*config.ConnDetail)

var cityDbConfig = make(map[string]map[string]string)
var muLabel sync.RWMutex

//Mongo 对外
type MysqlServiceInterface interface {
	//NewRs(label string) (MysqlResourcerInterface, error)
	GetPool(t string) (*gorm.DB, error)
	Get(ctx context.Context, args map[string]interface{}) error
	List(ctx context.Context, args map[string]interface{}) error
	Count(ctx context.Context, args map[string]interface{}) error
	Create(ctx context.Context, args map[string]interface{}) error
	UpdateOne(ctx context.Context, args map[string]interface{}) (int, error)
	DeleteOne(ctx context.Context, args map[string]interface{}) (int, error)
	GetLabelByCityBiz(city string, biz string) (string, error)
	GetDbByCityBiz(city string, biz string) (string, error)
	MysqlServiceByCityBiz(city string, biz string) (MysqlServiceInterface, error)
}

// 初始化
//InitMysql 初始化连接mongo
func InitMysqlService(hsm map[string][]*config.ConnDetail, cdc map[string]map[string]string) {
	muLabel.Lock()
	defer muLabel.Unlock()
	currentLabels = hsm
	cityDbConfig = cdc
	InitMysqlResource(hsm)
	for k, _ := range hsm {
		fmt.Println(k)
	}
}

// 实现方法
func NewRs(label string) (MysqlResourcerInterface, error) {
	return NewMysqlResourcer(label), nil
}

// 内部就结构体
type zgoMysqlService struct {
	label string
	res   MysqlResourcerInterface //使用resource另外的一个接口
}

func MysqlService(label string) (MysqlServiceInterface, error) {
	res, err := NewRs(label)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return &zgoMysqlService{label, res}, nil
}
func MysqlServiceByCityBiz(city string, biz string) (MysqlServiceInterface, error) {
	label, err := GetLabelByCityBiz(city, biz)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	res, err := NewRs(label)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return &zgoMysqlService{label, res}, nil
}
func GetDbByCityBiz(city string, biz string) (string, error) {
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

// 根据城市
func GetLabelByCityBiz(city string, biz string) (string, error) {
	if value, ok := cityDbConfig[biz]; ok {
		if value, ok := value[city]; ok {
			label := "mysql_" + biz + "_" + value
			return label, nil
		}
	}
	return "", errors.New(fmt.Sprintf("未知mysql label;city:%s;biz:%s;", city, biz))

}

// 对外接口 获取新的Service对象
func (c *zgoMysqlService) MysqlService(label string) (MysqlServiceInterface, error) {
	return MysqlService(label)
}

func (c *zgoMysqlService) MysqlServiceByCityBiz(city string, biz string) (MysqlServiceInterface, error) {
	return MysqlServiceByCityBiz(city, biz)
}

func (ms *zgoMysqlService) Get(ctx context.Context, args map[string]interface{}) error {
	err := ms.res.Get(ctx, args)
	return err
}

func (ms *zgoMysqlService) GetPool(t string) (*gorm.DB, error) {
	pool, err := ms.res.GetPool(t)
	return pool, err
}

func (ms *zgoMysqlService) List(ctx context.Context, args map[string]interface{}) error {
	return ms.res.List(ctx, args)
}

func (ms *zgoMysqlService) Count(ctx context.Context, args map[string]interface{}) error {
	return ms.res.Count(ctx, args)
}

func (ms *zgoMysqlService) Create(ctx context.Context, args map[string]interface{}) error {
	return ms.res.Get(ctx, args)
}

func (ms *zgoMysqlService) UpdateOne(ctx context.Context, args map[string]interface{}) (int, error) {
	return ms.res.UpdateOne(ctx, args)
}

func (ms *zgoMysqlService) DeleteOne(ctx context.Context, args map[string]interface{}) (int, error) {
	return ms.res.DeleteOne(ctx, args)
}

// 根据城市
func (c *zgoMysqlService) GetLabelByCityBiz(city string, biz string) (string, error) {
	return GetLabelByCityBiz(city, biz)
}

// 根据城市和业务 获取dbname和实例label
func (c *zgoMysqlService) GetDbByCityBiz(city string, biz string) (string, error) {
	return GetDbByCityBiz(city, biz)
}
