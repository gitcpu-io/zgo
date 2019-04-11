package zgomysql

import (
	"context"
	"errors"
	"fmt"
	"git.zhugefang.com/gocore/zgo/comm"
	"git.zhugefang.com/gocore/zgo/config"
	"github.com/jinzhu/gorm"
	"sync"
)

var currentLabels = make(map[string][]*config.ConnDetail)

var muLabel sync.RWMutex

//Mongo 对外
type Mysqler interface {
	//NewRs(label string) (MysqlResourcerInterface, error)
	New(label ...string) (Mysqler, error)
	GetPool(t string) (*gorm.DB, error)
	Get(ctx context.Context, args map[string]interface{}) error
	List(ctx context.Context, args map[string]interface{}) error
	Count(ctx context.Context, args map[string]interface{}) error
	//Create(ctx context.Context, args map[string]interface{}) error
	//UpdateOne(ctx context.Context, args map[string]interface{}) (int, error)
	//DeleteOne(ctx context.Context, args map[string]interface{}) (int, error)
	//GetLabelByCityBiz(city string, biz string) (string, error)
	//GetDbByCityBiz(city string, biz string) (string, error)
	//MysqlServiceByCityBiz(city string, biz string) (Mysqler, error)

	Create(ctx context.Context, obj MysqlBaser) error
	DeleteById(ctx context.Context, tableName string, id uint32) (int, error)
	DeleteByObj(ctx context.Context, obj MysqlBaser) (int, error)
	UpdateNotEmptyByObj(ctx context.Context, obj MysqlBaser) (int, error)
	UpdateByData(ctx context.Context, obj MysqlBaser, data map[string]interface{}) (int, error)
	UpdateByObj(ctx context.Context, obj MysqlBaser) (int, error)
	UpdateMany(ctx context.Context, tableName string, query string, args []interface{}, data map[string]interface{}) (int, error)
	Exec(ctx context.Context, sql string, values ...interface{}) (int, error)
}

// 内部就结构体
type zgoMysql struct {
	res MysqlResourcer //使用resource另外的一个接口
}

func Mysql(label string) Mysqler {
	return &zgoMysql{
		NewMysqlResourcer(label),
	}
}

// InitMysql 初始化连接mysql
func InitMysql(hsmIn map[string][]*config.ConnDetail, label ...string) chan *zgoMysql {
	muLabel.Lock()
	defer muLabel.Unlock()

	var hsm map[string][]*config.ConnDetail

	if len(label) > 0 && len(currentLabels) > 0 { //此时是destory操作,传入的hsm是nil
		//fmt.Println("--destory--前",currentLabels)
		for _, v := range label {
			delete(currentLabels, v)
		}
		hsm = currentLabels
		//fmt.Println("--destory--后",currentLabels)

	} else { //这是第一次创建操作或etcd中变更时init again操作
		hsm = hsmIn
		//currentLabels = hsm	//this operation is error
		for k, v := range hsm { //so big bug can't set hsm to currentLabels，must be for, may be have old label
			currentLabels[k] = v
		}
	}

	if len(hsm) == 0 {
		return nil
	}

	InitMysqlResource(hsm)

	//自动为变量初始化对象
	initLabel := ""
	for k, _ := range hsm {
		if k != "" {
			initLabel = k
			break
		}
	}
	out := make(chan *zgoMysql)
	go func() {
		in, err := GetMysql(initLabel)
		if err != nil {
			out <- nil
		}
		out <- in
		close(out)
	}()

	return out

}

// MysqlServiceByCityBiz
func MysqlServiceByCityBiz(city string, biz string) (Mysqler, error) {
	label, err := GetLabelByCityBiz(city, biz)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return &zgoMysql{NewMysqlResourcer(label)}, nil
}

// GetDbByCityBiz
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

// GetLabelByCityBiz根据城市
func GetLabelByCityBiz(city string, biz string) (string, error) {
	if value, ok := config.Conf.CityDbConfig[biz]; ok {
		if value, ok := value[city]; ok {
			label := "mysql_" + biz + "_" + value
			return label, nil
		}
	}
	return "", errors.New(fmt.Sprintf("未知mysql label;city:%s;biz:%s;", city, biz))

}

// GetMysql zgo内部获取一个连接mysql
func GetMysql(label ...string) (*zgoMysql, error) {
	l, err := comm.GetCurrentLabel(label, muLabel, currentLabels)
	if err != nil {
		return nil, err
	}
	return &zgoMysql{
		res: NewMysqlResourcer(l),
	}, nil
}

// New对外接口 获取新的Service对象
func (c *zgoMysql) New(label ...string) (Mysqler, error) {
	return GetMysql(label...)
}

// MysqlServiceByCityBiz
//func (c *zgoMysql) MysqlServiceByCityBiz(city string, biz string) (Mysqler, error) {
//	return MysqlServiceByCityBiz(city, biz)
//}

// Get
func (ms *zgoMysql) Get(ctx context.Context, args map[string]interface{}) error {
	db, err := ms.GetPool("r")
	if err != nil {
		return err
	}
	err = ms.res.Get(ctx, db, args)
	return err
}

// GetPool
func (ms *zgoMysql) GetPool(t string) (*gorm.DB, error) {
	if t == "r" {
		return ms.res.GetRPool()
	} else {
		return ms.res.GetWPool()
	}
}

func (ms *zgoMysql) getDB(ctx context.Context, T string, args map[string]interface{}) (*gorm.DB, error) {
	var (
		db  *gorm.DB
		err error
	)
	if t, ok := args["T"]; ok {
		db, err = ms.GetPool(t.(string))
	} else {
		db, err = ms.GetPool(T)
	}
	return db, err
}

// List
func (ms *zgoMysql) List(ctx context.Context, args map[string]interface{}) error {
	db, err := ms.getDB(ctx, "r", args)
	if err != nil {
		return err
	}
	return ms.res.List(ctx, db, args)
}

// Count
func (ms *zgoMysql) Count(ctx context.Context, args map[string]interface{}) error {
	return ms.res.Count(ctx, args)
}

// Create
//func (ms *zgoMysql) Create(ctx context.Context, args map[string]interface{}) error {
//	return ms.res.Create(ctx, args)
//}

// UpdateOne
//func (ms *zgoMysql) UpdateOne(ctx context.Context, args map[string]interface{}) (int, error) {
//	return ms.res.UpdateOne(ctx, args)
//}

// UpdateMany
//func (ms *zgoMysql) UpdateMany(ctx context.Context, args map[string]interface{}) (int, error) {
//	return ms.res.UpdateMany(ctx, args)
//}

// DeleteOne
//func (ms *zgoMysql) DeleteOne(ctx context.Context, args map[string]interface{}) (int, error) {
//	return ms.res.DeleteOne(ctx, args)
//}

//// GetLabelByCityBiz根据城市
//func (c *zgoMysql) GetLabelByCityBiz(city string, biz string) (string, error) {
//	return GetLabelByCityBiz(city, biz)
//}
//
//// GetDbByCityBiz根据城市和业务 获取dbname和实例label
//func (c *zgoMysql) GetDbByCityBiz(city string, biz string) (string, error) {
//	return GetDbByCityBiz(city, biz)
//}

func (ms *zgoMysql) Create(ctx context.Context, obj MysqlBaser) error {
	return ms.res.Create(ctx, obj)
}

func (ms *zgoMysql) DeleteById(ctx context.Context, tableName string, id uint32) (int, error) {
	db, err := ms.GetPool("w")
	if err != nil {
		return 0, err
	}
	return ms.res.DeleteById(ctx, db, tableName, id)
}

func (ms *zgoMysql) DeleteByObj(ctx context.Context, obj MysqlBaser) (int, error) {
	db, err := ms.GetPool("w")
	if err != nil {
		return 0, err
	}
	if obj.TableName() == "" {
		return 0, errors.New("表名不存在")
	}
	if obj.GetID() == 0 {
		return 0, errors.New("ID不能为0")
	}
	return ms.res.DeleteById(ctx, db, obj.TableName(), obj.GetID())
}

func (ms *zgoMysql) UpdateNotEmptyByObj(ctx context.Context, obj MysqlBaser) (int, error) {
	return ms.res.UpdateNotEmptyByObj(ctx, obj)
}

func (ms *zgoMysql) UpdateByData(ctx context.Context, obj MysqlBaser, data map[string]interface{}) (int, error) {
	return ms.res.UpdateByData(ctx, obj, data)
}

func (ms *zgoMysql) UpdateByObj(ctx context.Context, obj MysqlBaser) (int, error) {
	return ms.res.UpdateByObj(ctx, obj)
}

func (ms *zgoMysql) UpdateMany(ctx context.Context, tableName string, query string, args []interface{}, data map[string]interface{}) (int, error) {
	return ms.res.UpdateMany(ctx, tableName, query, args, data)
}

// Exec 执行原生sql
func (ms *zgoMysql) Exec(ctx context.Context, sql string, values ...interface{}) (int, error) {
	return ms.res.Exec(ctx, sql, values)
}
