package repo

import (
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
)

var session *mgo.Session

type Doc bson.M

type Repository struct {
	session  *mgo.Session
	database *mgo.Database
}

func NewRepository(dbName string) *Repository {
	sess := session.Copy()
	return &Repository{sess, sess.DB(dbName)}
}

func (this *Repository) Close() {
	this.session.Close()
}

func (this *Repository) C(name string) *mgo.Collection {
	return this.database.C(name)
}

func (this *Repository) Insert(name string, docs ...interface{}) error {
	return this.C(name).Insert(docs...)
}

func init() {
	session, _ = mgo.Dial("localhost")
}
