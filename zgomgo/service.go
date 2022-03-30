// zgomgo是对中间件Mongodb的封装，提供新建连接，各种CRUD接口
package zgomgo

import (
  "context"
  "errors"
  "fmt"
  "github.com/gitcpu-io/zgo/comm"
  "github.com/gitcpu-io/zgo/config"
  "go.mongodb.org/mongo-driver/bson"
  "go.mongodb.org/mongo-driver/bson/primitive"
  "go.mongodb.org/mongo-driver/mongo"
  "go.mongodb.org/mongo-driver/mongo/options"
  "sync"
)

var (
  currentLabels = make(map[string][]*config.ConnDetail) //用于存放label与具体Host:port的map
  muLabel       = &sync.RWMutex{}                            //用于并发读写上面的map
)

var (
  errNil = fmt.Errorf("%s", "文档实例不存在")
)

type MgoArgs struct {
  Document     interface{}              //保存时用到的结构体的指针
  Result       interface{}              //接受结构体的指针 比如: r := &User{} 这里的result就是r
  Filter       map[string]interface{}   //查询条件
  ArrayFilters []map[string]interface{} //子文档的查询条件
  Fields       map[string]interface{}   //字段筛选，形如SQL中的select选择字段
  Update       map[string]interface{}   //更新项 或 替换项
  Sort         map[string]interface{}   //排序 1是升序，-1是降序
  Limit        int64                    //限制数量
  Skip         int64                    //查询的offset，开区间，不包括这个skip的值
  Upsert       bool                     //当查询不到时，true表示插入一条新的
}

type MgoBulkWriteOperation = struct {
  Operation string
  MgoArgs   *MgoArgs
}

//Mgo 对外
type Mgoer interface {
  /*
   label: 可选，如果使用者，用了2个或多个label时，需要调用这个函数，传入label
  */
  // New 生产一条消息到Mgo
  New(label ...string) (*zgomgo, error)

  /*
   label: 可选，如果使用者，用了2个或多个label时，需要调用这个函数，传入label
  */
  // GetConnChan 获取原生的client，返回一个chan，使用者需要接收 <- chan
  GetConnChan(label ...string) (chan *mongo.Client, error)

  // GetCollection 获取文档实例
  GetCollection(dbName, collName string, label ...string) *mongo.Collection

  // FindById 通过Id查询文档，result为接受的结构体指针，比如: r := &User{} 这里的result就是r
  FindById(ctx context.Context, coll *mongo.Collection, result interface{}, id string) error

  // FindOne 通过条件只查询一条
  // 此时 args 可选如下：[
  // Filter: 查询条件
  // Fields: 字段筛选
  // Sort: 排序 1是升序，-1是降序
  // Skip: 查询的offset，开区间，不包括这个skip的值
  // ]
  FindOne(ctx context.Context, coll *mongo.Collection, args *MgoArgs) (uint8, error)

  // Find 查询多条,未查到返回空的[]
  // 此时 args 可选如下：[
  // Filter: 查询条件
  // Fields: 字段筛选
  // Sort: 排序 1是升序，-1是降序
  // Skip: 查询的offset，开区间，不包括这个skip的值
  // Limit: 限限的返回数量
  // ]
  Find(ctx context.Context, coll *mongo.Collection, args *MgoArgs) ([][]byte, error)

  // Count 查询数量,未查到返回0
  // 此时 args 可选如下：[
  // Filter: 查询条件
  // Skip: 查询的offset，开区间，不包括这个skip的值
  // Limit: 限限的返回数量
  // ]
  Count(ctx context.Context, coll *mongo.Collection, args *MgoArgs) (int64, error)

  // Insert 保存一条,返回objectId的string
  Insert(ctx context.Context, coll *mongo.Collection, document interface{}) (string, error)

  // InsertMany 保存多条,返回一个数组 []string,每一项的什是objectId的string
  InsertMany(ctx context.Context, coll *mongo.Collection, document []interface{}) ([]string, error)

  // UpdateById 通过查询ID来更新 传入的map的key表示要更新的字段，不写的字段不更新
  UpdateById(ctx context.Context, coll *mongo.Collection, update map[string]interface{}, id string) (uint8, error)

  // UpdateOne通过条件 更新一条, 返回第一个uint8表示更新的记录数, 第二个uint8表示插入的记录数,第三个string表示插入的id string
  // 更新与ReplaceOne的区别是：UpdateOne仅以args.update中的k,v覆盖当前记录的k对应的v值，如果这条记录以前有其它字段，不会更改原有的
  // 此时 args 可选如下：[
  // Filter: 查询条件
  // ArrayFilters: 子文档的查询条件
  // Update: 更新项
  // Upsert: 为true时插入一条
  // ]
  UpdateOne(ctx context.Context, coll *mongo.Collection, args *MgoArgs) (uint8, uint8, string, error)

  // ReplaceOne通过条件 替换一条, 返回第一个uint8表示替换的记录数, 第二个uint8表示插入的记录数,第三个string表示插入的id string
  // 替换与UpdateOne的区别是：ReplaceOne以args.update中的k,v重新为当前记录设置值，如果这条记录以前有其它字段，那么会删除以前的k、v
  // 但_id不变，执行结果是原来的_id，以及args.update中的k,v
  // 此时 args 可选如下：[
  // Filter: 查询条件
  // Update: 更新项
  // Upsert: 为true时插入一条
  // ]
  ReplaceOne(ctx context.Context, coll *mongo.Collection, args *MgoArgs) (uint8, uint8, string, error)

  // UpdateMany通过条件 更新多条, 返回第一个int64表示更新的记录数, 第二个int64表示插入的记录数,第三个int64表示查询到的记录数
  // 此时 args 可选如下：[
  // Filter: 查询条件
  // ArrayFilters: 子文档的查询条件
  // Update: 更新项
  // Upsert: 为true时插入一条
  // ]
  UpdateMany(ctx context.Context, coll *mongo.Collection, args *MgoArgs) (int64, int64, int64, error)

  // DeleteId 通过查询ID来删除一条
  DeleteById(ctx context.Context, coll *mongo.Collection, id string) (uint8, error)

  // DeleteOne通过条件 删除一条
  // 此时 args 可选如下：[
  // Filter: 查询条件
  // ]
  DeleteOne(ctx context.Context, coll *mongo.Collection, args *MgoArgs) (uint8, error)

  // DeleteMany通过条件 删除多条
  // 此时 args 可选如下：[
  // Filter: 查询条件
  // ]
  DeleteMany(ctx context.Context, coll *mongo.Collection, args *MgoArgs) (int64, error)

  // FindOneAndUpdate通过条件 查询并更新一条，默认返回最新的更新过的
  // 更新与FindOneAndReplace的区别是：FindOneAndUpdate仅以args.update中的k,v覆盖当前记录的k对应的v值，如果这条记录以前有其它字段，不会更改原有的
  // 此时 args 可选如下：[
  // Filter: 查询条件
  // Fields: 字段筛选
  // Result: 接受结构体的指针
  // Update: 更新项
  // Upsert: 为true时插入一条
  // Sort: 排序 1是升序，-1是降序
  // ]
  FindOneAndUpdate(ctx context.Context, coll *mongo.Collection, args *MgoArgs) error

  // FindOneAndReplace通过条件 查询并替换一条，默认返回最新的替换过的
  // 替换与FindOneAndUpdate的区别是：FindOneAndReplace以args.update中的k,v重新为当前记录设置值，如果这条记录以前有其它字段，那么会删除以前的k、v
  // 但_id不变，执行结果是原来的_id，以及args.update中的k,v
  // 此时 args 可选如下：[
  // Filter: 查询条件
  // ArrayFilters: 子文档的查询条件
  // Fields: 字段筛选
  // Result: 接受结构体的指针
  // Update: 更新项，此时为替换项
  // Upsert: 为true时插入一条
  // Sort: 排序 1是升序，-1是降序
  // ]
  FindOneAndReplace(ctx context.Context, coll *mongo.Collection, args *MgoArgs) error

  // FindOneAndDelete通过条件 查询并删除一条，返回当前删除的这条记录
  // 此时 args 可选如下：[
  // Filter: 查询条件
  // Fields: 字段筛选
  // Result: 接受结构体的指针
  // Sort: 排序 1是升序，-1是降序
  // ]
  FindOneAndDelete(ctx context.Context, coll *mongo.Collection, args *MgoArgs) error

  //Distinct 去重查询
  // fieldName是对哪个字段进行去重
  // filter是查询条件
  Distinct(ctx context.Context, coll *mongo.Collection, fieldName string, filter map[string]interface{}) ([]interface{}, error)

  // BulkWrite 多个并行计算
  // 1: 声明[]*MgoBulkWriteOperation
  // 2: 创建单个 &MgoBulkWriteOperation 结构体指针
  // {
  // 		Operation:"使用zgo.MgoBulkWriteOperation_xxxxxx" xxxxxx = insertOne/updateOne
  //      MgoArgs: 此时 args 可选如下：[
  // 			Document: 保存时用到的结构体的指针,用于insertOne
  //			Filter: 查询条件				用于updateOne/replaceOne/deleteOne/updateMany/deleteMany
  //			ArrayFilters: 子文档的查询条件 用于updateOne/updateMany
  //			Update: 更新项，此时为替换项	用于updateOne/replaceOne/updateMany
  //			Upsert: 为true时插入一条		用于updateOne/replaceOne/updateMany
  //		]
  // }
  // 3: append(第2步的结构体指针)到声明的[]中
  // 4: order 如果为true就是按第一步中[]的顺序执行，如果为false那么不管顺序，相当于并行计算
  // 使用实例请参考：
  // https://github.com/gitcpu-io/origin/blob/master/samples/demo_mgo/demo.go
  BulkWrite(ctx context.Context, coll *mongo.Collection, bulkWrites []*MgoBulkWriteOperation, order bool) (*mongo.BulkWriteResult, error)

  // Aggregate 聚合查询
  // 使用实例请参考：
  // https://github.com/gitcpu-io/origin/blob/master/samples/demo_mgo/demo.go
  Aggregate(ctx context.Context, coll *mongo.Collection, pipeline interface{}) ([][]byte, error)

  // Watch 监听
  Watch(ctx context.Context, coll *mongo.Collection, pipeline interface{}) (*mongo.ChangeStream, error)
}

// Mgo用于对zgo.Mgo这个全局变量赋值
func Mgo(label string) Mgoer {
  return &zgomgo{
    res: NewMgoResourcer(label),
  }
}

// zgomgo实现了Mgo的接口
type zgomgo struct {
  res MgoResourcer //使用resource另外的一个接口
}

// InitMgo 初始化连接mongo，用于使用者zgo.engine时，zgo init
func InitMgo(hsmIn map[string][]*config.ConnDetail, label ...string) chan *zgomgo {
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

  InitMgoResource(hsm)

  //自动为变量初始化对象
  initLabel := ""
  for k := range hsm {
    if k != "" {
      initLabel = k
      break
    }
  }
  out := make(chan *zgomgo)
  go func() {

    in, err := GetMgo(initLabel)
    if err != nil {
      panic(err)
    }
    out <- in
    close(out)
  }()

  return out

}

// GetMgo zgo内部获取一个连接mongo
func GetMgo(label ...string) (*zgomgo, error) {
  l, err := comm.GetCurrentLabel(label, muLabel, currentLabels)
  if err != nil {
    return nil, err
  }
  return &zgomgo{
    res: NewMgoResourcer(l),
  }, nil
}

// NewMgo获取一个Mgo生产者的client，用于发送数据
func (m *zgomgo) New(label ...string) (*zgomgo, error) {
  return GetMgo(label...)
}

//GetConnChan 供用户使用原生连接的chan
func (m *zgomgo) GetConnChan(label ...string) (chan *mongo.Client, error) {
  l, err := comm.GetCurrentLabel(label, muLabel, currentLabels)
  if err != nil {
    return nil, err
  }
  return m.res.GetConnChan(l), nil
}

//GetCollection 获取文档实例
func (m *zgomgo) GetCollection(dbName, collName string, label ...string) *mongo.Collection {
  return m.res.GetCollection(dbName, collName, label...)
}

// FindById 通过ID查询
func (m *zgomgo) FindById(ctx context.Context, coll *mongo.Collection, result interface{}, id string) error {
  if coll == nil {
    return errNil
  }
  oid, err := primitive.ObjectIDFromHex(id)
  if err != nil {
    return fmt.Errorf("%s:%s", "参数id不正确", err.Error())
  }
  bmq := bson.M{"_id": oid}
  err = m.res.FindById(ctx, coll, result, bmq)
  if err != nil {
    if err == mongo.ErrNoDocuments {
      return nil
    }
    return err
  }

  return err
}

// FindOne通过条件查询一条
func (m *zgomgo) FindOne(ctx context.Context, coll *mongo.Collection, args *MgoArgs) (uint8, error) {
  if coll == nil {
    return 0, errNil
  }

  m.dealObjectIdByString(args) //如果args.Filter中有_id就转为ObjectId

  bmq := bson.M(args.Filter)

  opts := &options.FindOneOptions{}
  if args.Fields != nil {
    opts.Projection = args.Fields
  }
  if args.Sort != nil {
    opts.Sort = args.Sort
  }
  if args.Skip > 0 {
    opts.Skip = &args.Skip
  }
  err := m.res.FindOne(ctx, coll, args.Result, bmq, opts)
  if err == mongo.ErrNoDocuments {
    //fmt.Println("没有结果")
    return 0, nil
  }
  return 1, nil
}

// Find通过条件查询多条
func (m *zgomgo) Find(ctx context.Context, coll *mongo.Collection, args *MgoArgs) ([][]byte, error) {
  if coll == nil {
    return nil, errNil
  }

  m.dealObjectIdByString(args) //如果args.Filter中有_id就转为ObjectId

  bmq := bson.M(args.Filter)

  opts := &options.FindOptions{}
  if args.Fields != nil {
    opts.Projection = args.Fields
  }
  if args.Sort != nil {
    opts.Sort = args.Sort
  }
  if args.Limit > 0 {
    opts.Limit = &args.Limit
  }
  if args.Skip > 0 {
    opts.Skip = &args.Skip
  }
  return m.res.Find(ctx, coll, bmq, opts)
}

// Count通过条件查询数量
func (m *zgomgo) Count(ctx context.Context, coll *mongo.Collection, args *MgoArgs) (int64, error) {
  if coll == nil {
    return 0, errNil
  }

  m.dealObjectIdByString(args) //如果args.Filter中有_id就转为ObjectId

  bmq := bson.M(args.Filter)

  opts := &options.CountOptions{}

  if args.Limit > 0 {
    opts.Limit = &args.Limit
  }
  if args.Skip > 0 {
    opts.Skip = &args.Skip
  }
  return m.res.Count(ctx, coll, bmq, opts)
}

// Insert 保存一条
func (m *zgomgo) Insert(ctx context.Context, coll *mongo.Collection, document interface{}) (string, error) {
  var id string
  if coll == nil {
    return id, errNil
  }
  var bv = true
  opts := &options.InsertOneOptions{
    BypassDocumentValidation: &bv,
  }

  result, err := m.res.Insert(ctx, coll, document, opts)
  if err != nil {
    return id, err
  }
  if val, ok := result.InsertedID.(primitive.ObjectID); ok && !val.IsZero() {
    id = val.Hex()
    return id, nil
  }
  return id, errors.New("not valid objectId")
}

// InsertMany 保存多条
func (m *zgomgo) InsertMany(ctx context.Context, coll *mongo.Collection, document []interface{}) ([]string, error) {
  var result []string
  if coll == nil {
    return result, errNil
  }
  var bv = true
  opts := &options.InsertManyOptions{
    BypassDocumentValidation: &bv,
  }

  res, err := m.res.InsertMany(ctx, coll, document, opts)
  if err != nil {
    return result, err
  }

  for _, v := range res.InsertedIDs {
    if val, ok := v.(primitive.ObjectID); ok && !val.IsZero() {
      result = append(result, val.Hex())
    }
  }

  return result, nil
}

// UpdateById 通过查询ID来更新
func (m *zgomgo) UpdateById(ctx context.Context, coll *mongo.Collection, update map[string]interface{}, id string) (uint8, error) {
  var modc uint8
  if coll == nil {
    return modc, errNil
  }
  oid, err := primitive.ObjectIDFromHex(id)
  if err != nil {
    return modc, fmt.Errorf("%s:%s", "参数id不正确", err.Error())
  }
  bmq := bson.M{"_id": oid}

  u := bson.M(update)

  opts := &options.UpdateOptions{}

  updateResult, err := m.res.UpdateOne(ctx, coll, bmq, u, opts)
  if err != nil {
    return modc, err
  }

  modc = uint8(updateResult.ModifiedCount)

  return modc, err
}

// UpdateOne通过条件 更新一条
func (m *zgomgo) UpdateOne(ctx context.Context, coll *mongo.Collection, args *MgoArgs) (uint8, uint8, string, error) {
  var (
    modc uint8
    upsc uint8
    upid string
  )
  if coll == nil {
    return modc, upsc, upid, errNil
  }

  m.dealObjectIdByString(args) //如果args.Filter中有_id就转为ObjectId

  bmq := bson.M(args.Filter)

  update := bson.M(args.Update)

  opts := &options.UpdateOptions{
    Upsert: &args.Upsert,
  }
  if args.ArrayFilters != nil && len(args.ArrayFilters) > 0 {
    tmp := options.ArrayFilters{}
    for _, v := range args.ArrayFilters {
      if len(v) > 0 {
        tmp.Filters = append(tmp.Filters, v)
      }
    }
    opts.ArrayFilters = &tmp
  }

  updateResult, err := m.res.UpdateOne(ctx, coll, bmq, update, opts)
  if err != nil {
    return modc, upsc, upid, err
  }

  if updateResult.UpsertedCount == 1 {
    if val, ok := updateResult.UpsertedID.(primitive.ObjectID); ok && !val.IsZero() {
      upid = val.Hex()
    }
  }

  modc = uint8(updateResult.ModifiedCount)
  upsc = uint8(updateResult.UpsertedCount)

  return modc, upsc, upid, nil
}

// ReplaceOne通过条件 替换一条
func (m *zgomgo) ReplaceOne(ctx context.Context, coll *mongo.Collection, args *MgoArgs) (uint8, uint8, string, error) {
  var (
    modc uint8
    upsc uint8
    upid string
  )
  if coll == nil {
    return modc, upsc, upid, errNil
  }

  m.dealObjectIdByString(args) //如果args.Filter中有_id就转为ObjectId

  bmq := bson.M(args.Filter)

  update := bson.M(args.Update)

  var bv = true
  opts := &options.ReplaceOptions{
    BypassDocumentValidation: &bv,
    Upsert:                   &args.Upsert,
  }

  updateResult, err := m.res.ReplaceOne(ctx, coll, bmq, update, opts)
  if err != nil {
    return modc, upsc, upid, err
  }

  if updateResult.UpsertedCount == 1 {
    if val, ok := updateResult.UpsertedID.(primitive.ObjectID); ok && !val.IsZero() {
      upid = val.Hex()
    }
  }

  modc = uint8(updateResult.ModifiedCount)
  upsc = uint8(updateResult.UpsertedCount)

  return modc, upsc, upid, nil
}

// UpdateMany通过条件 更新多条
func (m *zgomgo) UpdateMany(ctx context.Context, coll *mongo.Collection, args *MgoArgs) (int64, int64, int64, error) {
  var (
    modc int64
    upsc int64
    matc int64
  )
  if coll == nil {
    return modc, upsc, matc, errNil
  }

  m.dealObjectIdByString(args) //如果args.Filter中有_id就转为ObjectId

  bmq := bson.M(args.Filter)

  update := bson.M(args.Update)

  opts := &options.UpdateOptions{
    Upsert: &args.Upsert,
  }
  if args.ArrayFilters != nil && len(args.ArrayFilters) > 0 {
    tmp := options.ArrayFilters{}
    for _, v := range args.ArrayFilters {
      if len(v) > 0 {
        tmp.Filters = append(tmp.Filters, v)
      }
    }
    opts.ArrayFilters = &tmp
  }

  updateResult, err := m.res.UpdateMany(ctx, coll, bmq, update, opts)
  if err != nil {
    return modc, upsc, matc, err
  }

  modc = updateResult.ModifiedCount
  upsc = updateResult.UpsertedCount
  matc = updateResult.MatchedCount

  return modc, upsc, matc, nil
}

// DeleteById 通过查询ID来删除
func (m *zgomgo) DeleteById(ctx context.Context, coll *mongo.Collection, id string) (uint8, error) {
  var result uint8
  if coll == nil {
    return result, errNil
  }
  oid, err := primitive.ObjectIDFromHex(id)
  if err != nil {
    return result, fmt.Errorf("%s:%s", "参数id不正确", err.Error())
  }
  bmq := bson.M{"_id": oid}

  opts := &options.DeleteOptions{}

  deleteResult, err := m.res.DeleteOne(ctx, coll, bmq, opts)
  if err != nil {
    return result, err
  }

  result = uint8(deleteResult.DeletedCount)

  return result, err
}

// DeleteOne通过条件 删除一条
func (m *zgomgo) DeleteOne(ctx context.Context, coll *mongo.Collection, args *MgoArgs) (uint8, error) {
  var result uint8
  if coll == nil {
    return result, errNil
  }

  m.dealObjectIdByString(args) //如果args.Filter中有_id就转为ObjectId

  bmq := bson.M(args.Filter)

  opts := &options.DeleteOptions{}

  deleteResult, err := m.res.DeleteOne(ctx, coll, bmq, opts)
  if err != nil {
    return result, err
  }
  result = uint8(deleteResult.DeletedCount)
  return result, nil
}

// DeleteMany通过条件 删除多条
func (m *zgomgo) DeleteMany(ctx context.Context, coll *mongo.Collection, args *MgoArgs) (int64, error) {
  var result int64

  if coll == nil {
    return result, errNil
  }

  m.dealObjectIdByString(args) //如果args.Filter中有_id就转为ObjectId

  bmq := bson.M(args.Filter)

  opts := &options.DeleteOptions{}

  deleteResult, err := m.res.DeleteMany(ctx, coll, bmq, opts)
  if err != nil {
    return result, err
  }
  result = deleteResult.DeletedCount
  return result, nil
}

// FindOneAndUpdate通过条件 查询并更新一条
func (m *zgomgo) FindOneAndUpdate(ctx context.Context, coll *mongo.Collection, args *MgoArgs) error {

  if coll == nil {
    return errNil
  }

  m.dealObjectIdByString(args) //如果args.Filter中有_id就转为ObjectId

  bmq := bson.M(args.Filter)

  update := bson.M(args.Update)

  var bv = true
  var rd options.ReturnDocument = 1
  opts := &options.FindOneAndUpdateOptions{
    BypassDocumentValidation: &bv,
    Upsert:                   &args.Upsert,
    ReturnDocument:           &rd, //返回更新过的
  }
  if args.ArrayFilters != nil && len(args.ArrayFilters) > 0 {
    tmp := options.ArrayFilters{}
    for _, v := range args.ArrayFilters {
      if len(v) > 0 {
        tmp.Filters = append(tmp.Filters, v)
      }
    }
    opts.ArrayFilters = &tmp
  }
  if args.Fields != nil {
    opts.Projection = args.Fields
  }
  if args.Sort != nil {
    opts.Sort = args.Sort
  }

  return m.res.FindOneAndUpdate(ctx, coll, bmq, update, args.Result, opts)

}

// FindOneAndReplace通过条件 查询并替换一条
func (m *zgomgo) FindOneAndReplace(ctx context.Context, coll *mongo.Collection, args *MgoArgs) error {

  if coll == nil {
    return errNil
  }

  m.dealObjectIdByString(args) //如果args.Filter中有_id就转为ObjectId

  bmq := bson.M(args.Filter)

  update := bson.M(args.Update)

  var bv = true
  var rd options.ReturnDocument = 1
  opts := &options.FindOneAndReplaceOptions{
    BypassDocumentValidation: &bv,
    Upsert:                   &args.Upsert,
    ReturnDocument:           &rd, //返回更新过的
  }
  if args.Fields != nil {
    opts.Projection = args.Fields
  }
  if args.Sort != nil {
    opts.Sort = args.Sort
  }

  return m.res.FindOneAndReplace(ctx, coll, bmq, update, args.Result, opts)

}

// FindOneAndDelete通过条件 查询并删除一条，返回当前删除的这条记录
func (m *zgomgo) FindOneAndDelete(ctx context.Context, coll *mongo.Collection, args *MgoArgs) error {

  if coll == nil {
    return errNil
  }

  m.dealObjectIdByString(args) //如果args.Filter中有_id就转为ObjectId

  bmq := bson.M(args.Filter)

  opts := &options.FindOneAndDeleteOptions{}
  if args.Fields != nil {
    opts.Projection = args.Fields
  }
  if args.Sort != nil {
    opts.Sort = args.Sort
  }

  return m.res.FindOneAndDelete(ctx, coll, bmq, args.Result, opts)

}

// Distinct 去重
func (m *zgomgo) Distinct(ctx context.Context, coll *mongo.Collection, fieldName string, filter map[string]interface{}) ([]interface{}, error) {
  opts := &options.DistinctOptions{}

  bmq := bson.M(filter)

  return m.res.Distinct(ctx, coll, fieldName, bmq, opts)
}

// BulkWrite 多个并行计算
func (m *zgomgo) BulkWrite(ctx context.Context, coll *mongo.Collection, bulkWrites []*MgoBulkWriteOperation, order bool) (*mongo.BulkWriteResult, error) {
  if coll == nil {
    return nil, errNil
  }

  var writeModels []mongo.WriteModel

  for _, bukw := range bulkWrites {

    m.dealObjectIdByString(bukw.MgoArgs) //如果args.Filter中有_id就转为ObjectId

    switch bukw.Operation {
    case config.InsertOne:
      model := mongo.NewInsertOneModel()
      model.SetDocument(bukw.MgoArgs.Document)

      writeModels = append(writeModels, model)

    case config.UpdateOne:
      model := mongo.NewUpdateOneModel()
      if bukw.MgoArgs.ArrayFilters != nil && len(bukw.MgoArgs.ArrayFilters) > 0 {
        tmp := options.ArrayFilters{}
        for _, v := range bukw.MgoArgs.ArrayFilters {
          if len(v) > 0 {
            tmp.Filters = append(tmp.Filters, v)
          }
        }
        model.SetArrayFilters(tmp)
      }
      model.SetFilter(bson.M(bukw.MgoArgs.Filter))
      model.SetUpdate(bson.M(bukw.MgoArgs.Update))
      model.SetUpsert(bukw.MgoArgs.Upsert)

      writeModels = append(writeModels, model)

    case config.ReplaceOne:
      model := mongo.NewReplaceOneModel()
      model.SetFilter(bson.M(bukw.MgoArgs.Filter))
      model.SetReplacement(bson.M(bukw.MgoArgs.Update))
      model.SetUpsert(bukw.MgoArgs.Upsert)

      writeModels = append(writeModels, model)

    case config.DeleteOne:
      model := mongo.NewDeleteOneModel()
      model.SetFilter(bson.M(bukw.MgoArgs.Filter))

      writeModels = append(writeModels, model)

    case config.UpdateMany:
      model := mongo.NewUpdateManyModel()
      if bukw.MgoArgs.ArrayFilters != nil && len(bukw.MgoArgs.ArrayFilters) > 0 {
        tmp := options.ArrayFilters{}
        for _, v := range bukw.MgoArgs.ArrayFilters {
          if len(v) > 0 {
            tmp.Filters = append(tmp.Filters, v)
          }
        }
        model.SetArrayFilters(tmp)
      }
      model.SetFilter(bson.M(bukw.MgoArgs.Filter))
      model.SetUpdate(bson.M(bukw.MgoArgs.Update))
      model.SetUpsert(bukw.MgoArgs.Upsert)

      writeModels = append(writeModels, model)

    case config.DeleteMany:
      model := mongo.NewDeleteManyModel()
      model.SetFilter(bson.M(bukw.MgoArgs.Filter))

      writeModels = append(writeModels, model)

    }

  }

  opts := &options.BulkWriteOptions{
    Ordered: &order,
  }

  return m.res.BulkWrite(ctx, coll, writeModels, opts)
}

// Aggregate 聚合查询
func (m *zgomgo) Aggregate(ctx context.Context, coll *mongo.Collection, pipeline interface{}) ([][]byte, error) {
  var ad bool
  opts := &options.AggregateOptions{
    AllowDiskUse: &ad,
  }
  return m.res.Aggregate(ctx, coll, pipeline, opts)
}

// Watch 监听
func (m *zgomgo) Watch(ctx context.Context, coll *mongo.Collection, pipeline interface{}) (*mongo.ChangeStream, error) {
  var fd options.FullDocument = "updateLookup"
  opts := &options.ChangeStreamOptions{
    FullDocument: &fd,
  }
  return m.res.Watch(ctx, coll, pipeline, opts)
}

func (m *zgomgo) dealObjectIdByString(args *MgoArgs) {
  if args.Filter == nil {
    return
  }
  for k, v := range args.Filter {
    if k == "_id" {
      if id, ok := v.(string); ok { //如果是string就转objectId
        oid, err := primitive.ObjectIDFromHex(id)
        if err != nil {
          args.Filter[k] = id
        }
        args.Filter[k] = oid
      }
      break
    }
  }
}
