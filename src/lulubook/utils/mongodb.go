package utils

import
(
	"github.com/globalsign/mgo"
)

var (
	session      *mgo.Session
)

func GetMongoDBSession() *mgo.Session {
	if session == nil {
		var err error
		session, err = mgo.Dial(CONFIG.MongodbServer)
		if err != nil {
			panic(err) // no, not really
		}
	}
	return session.Clone()
}

func WithMongoDBCollection(database, collection string, f func(*mgo.Collection) error) error {
	session := GetMongoDBSession()
	defer func() {
		session.Close()
	}()
	c := session.DB(database).C(collection)
	return f(c)
}

func DropMongoDB(database string) error{
	session := GetMongoDBSession()
	defer func() {
		session.Close()
	}()
	err := session.DB(database).DropDatabase()
	return err
}