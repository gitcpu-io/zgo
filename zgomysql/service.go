package zgomysql

import (
  "context"
  "errors"
  "fmt"
  "github.com/gitcpu-io/zgo/comm"
  "github.com/gitcpu-io/zgo/config"
  "github.com/jinzhu/gorm"
  "sync"
)

var currentLabels = make(map[string][]*config.ConnDetail)

var muLabel = &sync.RWMutex{}

//Mongo 对外
type Mysqler interface {
  //NewRs(label string) (MysqlResourcerInterface, error)
  New(label ...string) (Mysqler, error)
  NewMysql(label ...string) Mysqler
  GetDB(ctx context.Context, T string) (*gorm.DB, error)
  // --- 事物方法
  Begin(label string) Mysqler
  Commit() error
  RollBack()

  // --- 事物方法结束
  // --- 查询方法
  Get(ctx context.Context, args map[string]interface{}) error
  List(ctx context.Context, args map[string]interface{}) error
  FindMaps(ctx context.Context, args map[string]interface{}) ([]map[string]interface{}, error)
  Count(ctx context.Context, args map[string]interface{}) error
  // --- 查询方法结束

  // --- 增删改方法
  Create(ctx context.Context, obj MysqlBaser) error
  DeleteById(ctx context.Context, tableName string, id uint32) (int, error)
  DeleteByObj(ctx context.Context, obj MysqlBaser) (int, error)
  UpdateNotEmptyByObj(ctx context.Context, obj MysqlBaser) (int, error)
  UpdateByData(ctx context.Context, obj MysqlBaser, data map[string]interface{}) (int, error)
  UpdateByObj(ctx context.Context, obj MysqlBaser) (int, error)
  UpdateMany(ctx context.Context, tableName string, query string, args []interface{}, data map[string]interface{}) (int, error)
  DeleteMany(ctx context.Context, tableName string, query string, args []interface{}) (int, error)
  // --- 增删改方法 结束

  Exec(ctx context.Context, sql string, values ...interface{}) (int, error)
  Raw(ctx context.Context, result interface{}, sql string, values ...interface{}) error
}

// 内部就结构体
type zgoMysql struct {
  res MysqlResourcer //使用resource另外的一个接口
  db  *gorm.DB
}

func Mysql(label string) Mysqler {
  return &zgoMysql{
    NewMysqlResourcer(label),
    nil,
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
  for k := range hsm {
    if k != "" {
      initLabel = k
      break
    }
  }
  out := make(chan *zgoMysql)
  go func() {
    in := GetMysql(initLabel)
    out <- in
    close(out)
  }()

  return out

}

// GetMysql zgo内部获取一个连接mysql
func GetMysql(label ...string) *zgoMysql {
  l, err := comm.GetCurrentLabel(label, muLabel, currentLabels)
  if err != nil {
    panic(err)
  }
  return &zgoMysql{
    res: NewMysqlResourcer(l),
  }
}

// Get
func (ms *zgoMysql) GetDB(ctx context.Context, T string) (*gorm.DB, error) {
  return ms.GetPool(T)
}

// New对外接口 获取新的Service对象
func (c *zgoMysql) New(label ...string) (Mysqler, error) {
  return GetMysql(label...), nil
}

// New对外接口 获取新的Service对象
func (c *zgoMysql) NewMysql(label ...string) Mysqler {
  return GetMysql(label...)
}

func (ms *zgoMysql) Begin(label string) Mysqler {
  mysql := GetMysql(label)
  db, err := mysql.GetPool("w")
  if err != nil {
    fmt.Println(err.Error())
    panic(err)
  }
  db = db.Begin()
  return &zgoMysql{
    mysql.res,
    db,
  }
}

func (ms *zgoMysql) Commit() error {
  ms.db.Commit()
  return ms.db.Error
}

func (ms *zgoMysql) RollBack() {
  ms.db.Rollback()
}

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
  if ms.db != nil {
    return ms.db, nil
  }
  if t == "r" {
    return ms.res.GetRPool()
  } else {
    return ms.res.GetWPool()
  }
}

func (ms *zgoMysql) getDB(ctx context.Context, T string, args map[string]interface{}) (*gorm.DB, error) {
  if ms.db != nil {
    return ms.db, nil
  }
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

// List
func (ms *zgoMysql) FindMaps(ctx context.Context, args map[string]interface{}) ([]map[string]interface{}, error) {
  db, err := ms.getDB(ctx, "r", args)
  if err != nil {
    return nil, err
  }
  return ms.res.FindMaps(ctx, db, args)
}

// Count
func (ms *zgoMysql) Count(ctx context.Context, args map[string]interface{}) error {
  db, err := ms.getDB(ctx, "r", args)
  if err != nil {
    return err
  }
  return ms.res.Count(ctx, db, args)
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
  db, err := ms.GetPool("w")
  if err != nil {
    return err
  }
  return ms.res.Create(ctx, db, obj)
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
  db, err := ms.GetPool("w")
  if err != nil {
    return 0, err
  }
  return ms.res.UpdateNotEmptyByObj(ctx, db, obj)
}

func (ms *zgoMysql) UpdateByData(ctx context.Context, obj MysqlBaser, data map[string]interface{}) (int, error) {
  db, err := ms.GetPool("w")
  if err != nil {
    return 0, err
  }
  return ms.res.UpdateByData(ctx, db, obj, data)
}

func (ms *zgoMysql) UpdateByObj(ctx context.Context, obj MysqlBaser) (int, error) {
  db, err := ms.GetPool("w")
  if err != nil {
    return 0, err
  }
  return ms.res.UpdateByObj(ctx, db, obj)
}

func (ms *zgoMysql) UpdateMany(ctx context.Context, tableName string, query string, args []interface{}, data map[string]interface{}) (int, error) {
  db, err := ms.GetPool("w")
  if err != nil {
    return 0, err
  }
  return ms.res.UpdateMany(ctx, db, tableName, query, args, data)
}

func (ms *zgoMysql) DeleteMany(ctx context.Context, tableName string, query string, args []interface{}) (int, error) {
  db, err := ms.GetPool("w")
  if err != nil {
    return 0, err
  }
  return ms.res.DeleteMany(ctx, db, tableName, query, args)
}

// Exec 执行原生sql
func (ms *zgoMysql) Exec(ctx context.Context, sql string, values ...interface{}) (int, error) {
  db, err := ms.GetPool("w")
  if err != nil {
    return 0, err
  }
  return ms.res.Exec(ctx, db, sql, values...)
}

// Exec 执行原生sql
func (ms *zgoMysql) Raw(ctx context.Context, result interface{}, sql string, values ...interface{}) error {
  db, err := ms.GetPool("r")
  if err != nil {
    return err
  }
  return ms.res.Raw(ctx, db, result, sql, values...)
}
