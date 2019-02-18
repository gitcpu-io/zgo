package zgo_db_mongo

import (
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

func Get(ch chan *mgo.Session, db, collection string, query bson.M) (chan interface{}, error) {
	out := make(chan interface{})
	go func() {
		//defer close(out)
		s := <-ch
		var res interface{}
		s.DB(db).C(collection).Find(&query).One(&res)
		out <- res
	}()
	return out, nil

}

func GetBySelect(ch chan *mgo.Session, db, collection string, query, bySelect bson.M) (interface{}, error) {
	s := <-ch
	var res interface{}
	err := s.DB(db).C(collection).Find(&query).Select(&bySelect).One(&res)
	return res, err
}

func List(ch chan *mgo.Session, db, collection string, query bson.M) ([]interface{}, error) {
	s := <-ch

	ress := make([]interface{}, 0)
	err := s.DB(db).C(collection).Find(&query).All(&ress)
	return ress, err
}

func ListByLimit(ch chan *mgo.Session, db, collection string, from, size int, query, bySelect bson.M, sort []string) ([]interface{}, error) {
	s := <-ch

	ress := make([]interface{}, 0)
	err := s.DB(db).C(collection).Find(&query).Select(&bySelect).Skip(from).Limit(size).Sort(sort...).All(&ress)
	return ress, err
}

func UpdateOne(ch chan *mgo.Session, db, col string, query bson.M, update bson.M) error {
	s := <-ch

	return s.DB(db).C(col).Update(query, update)
}

func UpdateAll(ch chan *mgo.Session, db, col string, query bson.M, update bson.M) error {
	s := <-ch

	_, err := s.DB(db).C(col).UpdateAll(query, update)
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
