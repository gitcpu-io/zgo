package zgomongo

import (
  "context"
  "errors"
  "fmt"
  "github.com/gitcpu-io/zgo/config"
  "github.com/globalsign/mgo"
  "github.com/globalsign/mgo/bson"
)

//NsqResourcer 给service使用
type MongoResourcer interface {
  GetConnChan(label string) chan *mgo.Session
  Get(ctx context.Context, args map[string]interface{}) error
  Create(ctx context.Context, args map[string]interface{}) error
  Insert(ctx context.Context, args map[string]interface{}) error
  FindOne(ctx context.Context, args map[string]interface{}) error
  FindPage(ctx context.Context, args map[string]interface{}) error
  Count(ctx context.Context, args map[string]interface{}) (int, error)
  Pipe(ctx context.Context, pipe interface{}, values interface{}, args map[string]interface{}) (interface{}, error)
  UpdateOne(ctx context.Context, args map[string]interface{}) error
  Upsert(ctx context.Context, args map[string]interface{}) error
  UpdateById(ctx context.Context, id interface{}, args map[string]interface{}) error
  UpdateAll(ctx context.Context, args map[string]interface{}) error
  DeleteOne(ctx context.Context, args map[string]interface{}) error
  DeleteById(ctx context.Context, _id interface{}, args map[string]interface{}) error
  DeleteAll(ctx context.Context, args map[string]interface{}) error
}

//内部结构体
type mongoResource struct {
  label    string
  //mu       sync.RWMutex
  connpool ConnPooler
}

func NewMongoResourcer(label string) MongoResourcer {
  return &mongoResource{
    label:    label,
    connpool: NewConnPool(label), //使用connpool
  }
}

func InitMongoResource(hsm map[string][]*config.ConnDetail) {
  InitConnPool(hsm)
}

//GetConnChan 返回存放连接的chan
func (m *mongoResource) GetConnChan(label string) chan *mgo.Session {
  return m.connpool.GetConnChan(label)
}

func (m *mongoResource) Login(ctx context.Context, db, user, pass string) (*mgo.Session, error) {
  s := <-m.connpool.GetConnChan(m.label)

  err := s.DB(db).Login(user, pass)
  return s, err
}

func (m *mongoResource) Create(ctx context.Context, args map[string]interface{}) error {
  s := <-m.connpool.GetConnChan(m.label)
  return s.DB(args["db"].(string)).C(args["table"].(string)).Insert(args["items"])
}

func (m *mongoResource) Insert(ctx context.Context, args map[string]interface{}) error {
  s := <-m.connpool.GetConnChan(m.label)
  return s.DB(args["db"].(string)).C(args["table"].(string)).Insert(args["items"])
}

//func (m *mongoResource) InsertMany(ctx context.Context, args map[string]interface{}, docs ...interface{}) error {
//	s := <-m.connpool.GetConnChan(m.label)
//	collection := s.DB(args["db"].(string)).C(args["table"].(string))
//	return collection.Insert(docs)
//}

// type bson.M map[string]interface{}
// test := bson.M{"name": "Jim"}
// test := bson.M{"name": bson.M{"first": "Jim"}}

//获取一条数据，期望参数db table update， query参数不传则为全部查询
func (m *mongoResource) Get(ctx context.Context, args map[string]interface{}) error {
  s := <-m.connpool.GetConnChan(m.label)
  //judge if not exist
  err := judgeMongo(args)
  if err != nil {
    return err
  }
  if res, ok := args["obj"]; ok {
    preWorkForMongo(args)
    return s.DB(args["db"].(string)).C(args["table"].(string)).Find(args["query"].(bson.M)).Select(args["select"].(bson.M)).One(res)
  }
  return errors.New("未传入赋值对象")

}

func (m *mongoResource) Pipe(ctx context.Context, pipe interface{}, values interface{}, args map[string]interface{}) (interface{}, error) {
  s := <-m.connpool.GetConnChan(m.label)
  err := s.DB(args["db"].(string)).C(args["table"].(string)).Pipe(pipe).All(&values)
  if err != nil {
    return nil, err
  }
  return values, nil
}

func (m *mongoResource) FindOne(ctx context.Context, args map[string]interface{}) error {
  s := <-m.connpool.GetConnChan(m.label)
  if res, ok := args["obj"]; ok {
    err := judgeMongo(args)
    if err != nil {
      return err
    }
    preWorkForMongo(args)
    err = s.DB(args["db"].(string)).C(args["table"].(string)).Find(args["query"].(bson.M)).Select(args["select"].(bson.M)).Skip(args["from"].(int)).Limit(args["limit"].(int)).Sort(args["sort"].([]string)...).One(res)
    if err != nil {
      fmt.Println(err)
      return nil
    }
    return nil
  }
  return errors.New("未传入赋值对象")
}

//条件查询数量
func (m *mongoResource) Count(ctx context.Context, args map[string]interface{}) (int, error) {
  s := <-m.connpool.GetConnChan(m.label)
  err := judgeMongo(args)
  if err != nil {
    return 0, err
  }
  preWorkForMongo(args)
  return s.DB(args["db"].(string)).C(args["table"].(string)).Find(args["query"].(bson.M)).Count()
}

//分页查询 期望参数db table update query参数不传则为全部查询 select 为期望返回字段，若不填则自动为空
func (m *mongoResource) FindPage(ctx context.Context, args map[string]interface{}) error {
  if err := judgeMongo(args); err != nil {
    return err
  }
  s := <-m.connpool.GetConnChan(m.label)
  preWorkForMongo(args)
  ress := args["obj"]
  err := s.DB(args["db"].(string)).C(args["table"].(string)).Find(args["query"].(bson.M)).
    Select(args["select"].(bson.M)).Skip(args["from"].(int)).Limit(args["limit"].(int)).Sort(args["sort"].([]string)...).All(ress)
  return err
}

//修改多条数据，期望参数db table update query参数不传则为最后一条更新
func (m *mongoResource) UpdateOne(ctx context.Context, args map[string]interface{}) error {
  s := <-m.connpool.GetConnChan(m.label)
  err := judgeMongo(args)
  if err != nil {
    fmt.Println(err)
    return nil
  }
  preWorkForMongo(args)
  return s.DB(args["db"].(string)).C(args["table"].(string)).Update(args["query"].(bson.M), args["update"].(bson.M))
}

func (m *mongoResource) Upsert(ctx context.Context, args map[string]interface{}) error {
  s := <-m.connpool.GetConnChan(m.label)
  err := judgeMongo(args)
  if err != nil {
    fmt.Println(err)
    return nil
  }
  preWorkForMongo(args)
  _, err = s.DB(args["db"].(string)).C(args["table"].(string)).Upsert(args["query"].(bson.M), args["update"].(bson.M))
  return err
}

func (m *mongoResource) UpdateById(ctx context.Context, _id interface{}, args map[string]interface{}) error {
  s := <-m.connpool.GetConnChan(m.label)
  if args["db"] == nil {
    return errors.New("db is nil")
  } else if args["table"] == nil {
    return errors.New("table is nil")
  }
  if args["update"] != nil {
    update := updateMap2Bson(args["update"].(map[string]interface{}))
    args["update"] = update
  }
  var objId bson.ObjectId
  if val, ok := _id.(string); !ok {
    return errors.New("_id must be string")
  } else {
    objId = bson.ObjectIdHex(val)
  }
  return s.DB(args["db"].(string)).C(args["table"].(string)).UpdateId(objId, args["update"])
}

//修改多条数据，期望参数db table update query参数不传则为全部查询
func (m *mongoResource) UpdateAll(ctx context.Context, args map[string]interface{}) error {
  s := <-m.connpool.GetConnChan(m.label)
  if args["query"] == nil {
    return errors.New("query cannot be empty")
  }
  err := judgeMongo(args)
  if err != nil {
    fmt.Println(err)
    return nil
  }
  preWorkForMongo(args)
  _, err = s.DB(args["db"].(string)).C(args["table"].(string)).UpdateAll(args["query"].(bson.M), args["update"].(bson.M))
  return err
}

func (m *mongoResource) DeleteOne(ctx context.Context, args map[string]interface{}) error {
  err := judgeMongo(args)
  if err != nil {
    fmt.Println(err)
    return nil
  }
  s := <-m.connpool.GetConnChan(m.label)
  preWorkForMongo(args)
  return s.DB(args["db"].(string)).C(args["table"].(string)).Remove(args["query"])
}

func (m *mongoResource) DeleteById(ctx context.Context, _id interface{}, args map[string]interface{}) error {
  if args["db"] == nil {
    return errors.New("db is nil")
  } else if args["table"] == nil {
    return errors.New("table is nil")
  }
  s := <-m.connpool.GetConnChan(m.label)
  var objId bson.ObjectId
  if val, ok := _id.(string); !ok {
    return errors.New("_id must be string")
  } else {
    objId = bson.ObjectIdHex(val)
  }
  return s.DB(args["db"].(string)).C(args["table"].(string)).RemoveId(objId)
}

func (m *mongoResource) DeleteAll(ctx context.Context, args map[string]interface{}) error {
  err := judgeMongo(args)
  if err != nil {
    fmt.Println(err)
    return nil
  }
  s := <-m.connpool.GetConnChan(m.label)
  if args["query"] == nil {
    return errors.New("query param can not empty when operate delete  ")
  }
  preWorkForMongo(args)
  _, err = s.DB(args["db"].(string)).C(args["table"].(string)).RemoveAll(args["query"])
  return err
}

func judgeMongo(args map[string]interface{}) error {
  //switch args {
  //case args["db"] == nil:
  //	return errors.New("db is nil")
  //case args["table"] == nil:
  //	return errors.New("table is nil")
  //case args["query"] == nil:
  //	return errors.New("query is nil")
  //}
  if args["db"] == nil {
    return errors.New("db is nil")
  } else if args["table"] == nil {
    return errors.New("table is nil")
  } else if args["query"] == nil {
    return errors.New("query is nil")
  }
  return nil
}

//map转化成bson 内置函数
//func map2Bson(arg map[string]interface{}) bson.M {
func map2Bson(arg interface{}) bson.M {
  switch val := arg.(type) {
  case bson.M:
    return arg.(bson.M)
  case map[string]interface{}:
    return (bson.M)(arg.(map[string]interface{}))
  default:
    fmt.Printf("arg type is %v not bson.M or map[string]interface{} \n", val)
    return nil
  }
  //data, err := bson.Marshal(&arg)
  //if err != nil {
  //	fmt.Println(err)
  //}
  //mmap := bson.M{}
  //bson.Unmarshal(data, &mmap)
  //return mmap
}

//map转换成bson 并且添加$set  避免修改导致对删除其他字段  内置函数
func updateMap2Bson(arg map[string]interface{}) bson.M {
  set := map[string]map[string]interface{}{"$set": arg}
  data, err := bson.Marshal(&set)
  if err != nil {
    fmt.Println(err)
  }
  mmap := bson.M{}
  err = bson.Unmarshal(data, &mmap)
  if err != nil {
    fmt.Println(err)
    return nil
  }
  return mmap
}

//对mongo操作对前置工作 讲map类型转化为bson 内置函数
func preWorkForMongo(args map[string]interface{}) {
  if args["query"] != nil {
    query := map2Bson(args["query"])
    args["query"] = query
  } else {
    args["query"] = map[string]interface{}{}
    query := map2Bson(args["query"].(map[string]interface{}))
    args["query"] = query
  }

  if args["update"] != nil {
    update := updateMap2Bson(args["update"].(map[string]interface{}))
    args["update"] = update
  }

  if args["select"] != nil {
    bySelect := map2Bson(args["select"])
    args["select"] = bySelect
  } else {
    args["select"] = map[string]interface{}{}
    bySelect := map2Bson(args["select"].(map[string]interface{}))
    args["select"] = bySelect
  }

}
