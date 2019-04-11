package zgomysql

import (
	"github.com/jinzhu/gorm"
)

// 对外接口
type Txer interface {

	//List(ctx context.Context, args map[string]interface{}) error
	//Count(ctx context.Context, args map[string]interface{}) error
	//Get(ctx context.Context, args map[string]interface{}) error
	//
	//Create(ctx context.Context, obj MysqlBaser) error
	//DeleteById(ctx context.Context, tableName string, id uint32) (int64, error)
	//UpdateNotEmptyByObj(ctx context.Context, obj MysqlBaser) (int, error)
	//UpdateByData(ctx context.Context, obj MysqlBaser, data map[string]interface{}) (int, error)
	//UpdateByObj(ctx context.Context, obj MysqlBaser) (int, error)
	//UpdateMany(ctx context.Context, tableName string, query string, args []interface{}, data map[string]interface{}) (int, error)
	//Exec(ctx context.Context, sql string, values ...interface{}) (int, error)
}

//内部结构体
type Tx struct {
	db *gorm.DB
}

// 对外函数 -- 创建mysqlResourcer对象
func NewTx(db *gorm.DB) Txer {
	db.Begin()
	return &Tx{
		db,
	}
}

//
//// 查询单个数据
//func (t *Tx) Get(ctx context.Context, args map[string]interface{}) (Tx,error) {
//	errv := mr.validate(args, "table", "query", "args", "obj")
//	if errv != nil {
//		return errv
//	}
//	var (
//		gormPool *gorm.DB
//		err      error
//	)
//	gormPool.Begin()
//	if T, ok := args["T"]; ok {
//		gormPool, err = mr.GetPool(T.(string))
//	} else {
//		gormPool, err = mr.GetRPool()
//	}
//
//	if err != nil {
//		return err
//	}
//	gormPool = gormPool.Table(args["table"].(string))
//	if sel, ok := args["select"]; ok {
//		gormPool = gormPool.Select(sel)
//	}
//	if order, ok := args["order"]; ok {
//		gormPool = gormPool.Order(order)
//	}
//	if group, ok := args["group"]; ok {
//		gormPool = gormPool.Group(group.(string))
//	}
//	err = gormPool.Where(args["query"], args["args"].([]interface{})...).First(args["obj"]).Error
//	return t,err
//}
//
//// 查询列表数据
//func (t *Tx) List(ctx context.Context, args map[string]interface{}) error {
//	errv := mr.validate(args, "table", "query", "args", "obj")
//	if errv != nil {
//		return errv
//	}
//	var (
//		gormPool *gorm.DB
//		err      error
//	)
//	if T, ok := args["T"]; ok {
//		gormPool, err = mr.GetPool(T.(string))
//	} else {
//		gormPool, err = mr.GetRPool()
//	}
//
//	if err != nil {
//		return err
//	}
//	gormPool = gormPool.Table(args["table"].(string))
//	if sel, ok := args["select"]; ok {
//		gormPool = gormPool.Select(sel)
//	}
//	gormPool = gormPool.Where(args["query"], args["args"].([]interface{})...)
//	currentLimit := 30
//	if limit, ok := args["limit"]; ok {
//		gormPool = gormPool.Limit(limit)
//		currentLimit = limit.(int)
//	} else {
//		gormPool = gormPool.Limit(currentLimit)
//	}
//	if page, ok := args["page"]; ok {
//		gormPool = gormPool.Offset((page.(int) - 1) * currentLimit)
//	} else if offset, ok := args["offset"]; ok {
//		gormPool = gormPool.Offset(offset)
//	}
//	if order, ok := args["order"]; ok {
//		gormPool = gormPool.Order(order)
//	}
//	if group, ok := args["group"]; ok {
//		gormPool = gormPool.Group(group.(string))
//	}
//	err = gormPool.Find(args["obj"]).Error
//	return err
//}
//
//// 查询数量
//func (t *Tx) Count(ctx context.Context, args map[string]interface{}) error {
//	errv := mr.validate(args, "table", "query", "args", "count")
//	if errv != nil {
//		return errv
//	}
//	var (
//		gormPool *gorm.DB
//		err      error
//	)
//	if T, ok := args["T"]; ok {
//		gormPool, err = mr.GetPool(T.(string))
//	} else {
//		gormPool, err = mr.GetRPool()
//	}
//
//	if err != nil {
//		return err
//	}
//	gormPool = gormPool.Table(args["table"].(string)).Where(args["query"], args["args"].([]interface{})...)
//	err = gormPool.Count(args["count"]).Error
//	return err
//}
//
//// 修改一条数据
//func (t *Tx) UpdateOne(ctx context.Context, args map[string]interface{}) (int, error) {
//	errv := mr.validate(args, "table", "data")
//	if errv != nil {
//		return 0, errv
//	}
//	gormPool, err := mr.GetWPool()
//	if err != nil {
//		return 0, err
//	}
//	if model, ok := args["model"]; ok {
//		// args["data"] = map[string]interface{}{"name": "hello", "age": 18}
//		gormPool = gormPool.Model(model)
//	}
//	gormPool = gormPool.Table(args["table"].(string))
//	if _, ok := args["id"]; ok {
//		// args["data"] = map[string]interface{}{"name": "hello", "age": 18}
//		gormPool = gormPool.Where(" id = ? ", args["id"])
//
//	}
//	if data, ok := args["data"]; ok {
//		db := gormPool.Updates(data)
//		count := db.RowsAffected
//		err = db.Error
//		return int(count), err
//	}
//
//	return 0, errors.New("mysql updateOne method : id not allow null or 0")
//}
//
//// 根据Id删除
//func (t *Tx) DeleteById(ctx context.Context, tableName string, id uint32) (int64, error) {
//	// 根据id删除
//	if id > 0 {
//		t.db.Table(tableName).Delete(nil, map[string]uint32{"id": id})
//		return t.db.RowsAffected, t.db.Error
//	}
//	return 0, errors.New("mysql deleteOne method : id not allow null or 0")
//}
//
//// 新增数据
//func (t *Tx) Create(ctx context.Context, obj MysqlBaser) error {
//	gormPool, err := mr.GetWPool()
//	if err != nil {
//		return err
//	}
//	if gormPool.NewRecord(obj) {
//		err = gormPool.Create(obj).Error
//		return err
//	} else {
//		return errors.New("被创建对象不能有主键")
//	}
//}
//
//// UpdateOneByData 根据data修改值
//func (t *Tx) UpdateByData(ctx context.Context, obj MysqlBaser, data map[string]interface{}) (int, error) {
//	gormPool, err := mr.GetWPool()
//	if err != nil {
//		return 0, err
//	}
//	if obj.GetID() == 0 {
//		return 0, errors.New("id不能为空")
//	}
//	gormPool = gormPool.Model(obj).Updates(data)
//	count := gormPool.RowsAffected
//	err = gormPool.Error
//	return int(count), err
//}
//
//// 更新数据，只更新非空字段
//func (t *Tx) UpdateNotEmptyByObj(ctx context.Context, obj MysqlBaser) (int, error) {
//	gormPool, err := mr.GetWPool()
//	if err != nil {
//		return 0, err
//	}
//	if obj.GetID() == 0 {
//		return 0, errors.New("id不能为空")
//	}
//	gormPool = gormPool.Model(obj).Update(obj)
//	count := gormPool.RowsAffected
//	err = gormPool.Error
//	return int(count), err
//}
//
//// 更新所有字段，不考虑非空 走回调方法
//func (t *Tx) UpdateByObj(ctx context.Context, obj MysqlBaser) (int, error) {
//	gormPool, err := mr.GetWPool()
//	if err != nil {
//		return 0, err
//	}
//	if obj.GetID() == 0 {
//		return 0, errors.New("id不能为空")
//	}
//	gormPool = gormPool.Model(obj).Omit(obj.Omit()).Save(obj)
//	count := gormPool.RowsAffected
//	err = gormPool.Error
//	return int(count), err
//}
//
//// UpdateMany 根据筛选条件批量修改数据 不支持回调方法
//func (t *Tx) UpdateMany(ctx context.Context, tableName string, query string, args []interface{}, data map[string]interface{}) (int, error) {
//	gormPool, err := mr.GetWPool()
//	if err != nil {
//		return 0, err
//	}
//	gormPool = gormPool.Table(tableName).Where(query, args).Updates(data)
//	count := gormPool.RowsAffected
//	err = gormPool.Error
//	return int(count), err
//}
//
//// Exec 执行sql语句
//func (t *Tx) Exec(ctx context.Context, sql string, values ...interface{}) (int, error) {
//	gormPool, err := mr.GetWPool()
//	if err != nil {
//		return 0, err
//	}
//	gormPool = gormPool.Exec(sql, values)
//	count := gormPool.RowsAffected
//	err = gormPool.Error
//	return int(count), err
//}
//
//// 校验参数是否齐全
//func (t *Tx) validate(args map[string]interface{}, fields ...string) error {
//	for _, v := range fields {
//		if _, ok := args[v]; !ok {
//			return errors.New(fmt.Sprintf("参数错误，%s 不存在", v))
//		}
//	}
//	return nil
//}