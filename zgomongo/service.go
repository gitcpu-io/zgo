package zgomongo

import (
	"context"
	"git.zhugefang.com/gocore/zgo.git/comm"
	"github.com/globalsign/mgo"
	"sync"
)

var (
	currentLabels = make(map[string][]string)
	muLabel       sync.RWMutex
)

//Mongo 对外
type Mongo interface {
	NewMongo(label ...string) (*zgomongo, error)
	GetConnChan(label ...string) (chan *mgo.Session, error)

	Create(ctx context.Context, args map[string]interface{}) (interface{}, error)
	Update(ctx context.Context, args map[string]interface{}) (interface{}, error)
	UpdateAll(ctx context.Context, args map[string]interface{}) (interface{}, error)
	Delete(ctx context.Context, args map[string]interface{}) (interface{}, error)
	List(ctx context.Context, args map[string]interface{}) ([]interface{}, error)
	Get(ctx context.Context, args map[string]interface{}) (interface{}, error)
}

//zgomong实现了Mongo的接口
type zgomongo struct {
	res MongoResourcer //使用resource另外的一个接口
}

//InitMongo 初始化连接mongo
func InitMongo(hsm map[string][]string) {
	muLabel.Lock()
	defer muLabel.Unlock()

	currentLabels = hsm
	InitMongoResource(hsm)
}

//GetMongo zgo内部获取一个连接mongo
func GetMongo(label ...string) (*zgomongo, error) {
	l, err := comm.GetCurrentLabel(label, muLabel, currentLabels)
	if err != nil {
		return nil, err
	}
	return &zgomongo{
		res: NewMongoResource(l), //interface
	}, nil
}

func (n *zgomongo) NewMongo(label ...string) (*zgomongo, error) {
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

func (m *zgomongo) Create(ctx context.Context, args map[string]interface{}) (interface{}, error) {
	return m.res.Create(ctx, args)
}

func (m *zgomongo) Update(ctx context.Context, args map[string]interface{}) (interface{}, error) {
	err := m.res.UpdateOne(ctx, args)
	if err != nil {
		return nil, err
	}
	return "success", err
}

func (m *zgomongo) UpdateAll(ctx context.Context, args map[string]interface{}) (interface{}, error) {
	err := m.res.UpdateAll(ctx, args)
	if err != nil {
		return nil, err
	}
	return "success", err
}

func (m *zgomongo) Delete(ctx context.Context, args map[string]interface{}) (interface{}, error) {

	return nil, m.res.DeleteOne(ctx, args)
}

func (m *zgomongo) List(ctx context.Context, args map[string]interface{}) ([]interface{}, error) {
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
	return m.res.List(ctx, args)
	//return zgo_db_mongo.List(ch, args["db"].(string), args["collection"].(string),
	//	args["query"].(bson.M))
}

func (m *zgomongo) Get(ctx context.Context, args map[string]interface{}) (interface{}, error) {
	return m.res.Get(ctx, args)
}
