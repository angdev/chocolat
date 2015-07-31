package query

import (
	"labix.org/v2/mgo"
)

type RawExpr map[string]interface{}
type Result map[string]interface{}

func New(c *mgo.Collection, arel *Arel) *Query {
	return &Query{Arel: arel, collection: c}
}

type Query struct {
	Arel       *Arel
	collection *mgo.Collection
}

func (this *Query) Execute() (interface{}, error) {
	if this.Arel.GroupByGiven() {
		return this.executeGroupBy()
	}
	return this.execute()
}

func (this *Query) executeGroupBy() ([]Result, error) {
	var result []Result
	p := this.Arel.Pipeline()
	if err := this.collection.Pipe(p).All(&result); err != nil {
		return nil, err
	} else {
		return result, nil
	}
}

func (this *Query) execute() (Result, error) {
	var result Result
	p := this.Arel.Pipeline()
	if err := this.collection.Pipe(p).One(&result); err != nil {
		return nil, err
	} else {
		return result, nil
	}
}
