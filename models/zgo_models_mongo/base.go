package zgo_models_mongo

import (
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

type mongoResource struct {
	sess *mgo.Session
}

func (mongo mongoResource) MongoClient() *mgo.Session {
	return mongo.sess
}

func NewMongoResource(session *mgo.Session) (MongoResourcer, error) {
	return &mongoResource{
		sess: session,
	}, nil
}

func Login(mongoResourcer MongoResourcer, db, user, pass string) (*mgo.Session, error) {
	s := mongoResourcer.MongoClient().Copy()
	err := s.DB(db).Login(user, pass)
	return s, err
}

func Create(mongoResourcer MongoResourcer, db, collection string, args ...interface{}) (interface{}, error) {
	s := mongoResourcer.MongoClient().Copy()
	defer s.Close()
	return nil, s.DB(db).C(collection).Insert(args...)
}

// type bson.M map[string]interface{}
// test := bson.M{"name": "Jim"}
// test := bson.M{"name": bson.M{"first": "Jim"}}

func Get(mongoResourcer MongoResourcer, db, collection string, query bson.M) (interface{}, error) {
	s := mongoResourcer.MongoClient().Copy()
	defer s.Close()

	var res interface{}
	err := s.DB(db).C(collection).Find(&query).One(&res)
	return res, err
}

func GetBySelect(mongoResourcer MongoResourcer, db, collection string, query, bySelect bson.M) (interface{}, error) {
	s := mongoResourcer.MongoClient().Copy()
	defer s.Close()

	var res interface{}
	err := s.DB(db).C(collection).Find(&query).Select(&bySelect).One(&res)
	return res, err
}

func List(mongoResourcer MongoResourcer, db, collection string, query bson.M) ([]interface{}, error) {
	s := mongoResourcer.MongoClient().Copy()
	defer s.Close()

	ress := make([]interface{}, 0)
	err := s.DB(db).C(collection).Find(&query).All(&ress)
	return ress, err
}

func ListByLimit(mongoResourcer MongoResourcer, db, collection string, from, size int, query, bySelect bson.M, sort []string) ([]interface{}, error) {
	s := mongoResourcer.MongoClient().Copy()
	defer s.Close()

	ress := make([]interface{}, 0)
	err := s.DB(db).C(collection).Find(&query).Select(&bySelect).Skip(from).Limit(size).Sort(sort...).All(&ress)
	return ress, err
}

func UpdateOne(mongoResourcer MongoResourcer, db, col string, query bson.M, update bson.M) error {
	sess := mongoResourcer.MongoClient().Copy()
	defer sess.Close()

	return sess.DB(db).C(col).Update(query, update)
}

func UpdateAll(mongoResourcer MongoResourcer, db, col string, query bson.M, update bson.M) error {
	sess := mongoResourcer.MongoClient().Copy()
	defer sess.Close()

	_, err := sess.DB(db).C(col).UpdateAll(query, update)
	return err
}

func DeleteOne(mongoResourcer MongoResourcer, db, col string, query bson.M) error {
	sess := mongoResourcer.MongoClient().Copy()
	defer sess.Close()

	return sess.DB(db).C(col).Remove(query)
}

func DeleteAll(mongoResourcer MongoResourcer, db, col string, query bson.M) error {
	sess := mongoResourcer.MongoClient().Copy()
	defer sess.Close()

	_, err := sess.DB(db).C(col).RemoveAll(query)
	return err
}
