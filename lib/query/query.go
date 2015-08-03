package query

import (
	"labix.org/v2/mgo"
)

type RawExpr map[string]interface{}

func New(c *mgo.Collection, arel *Arel) *Query {
	return &Query{Arel: arel, collection: c}
}

type Query struct {
	Arel       *Arel
	collection *mgo.Collection
}

func (this *Query) Execute(result interface{}) error {
	if this.Arel.GroupByGiven() {
		return this.executeGroupBy(result)
	}
	return this.execute(result)
}

func (this *Query) executeGroupBy(result interface{}) error {
	p := this.Arel.Pipeline()
	if err := this.collection.Pipe(p).All(result); err != nil {
		return err
	} else {
		return nil
	}
}

func (this *Query) execute(result interface{}) error {
	p := this.Arel.Pipeline()
	if err := this.collection.Pipe(p).One(result); err != nil {
		return err
	} else {
		return nil
	}
}
