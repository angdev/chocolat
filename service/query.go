package service

// import (
// 	"github.com/angdev/chocolat/lib/query"
// 	"gopkg.in/fatih/set.v0"
// 	"labix.org/v2/mgo"
// )

// type Query struct {
// 	collection *mgo.Collection
// 	arel       *query.Arel
// }

// func NewQuery(c *mgo.Collection, arel *query.Arel) *Query {
// 	return &Query{collection: c, arel: arel}
// }

// func (this *Query) Execute() (interface{}, error) {
// 	var result interface{}
// 	var err error

// 	if this.arel.group.IsGiven() {
// 		result, err = this.executeGroupBy()
// 	} else {
// 		result, err = this.execute()
// 	}

// 	return result, err
// }

// func (this *Query) ExecuteWithInterval(t TimeFrame, i Interval) (interface{}, error) {
// 	var results []IntervalResult

// 	start := t.Start
// 	frameEnd := t.End

// 	for start.Before(frameEnd) {
// 		end := i.NextTime(start)
// 		timeframe := TimeFrame{Start: start, End: end}
// 		arel := this.arel.Clone().TimeFrame(timeframe)
// 		q := NewQuery(this.collection, arel)

// 		if result, err := q.Execute(); err != nil {
// 			return nil, err
// 		} else {
// 			results = append(results, IntervalResult{Value: result, TimeFrame: &timeframe})
// 		}

// 		start = end
// 	}

// 	if this.arel.group.IsGiven() {
// 		results = this.flattenResults(results)
// 	}
// 	return results, nil
// }

// func (this *Query) flattenResults(results []IntervalResult) []IntervalResult {
// 	// Support only one group by field
// 	groupColumn := this.arel.group.groups[0]
// 	values := set.New()
// 	for _, result := range results {
// 		values = set.Union(values, result.Value.(GroupQueryResult).GroupValues(groupColumn)).(*set.Set)
// 	}

// 	var r []IntervalResult
// 	for _, result := range results {
// 		result.Value = result.Value.(GroupQueryResult).ensureGroupValues(groupColumn, values)
// 		r = append(r, result)
// 	}

// 	return r
// }

// func (this *Query) fetchResult() ([]QueryResult, error) {
// 	var result []QueryResult
// 	p := this.arel.Pipeline()

// 	if err := this.collection.Pipe(p.Stages()).All(&result); err != nil {
// 		return nil, err
// 	} else {
// 		return result, nil
// 	}
// }

// func (this *Query) execute() (interface{}, error) {
// 	if result, err := this.fetchResult(); err != nil {
// 		return nil, err
// 	} else if result != nil {
// 		return QueryResult(result[0]), nil
// 	} else {
// 		return QueryResult(nil), nil
// 	}
// }

// func (this *Query) executeGroupBy() (interface{}, error) {
// 	if result, err := this.fetchResult(); err != nil {
// 		return nil, err
// 	} else {
// 		return GroupQueryResult(result), nil
// 	}
// }

// type QueryResult map[string]interface{}

// func (this QueryResult) Value() interface{} {
// 	if this == nil {
// 		return nil
// 	} else {
// 		return this["result"]
// 	}
// }

// type GroupQueryResult []QueryResult

// func (this GroupQueryResult) Value() interface{} {
// 	return this
// }

// func (this GroupQueryResult) GroupValues(groupName string) *set.Set {
// 	values := set.New()
// 	for _, result := range this {
// 		values.Add(result[groupName])
// 	}
// 	return values
// }

// func (this GroupQueryResult) ensureGroupValues(groupName string, values *set.Set) GroupQueryResult {
// 	notExisting := set.Difference(values, this.GroupValues(groupName)).List()
// 	ensured := this

// 	for _, value := range notExisting {
// 		ensured = append(ensured, QueryResult{
// 			"result":  nil,
// 			groupName: value,
// 		})
// 	}
// 	return ensured
// }

// type IntervalQueryResult struct {
// 	Value     interface{}
// 	TimeFrame TimeFrame
// }
