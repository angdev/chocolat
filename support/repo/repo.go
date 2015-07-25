package repo

import (
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
)

type Doc bson.M

type Repository struct {
	session  *mgo.Session
	database *mgo.Database
}

func NewRepository(dbName string) (*Repository, error) {
	session, err := mgo.Dial("localhost")
	if err != nil {
		return nil, err
	}

	return &Repository{session, session.DB(dbName)}, nil
}

func (this *Repository) C(name string) *mgo.Collection {
	return this.database.C(name)
}

func (this *Repository) Insert(name string, docs ...interface{}) error {
	return this.C(name).Insert(docs...)
}
