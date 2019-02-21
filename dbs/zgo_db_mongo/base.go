package zgo_db_mongo

import (
	"context"
	"errors"
	"fmt"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

type mongoResource struct {
	MongoChan chan *mgo.Session
}

func NewMongoResource() *mongoResource {
	return &mongoResource{
		MongoChan: MongoClientChan(),
	}
}

func Login(ch chan *mgo.Session, db, user, pass string) (*mgo.Session, error) {
	s := <-ch
	err := s.DB(db).Login(user, pass)
	return s, err
}

func Create(ch chan *mgo.Session, db, collection string, args ...interface{}) (interface{}, error) {
	s := <-ch
	return nil, s.DB(db).C(collection).Insert(args...)
}

// type bson.M map[string]interface{}
// test := bson.M{"name": "Jim"}
// test := bson.M{"name": bson.M{"first": "Jim"}}

func Get(ctx context.Context, ch chan *mgo.Session, args map[string]interface{} )(interface{}, error) {
	//out := make(chan interface{})
	//go func(ctx context.Context) {
	//	//defer close(out)
	//	s := <-ch
	//	var res interface{}
	//
	//	s.DB(db).C(collection).Find(&query).One(&res)
	//	out <- res
	//}(ctx)
	s := <-ch
	var res interface{}
	//judge if not exist
	err := judgeMongo(args)
	if err != nil{
		return nil, err
	}
	s.DB(args["db"].(string)).C(args["collection"].(string)).Find(args["query"].(bson.M)).One(&res)
	//out <- res
	return res, nil

}

func GetBySelect(ch chan *mgo.Session, args map[string]interface{}) (interface{}, error) {
	judgeMongo(args)
	s := <-ch
	var res interface{}
	query := map2Bson(args["query"].(map[string]interface{}))
	err := s.DB(args["db"].(string)).C(args["collection"].(string)).Find(query).
		Select(args["select"].(bson.M)).One(&res)
	return res, err
}

func List(ch chan *mgo.Session, args map[string]interface{}) ([]interface{}, error) {
	judgeMongo(args)
	s := <-ch
	//ress := make([]interface{}, 0)
	ress := []interface{}{}
	err := s.DB(args["db"].(string)).C(args["collection"].(string)).Find(args["query"].(bson.M)).
		Select(args["select"].(bson.M)).Skip(args["from"].(int)).Limit(args["limit"].(int)).Sort(args["sort"].([]string)...).All(&ress)
	return ress, err
}

func ListByLimit(ch chan *mgo.Session, db, collection string, from, size int, query, bySelect bson.M, sort []string) ([]interface{}, error) {
	s := <-ch

	ress := make([]interface{}, 0)
	err := s.DB(db).C(collection).Find(&query).Select(&bySelect).Skip(from).Limit(size).Sort(sort...).All(&ress)
	return ress, err
}

func UpdateOne(ch chan *mgo.Session, args map[string]interface{}) error {
	s := <-ch
	return s.DB(args["db"].(string)).C(args["collection"].(string)).Update(args["query"].(bson.M), args["update"].(bson.M))
}

func UpdateAll(ch chan *mgo.Session, args map[string]interface{}) error {
	s := <-ch

	_, err := s.DB(args["db"].(string)).C(args["collection"].(string)).UpdateAll(args["query"].(bson.M), args["update"].(bson.M))
	return err
}

func DeleteOne(ch chan *mgo.Session, db, col string, query bson.M) error {
	s := <-ch

	return s.DB(db).C(col).Remove(query)
}

func DeleteAll(ch chan *mgo.Session, db, col string, query bson.M) error {
	s := <-ch

	_, err := s.DB(db).C(col).RemoveAll(query)
	return err
}

func judgeMongo(args map[string]interface{})error{
	//switch args {
	//case args["db"] == nil:
	//	return errors.New("db is nil")
	//case args["collection"] == nil:
	//	return errors.New("collection is nil")
	//case args["query"] == nil:
	//	return errors.New("query is nil")
	//}
	if args["db"] == nil{
		return errors.New("db is nil")
	}else if  args["collection"] == nil{
		return errors.New("collection is nil")
	}else if  args["query"] == nil{
		return errors.New("query is nil")
	}
	return nil
}


func map2Bson(arg map[string]interface{})*bson.M{
	data, err := bson.Marshal(&arg)
	if err != nil{
		fmt.Println(err)
	}
	mmap := bson.M{}
	bson.Unmarshal(data, &mmap)
	return &mmap
}