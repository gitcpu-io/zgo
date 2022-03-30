package zgomysql

import (
  "context"
  "errors"
  "fmt"
  "github.com/gitcpu-io/zgo/config"
  "github.com/jinzhu/gorm"
  "reflect"
)

// 初始化 连接池
func InitMysqlResource(hsm map[string][]*config.ConnDetail) {
  InitConnPool(hsm)
}

// 对外接口
type MysqlResourcer interface {
  GetPool(t string) (*gorm.DB, error)
  GetRPool() (*gorm.DB, error)
  GetWPool() (*gorm.DB, error)
  List(ctx context.Context, gormPool *gorm.DB, args map[string]interface{}) error
  FindMaps(ctx context.Context, gormPool *gorm.DB, args map[string]interface{}) ([]map[string]interface{}, error)
  Count(ctx context.Context, gormPool *gorm.DB, args map[string]interface{}) error
  Get(ctx context.Context, gormPool *gorm.DB, args map[string]interface{}) error

  Create(ctx context.Context, gormPool *gorm.DB, obj MysqlBaser) error
  DeleteById(ctx context.Context, gormPool *gorm.DB, tableName string, id uint32) (int, error)
  DeleteMany(ctx context.Context, gormPool *gorm.DB, tableName string, query string, args []interface{}) (int, error)
  UpdateNotEmptyByObj(ctx context.Context, gormPool *gorm.DB, obj MysqlBaser) (int, error)
  UpdateByData(ctx context.Context, gormPool *gorm.DB, obj MysqlBaser, data map[string]interface{}) (int, error)
  UpdateByObj(ctx context.Context, gormPool *gorm.DB, obj MysqlBaser) (int, error)
  UpdateMany(ctx context.Context, gormPool *gorm.DB, tableName string, query string, args []interface{}, data map[string]interface{}) (int, error)
  Exec(ctx context.Context, gormPool *gorm.DB, sql string, values ...interface{}) (int, error)
  Raw(ctx context.Context, gormPool *gorm.DB, result interface{}, sql string, values ...interface{}) error
}

//内部结构体
type mysqlResource struct {
  label string
  //connpool *gorm.DB
}

// 对外函数 -- 创建mysqlResourcer对象
func NewMysqlResourcer(label string) MysqlResourcer {
  return &mysqlResource{
    label: label,
  }
}

// mysqlResourcer 实现方法
func (mr *mysqlResource) GetPool(t string) (*gorm.DB, error) {
  return GetPool(mr.label, t)
}

// mysqlResourcer 实现方法
func (mr *mysqlResource) GetRPool() (*gorm.DB, error) {
  c, err := GetPool(mr.label, "r")
  if err != nil {
    //zgo.Log.Error(err.Error())
    return GetPool(mr.label, "w")
  }
  return c, err
}

func (mr *mysqlResource) GetWPool() (*gorm.DB, error) {
  return GetPool(mr.label, "w")
}

// 查询单个数据
func (mr *mysqlResource) Get(ctx context.Context, gormPool *gorm.DB, args map[string]interface{}) error {
  errv := mr.validate(args, "table", "query", "args", "obj")
  if errv != nil {
    return errv
  }
  gormPool = gormPool.Table(args["table"].(string))
  if sel, ok := args["select"]; ok {
    gormPool = gormPool.Select(sel)
  }
  if join, ok := args["join"]; ok {
    gormPool = gormPool.Joins(join.(string))
  }
  if order, ok := args["order"]; ok {
    gormPool = gormPool.Order(order)
  }
  if group, ok := args["group"]; ok {
    gormPool = gormPool.Group(group.(string))
  }
  err := gormPool.Where(args["query"], args["args"].([]interface{})...).First(args["obj"]).Error
  if err != nil && err.Error() != "record not found" {
    return err
  }
  return nil
}

// 查询列表数据
func (mr *mysqlResource) List(ctx context.Context, gormPool *gorm.DB, args map[string]interface{}) error {
  errv := mr.validate(args, "table", "query", "args", "obj")
  if errv != nil {
    return errv
  }
  gormPool = gormPool.Table(args["table"].(string))
  if join, ok := args["join"]; ok {
    gormPool = gormPool.Joins(join.(string))
  }
  if sel, ok := args["select"]; ok {
    gormPool = gormPool.Select(sel)
  }
  gormPool = gormPool.Where(args["query"], args["args"].([]interface{})...)
  currentLimit := 30
  if limit, ok := args["limit"]; ok {
    gormPool = gormPool.Limit(limit)
    currentLimit = limit.(int)
  } else {
    gormPool = gormPool.Limit(currentLimit)
  }
  if page, ok := args["page"]; ok {
    gormPool = gormPool.Offset((page.(int) - 1) * currentLimit)
  } else if offset, ok := args["offset"]; ok {
    gormPool = gormPool.Offset(offset)
  }
  if order, ok := args["order"]; ok {
    gormPool = gormPool.Order(order)
  }
  if group, ok := args["group"]; ok {
    gormPool = gormPool.Group(group.(string))
  }
  err := gormPool.Find(args["obj"]).Error
  if err == gorm.ErrRecordNotFound {
    return nil
  }
  return err
}

// 查询列表数据
func (mr *mysqlResource) FindMaps(ctx context.Context, gormPool *gorm.DB, args map[string]interface{}) ([]map[string]interface{}, error) {
  errv := mr.validate(args, "table", "query", "args")
  if errv != nil {
    return nil, errv
  }
  gormPool = gormPool.Table(args["table"].(string))
  if join, ok := args["join"]; ok {
    gormPool = gormPool.Joins(join.(string))
  }
  if sel, ok := args["select"]; ok {
    gormPool = gormPool.Select(sel)
  }
  gormPool = gormPool.Where(args["query"], args["args"].([]interface{})...)
  currentLimit := 30
  if limit, ok := args["limit"]; ok {
    gormPool = gormPool.Limit(limit)
    currentLimit = limit.(int)
  } else {
    gormPool = gormPool.Limit(currentLimit)
  }
  if page, ok := args["page"]; ok {
    gormPool = gormPool.Offset((page.(int) - 1) * currentLimit)
  } else if offset, ok := args["offset"]; ok {
    gormPool = gormPool.Offset(offset)
  }
  if order, ok := args["order"]; ok {
    gormPool = gormPool.Order(order)
  }
  if group, ok := args["group"]; ok {
    gormPool = gormPool.Group(group.(string))
  }
  rows, err := gormPool.Rows()
  if err != nil {
    fmt.Println(err.Error())
    return nil, err
  }
  defer rows.Close()
  columns, _ := rows.Columns()
  //columnsT, _ := rows.ColumnTypes()
  //fmt.Println(columnsT)
  length := len(columns)
  results := make([]map[string]interface{}, 0)
  for rows.Next() {
    current := makeResultReceiver(length)
    if err := rows.Scan(current...); err != nil {
      panic(err)
    }
    value := make(map[string]interface{})
    for i := 0; i < length; i++ {
      key := columns[i]
      val := *(current[i]).(*interface{})
      if val == nil {
        value[key] = nil
        continue
      }
      vType := reflect.TypeOf(val)
      switch vType.String() {
      //case "int64":
      //	value[key] = val.(int64)
      //case "string":
      //	value[key] = val.(string)
      //case "time.Time":
      //	value[key] = val.(time.Time)
      case "[]uint8":
        //fmt.Printf("unsupport data type '%s' now\n", vType)
        value[key] = string(val.([]uint8))
      default:
        value[key] = val
        //fmt.Printf("unsupport data type '%s' now\n", vType)
        // TODO remember add other data type
      }
    }
    results = append(results, value)
  }
  return results, err
}

// 查询数量
func (mr *mysqlResource) Count(ctx context.Context, gormPool *gorm.DB, args map[string]interface{}) error {
  errv := mr.validate(args, "table", "query", "args", "count")
  if errv != nil {
    return errv
  }
  gormPool = gormPool.Table(args["table"].(string)).Where(args["query"], args["args"].([]interface{})...)
  if join, ok := args["join"]; ok {
    gormPool = gormPool.Joins(join.(string))
  }
  if group, ok := args["group"]; ok {
    gormPool = gormPool.Group(group.(string))
  }
  err := gormPool.Count(args["count"]).Error
  return err
}

// 修改一条数据
func (mr *mysqlResource) UpdateOne(ctx context.Context, gormPool *gorm.DB, args map[string]interface{}) (int, error) {
  errv := mr.validate(args, "table", "data")
  if errv != nil {
    return 0, errv
  }
  if model, ok := args["model"]; ok {
    // args["data"] = map[string]interface{}{"name": "hello", "age": 18}
    gormPool = gormPool.Model(model)
  }
  gormPool = gormPool.Table(args["table"].(string))
  if _, ok := args["id"]; ok {
    // args["data"] = map[string]interface{}{"name": "hello", "age": 18}
    gormPool = gormPool.Where(" id = ? ", args["id"])

  }
  if data, ok := args["data"]; ok {
    db := gormPool.Updates(data)
    count := db.RowsAffected
    err := db.Error
    return int(count), err
  }

  return 0, errors.New("mysql updateOne method : id not allow null or 0")
}

// 根据Id删除
func (mr *mysqlResource) DeleteById(ctx context.Context, gormPool *gorm.DB, tableName string, id uint32) (int, error) {
  // 根据id删除
  if id > 0 {
    gormPool = gormPool.Table(tableName).Delete(nil, "id = ?", id)
    return int(gormPool.RowsAffected), gormPool.Error
  }
  return 0, errors.New("mysql deleteOne method : id not allow null or 0")
}

// 新增数据
func (mr *mysqlResource) Create(ctx context.Context, gormPool *gorm.DB, obj MysqlBaser) error {
  //if gormPool.NewRecord(obj) {
  err := gormPool.Create(obj).Error
  return err
  //} else {
  //	return errors.New("被创建对象不能有主键")
  //}
}

// UpdateOneByData 根据data修改值
func (mr *mysqlResource) UpdateByData(ctx context.Context, gormPool *gorm.DB, obj MysqlBaser, data map[string]interface{}) (int, error) {
  if obj.GetID() == 0 {
    return 0, errors.New("id不能为空")
  }
  gormPool = gormPool.Model(obj).Updates(data)
  count := gormPool.RowsAffected
  err := gormPool.Error
  return int(count), err
}

// 更新数据，只更新非空字段
func (mr *mysqlResource) UpdateNotEmptyByObj(ctx context.Context, gormPool *gorm.DB, obj MysqlBaser) (int, error) {
  if obj.GetID() == 0 {
    return 0, errors.New("id不能为空")
  }
  gormPool = gormPool.Model(obj).Update(obj)
  count := gormPool.RowsAffected
  err := gormPool.Error
  return int(count), err
}

// 更新所有字段，不考虑非空 走回调方法
func (mr *mysqlResource) UpdateByObj(ctx context.Context, gormPool *gorm.DB, obj MysqlBaser) (int, error) {
  if obj.GetID() == 0 {
    return 0, errors.New("id不能为空")
  }
  gormPool = gormPool.Model(obj).Omit(obj.Omit()).Save(obj)
  count := gormPool.RowsAffected
  err := gormPool.Error
  return int(count), err
}

// UpdateMany 根据筛选条件批量修改数据 不支持回调方法
func (mr *mysqlResource) UpdateMany(ctx context.Context, gormPool *gorm.DB, tableName string, query string, args []interface{}, data map[string]interface{}) (int, error) {
  gormPool = gormPool.Table(tableName).Where(query, args...).Updates(data)
  count := gormPool.RowsAffected
  err := gormPool.Error
  return int(count), err
}

// UpdateMany 根据筛选条件批量修改数据 不支持回调方法
func (mr *mysqlResource) DeleteMany(ctx context.Context, gormPool *gorm.DB, tableName string, query string, args []interface{}) (int, error) {
  gormPool = gormPool.Table(tableName).Where(query, args...).Delete(nil)
  count := gormPool.RowsAffected
  err := gormPool.Error
  return int(count), err
}

// Exec 执行sql语句
func (mr *mysqlResource) Exec(ctx context.Context, gormPool *gorm.DB, sql string, values ...interface{}) (int, error) {
  gormPool = gormPool.Exec(sql, values...)
  count := gormPool.RowsAffected
  err := gormPool.Error
  return int(count), err
}

// Exec 执行sql语句
func (mr *mysqlResource) Raw(ctx context.Context, gormPool *gorm.DB, result interface{}, sql string, values ...interface{}) error {
  gormPool = gormPool.Raw(sql, values...)
  gormPool.Scan(result)
  err := gormPool.Error
  if err != nil {
    fmt.Println(err.Error())
  }
  return err
}

// 校验参数是否齐全
func (mr *mysqlResource) validate(args map[string]interface{}, fields ...string) error {
  for _, v := range fields {
    if _, ok := args[v]; !ok {
      return fmt.Errorf("参数错误，%s 不存在", v)
    }
  }
  return nil
}
func makeResultReceiver(length int) []interface{} {
  result := make([]interface{}, 0, length)
  for i := 0; i < length; i++ {
    var current = struct{}{}
    result = append(result, &current)
  }
  return result
}
