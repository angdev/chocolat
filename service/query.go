package service

import (
	"labix.org/v2/mgo"
)

type Query struct {
	collection *mgo.Collection
	arel       *Arel
}

func NewQuery(c *mgo.Collection, arel *Arel) *Query {
	return &Query{collection: c, arel: arel}
}

func (this *Query) Execute() (interface{}, error) {
	var result interface{}
	var err error

	if this.arel.group.IsGiven() {
		result, err = this.executeGroupBy()
	} else {
		result, err = this.execute()
	}

	return result, err
}

func (this *Query) fetchResult() ([]QueryResult, error) {
	var result []QueryResult
	p := this.arel.Pipeline()

	if err := this.collection.Pipe(p.Stages()).All(&result); err != nil {
		return nil, err
	} else {
		return result, nil
	}
}

func (this *Query) execute() (interface{}, error) {
	result, err := this.fetchResult()
	if err != nil {
		return nil, err
	} else if result != nil {
		return result[0].Value(), nil
	} else {
		return 0, nil
	}
}

func (this *Query) executeGroupBy() (interface{}, error) {
	return this.fetchResult()
}

type QueryResult map[string]interface{}

func (this QueryResult) Value() interface{} {
	return this["result"]
}
