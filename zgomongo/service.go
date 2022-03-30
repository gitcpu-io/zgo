package zgomongo

import (
  "context"
  "github.com/gitcpu-io/zgo/comm"
  "github.com/gitcpu-io/zgo/config"
  "github.com/globalsign/mgo"
  "sync"
)

var (
  currentLabels = make(map[string][]*config.ConnDetail)
  muLabel       = &sync.RWMutex{}
)

//Mongo 对外
type Mongoer interface {
  New(label ...string) (*zgomongo, error)
  GetConnChan(label ...string) (chan *mgo.Session, error)
  Create(ctx context.Context, args map[string]interface{}) error
  Update(ctx context.Context, args map[string]interface{}) error
  Upsert(ctx context.Context, args map[string]interface{}) error
  UpdateById(ctx context.Context, _id interface{}, args map[string]interface{}) error
  UpdateAll(ctx context.Context, args map[string]interface{}) error
  Delete(ctx context.Context, args map[string]interface{}) error
  DeleteById(ctx context.Context, _id interface{}, args map[string]interface{}) error
  DeleteAll(ctx context.Context, args map[string]interface{}) error
  FindOne(ctx context.Context, args map[string]interface{}) error
  FindPage(ctx context.Context, args map[string]interface{}) error
  Pipe(ctx context.Context, pipe interface{}, values interface{}, args map[string]interface{}) (interface{}, error)
  Count(ctx context.Context, args map[string]interface{}) (int, error)
  Get(ctx context.Context, args map[string]interface{}) error
  Insert(ctx context.Context, args map[string]interface{}) error
  //InsertMany(ctx context.Context, args map[string]interface{}, docs ...interface{}) error
}

func Mongo(l string) Mongoer {
  return &zgomongo{
    res: NewMongoResourcer(l),
  }
}

//zgomong实现了Mongo的接口
type zgomongo struct {
  res MongoResourcer //使用resource另外的一个接口
}

//InitMongo 初始化连接mongo
func InitMongo(hsmIn map[string][]*config.ConnDetail, label ...string) chan *zgomongo {
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

  InitMongoResource(hsm)

  //自动为变量初始化对象
  initLabel := ""
  for k := range hsm {
    if k != "" {
      initLabel = k
      break
    }
  }
  out := make(chan *zgomongo)
  go func() {
    in, err := GetMongo(initLabel)
    if err != nil {
      out <- nil
    }
    out <- in
    close(out)
  }()
  return out

}

//GetMongo zgo内部获取一个连接mongo
func GetMongo(label ...string) (*zgomongo, error) {
  l, err := comm.GetCurrentLabel(label, muLabel, currentLabels)
  if err != nil {
    return nil, err
  }

  return &zgomongo{
    res: NewMongoResourcer(l), //interface
  }, nil
}

func (n *zgomongo) New(label ...string) (*zgomongo, error) {
  return GetMongo(label...)
}

func (m *zgomongo) GetConnChan(label ...string) (chan *mgo.Session, error) {
  //label用来查找对应的库
  l, err := comm.GetCurrentLabel(label, muLabel, currentLabels)
  if err != nil {
    return nil, err
  }
  return m.res.GetConnChan(l), nil
}

func (m *zgomongo) Create(ctx context.Context, args map[string]interface{}) error {
  return m.res.Create(ctx, args)
}

func (m *zgomongo) Update(ctx context.Context, args map[string]interface{}) error {
  return m.res.UpdateOne(ctx, args)
}

func (m *zgomongo) Upsert(ctx context.Context, args map[string]interface{}) error {
  return m.res.Upsert(ctx, args)
}

func (m *zgomongo) UpdateById(ctx context.Context, _id interface{}, args map[string]interface{}) error {
  return m.res.UpdateById(ctx, _id, args)
}

func (m *zgomongo) UpdateAll(ctx context.Context, args map[string]interface{}) error {
  return m.res.UpdateAll(ctx, args)
}

func (m *zgomongo) Delete(ctx context.Context, args map[string]interface{}) error {
  return m.res.DeleteOne(ctx, args)
}

func (m *zgomongo) DeleteById(ctx context.Context, _id interface{}, args map[string]interface{}) error {
  return m.res.DeleteById(ctx, _id, args)
}

func (m *zgomongo) FindPage(ctx context.Context, args map[string]interface{}) error {
  //sort := args["sort"]
  if args["from"] == nil {
    args["from"] = 0
  }
  if args["size"] == nil {
    args["size"] = 10
  }
  if args["limit"] == nil {
    args["limit"] = 10
  }
  if args["sort"] == nil {
    args["sort"] = []string{}
  }
  return m.res.FindPage(ctx, args)
  //return zgo_db_mongo.List(ch, args["db"].(string), args["collection"].(string),
  //	args["query"].(bson.M))
}

func (m *zgomongo) Get(ctx context.Context, args map[string]interface{}) error {
  return m.res.Get(ctx, args)
}

func (m *zgomongo) Insert(ctx context.Context, args map[string]interface{}) error {
  return m.res.Insert(ctx, args)
}

func (m *zgomongo) FindOne(ctx context.Context, args map[string]interface{}) error {
  return m.res.FindOne(ctx, args)
}

func (m *zgomongo) Count(ctx context.Context, args map[string]interface{}) (int, error) {
  return m.res.Count(ctx, args)
}

func (m *zgomongo) Pipe(ctx context.Context, pipe interface{}, values interface{}, args map[string]interface{}) (interface{}, error) {
  return m.res.Pipe(ctx, pipe, values, args)
}
func (m *zgomongo) DeleteAll(ctx context.Context, args map[string]interface{}) error {
  return m.res.DeleteAll(ctx, args)
}

//func (m *zgomongo) InsertMany(ctx context.Context, args map[string]interface{}, docs ...interface{}) error {
//	return m.res.InsertMany(ctx, args, docs)
//}
