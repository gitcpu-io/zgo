package zgomgo

import (
	"context"
	"github.com/gitcpu-io/zgo/config"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"sync"
)

//MgoResourcer 给service使用
type MgoResourcer interface {
	GetConnChan(label string) chan *mongo.Client
	GetCollection(dbName, collName string, label ...string) *mongo.Collection
	WithCollection(dbName, collName string, cb func(*mongo.Collection) error) error
	FindById(ctx context.Context, coll *mongo.Collection, result interface{}, filter bson.M) error
	FindOne(ctx context.Context, coll *mongo.Collection, result interface{}, filter bson.M, opts *options.FindOneOptions) error
	Find(ctx context.Context, coll *mongo.Collection, filter bson.M, opts *options.FindOptions) ([][]byte, error)
	Count(ctx context.Context, coll *mongo.Collection, filter bson.M, opts *options.CountOptions) (int64, error)
	Insert(ctx context.Context, coll *mongo.Collection, document interface{}, opts *options.InsertOneOptions) (*mongo.InsertOneResult, error)
	InsertMany(ctx context.Context, coll *mongo.Collection, document []interface{}, opts *options.InsertManyOptions) (*mongo.InsertManyResult, error)
	UpdateOne(ctx context.Context, coll *mongo.Collection, filter bson.M, update bson.M, opts *options.UpdateOptions) (*mongo.UpdateResult, error)
	ReplaceOne(ctx context.Context, coll *mongo.Collection, filter bson.M, replacement bson.M, opts *options.ReplaceOptions) (*mongo.UpdateResult, error)
	UpdateMany(ctx context.Context, coll *mongo.Collection, filter bson.M, update bson.M, opts *options.UpdateOptions) (*mongo.UpdateResult, error)
	DeleteOne(ctx context.Context, coll *mongo.Collection, filter bson.M, opts *options.DeleteOptions) (*mongo.DeleteResult, error)
	DeleteMany(ctx context.Context, coll *mongo.Collection, filter bson.M, opts *options.DeleteOptions) (*mongo.DeleteResult, error)

	FindOneAndUpdate(ctx context.Context, coll *mongo.Collection, filter bson.M, update bson.M, result interface{}, opts *options.FindOneAndUpdateOptions) error
	FindOneAndReplace(ctx context.Context, coll *mongo.Collection, filter bson.M, replacement bson.M, result interface{}, opts *options.FindOneAndReplaceOptions) error
	FindOneAndDelete(ctx context.Context, coll *mongo.Collection, filter bson.M, result interface{}, opts *options.FindOneAndDeleteOptions) error

	Distinct(ctx context.Context, coll *mongo.Collection, fieldName string, filter bson.M, opts *options.DistinctOptions) ([]interface{}, error)
	BulkWrite(ctx context.Context, coll *mongo.Collection, models []mongo.WriteModel, opts *options.BulkWriteOptions) (*mongo.BulkWriteResult, error)
	Aggregate(ctx context.Context, coll *mongo.Collection, pipeline interface{}, opts *options.AggregateOptions) ([][]byte, error)
	Watch(ctx context.Context, coll *mongo.Collection, pipeline interface{}, opts *options.ChangeStreamOptions) (*mongo.ChangeStream, error)
}

//内部结构体
type MgoResource struct {
	label    string
	mu       sync.RWMutex
	connpool ConnPooler
}

func NewMgoResourcer(label string) MgoResourcer {
	return &MgoResource{
		label:    label,
		connpool: NewConnPool(label), //使用connpool
	}
}

func InitMgoResource(hsm map[string][]*config.ConnDetail) {
	InitConnPool(hsm)
}

//GetConnChan 返回存放连接的chan
func (m *MgoResource) GetConnChan(label string) chan *mongo.Client {
	return m.connpool.GetConnChan(label)
}

func (m *MgoResource) GetCollection(dbName, collName string, label ...string) *mongo.Collection {
	var l string
	if len(label) > 0 {
		l = label[0]
	}
	if l == "" {
		l = m.label
	}
	client := <-m.connpool.GetConnChan(l)
	return client.Database(dbName).Collection(collName)
}

func (m *MgoResource) WithCollection(dbName, collName string, cb func(*mongo.Collection) error) error {
	client := <-m.connpool.GetConnChan(m.label)
	c := client.Database(dbName).Collection(collName)
	return cb(c)
}

func (m *MgoResource) FindById(ctx context.Context, coll *mongo.Collection, result interface{}, filter bson.M) error {

	//exop := func(c *mongo.Collection) error {
	//	return c.FindOne(ctx, filter).Decode(result)
	//}
	//
	//return m.WithCollection(dbName, collName, exop)

	return coll.FindOne(ctx, filter).Decode(result)
}

func (m *MgoResource) FindOne(ctx context.Context, coll *mongo.Collection, result interface{}, filter bson.M, opts *options.FindOneOptions) error {
	return coll.FindOne(ctx, filter, opts).Decode(result)
}

func (m *MgoResource) Find(ctx context.Context, coll *mongo.Collection, filter bson.M, opts *options.FindOptions) ([][]byte, error) {
	cur, err := coll.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var results [][]byte
	for cur.Next(ctx) {
		res := bson.M{}
		err := cur.Decode(&res)
		if err != nil {
			continue
		}
		bytes, err := bson.Marshal(res)
		if err != nil {
			continue
		}
		results = append(results, bytes)
	}
	if err := cur.Err(); err != nil {
		return results, err
	}
	return results, nil
}

func (m *MgoResource) Count(ctx context.Context, coll *mongo.Collection, filter bson.M, opts *options.CountOptions) (int64, error) {
	return coll.CountDocuments(ctx, filter, opts)
}

func (m *MgoResource) Insert(ctx context.Context, coll *mongo.Collection, document interface{}, opts *options.InsertOneOptions) (*mongo.InsertOneResult, error) {
	oneResult, err := coll.InsertOne(ctx, document, opts)
	if err != nil {
		return nil, err
	}
	return oneResult, nil
}

func (m *MgoResource) InsertMany(ctx context.Context, coll *mongo.Collection, document []interface{}, opts *options.InsertManyOptions) (*mongo.InsertManyResult, error) {
	manyResult, err := coll.InsertMany(ctx, document, opts)
	if err != nil {
		return nil, err
	}
	return manyResult, nil
}

func (m *MgoResource) UpdateOne(ctx context.Context, coll *mongo.Collection, filter bson.M, update bson.M, opts *options.UpdateOptions) (*mongo.UpdateResult, error) {
	updateResult, err := coll.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return nil, err
	}
	return updateResult, nil
}

func (m *MgoResource) ReplaceOne(ctx context.Context, coll *mongo.Collection, filter bson.M, replacement bson.M, opts *options.ReplaceOptions) (*mongo.UpdateResult, error) {
	replaceOne, err := coll.ReplaceOne(ctx, filter, replacement, opts)
	if err != nil {
		return nil, err
	}
	return replaceOne, nil
}

func (m *MgoResource) UpdateMany(ctx context.Context, coll *mongo.Collection, filter bson.M, update bson.M, opts *options.UpdateOptions) (*mongo.UpdateResult, error) {
	updateMany, err := coll.UpdateMany(ctx, filter, update, opts)
	if err != nil {
		return nil, err
	}
	return updateMany, nil
}

func (m *MgoResource) DeleteOne(ctx context.Context, coll *mongo.Collection, filter bson.M, opts *options.DeleteOptions) (*mongo.DeleteResult, error) {
	deleteResult, err := coll.DeleteOne(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	return deleteResult, nil
}

func (m *MgoResource) DeleteMany(ctx context.Context, coll *mongo.Collection, filter bson.M, opts *options.DeleteOptions) (*mongo.DeleteResult, error) {
	deleteMany, err := coll.DeleteMany(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	return deleteMany, nil
}

func (m *MgoResource) FindOneAndUpdate(ctx context.Context, coll *mongo.Collection, filter bson.M, update bson.M, result interface{}, opts *options.FindOneAndUpdateOptions) error {
	return coll.FindOneAndUpdate(ctx, filter, update, opts).Decode(result)
}

func (m *MgoResource) FindOneAndReplace(ctx context.Context, coll *mongo.Collection, filter bson.M, replacement bson.M, result interface{}, opts *options.FindOneAndReplaceOptions) error {
	return coll.FindOneAndReplace(ctx, filter, replacement, opts).Decode(result)
}

func (m *MgoResource) FindOneAndDelete(ctx context.Context, coll *mongo.Collection, filter bson.M, result interface{}, opts *options.FindOneAndDeleteOptions) error {
	return coll.FindOneAndDelete(ctx, filter, opts).Decode(result)
}

func (m *MgoResource) Distinct(ctx context.Context, coll *mongo.Collection, fieldName string, filter bson.M, opts *options.DistinctOptions) ([]interface{}, error) {
	return coll.Distinct(ctx, fieldName, filter, opts)
}

func (m *MgoResource) BulkWrite(ctx context.Context, coll *mongo.Collection, models []mongo.WriteModel, opts *options.BulkWriteOptions) (*mongo.BulkWriteResult, error) {
	return coll.BulkWrite(ctx, models, opts)
}

func (m *MgoResource) Aggregate(ctx context.Context, coll *mongo.Collection, pipeline interface{}, opts *options.AggregateOptions) ([][]byte, error) {
	cur, err := coll.Aggregate(ctx, pipeline, opts)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var results [][]byte
	for cur.Next(ctx) {
		res := bson.M{}
		err := cur.Decode(&res)
		if err != nil {
			continue
		}
		bytes, err := bson.Marshal(res)
		if err != nil {
			continue
		}
		results = append(results, bytes)
	}
	if err := cur.Err(); err != nil {
		return results, err
	}
	return results, nil

}

func (m *MgoResource) Watch(ctx context.Context, coll *mongo.Collection, pipeline interface{}, opts *options.ChangeStreamOptions) (*mongo.ChangeStream, error) {
	return coll.Watch(ctx, pipeline, opts)
}
