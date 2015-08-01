package api

import (
	"github.com/angdev/chocolat/lib/query"
	"github.com/deckarep/golang-set"
)

func NewPresenter(q *query.Query) *Presenter {
	return &Presenter{query: q}
}

type Presenter struct {
	query *query.Query
}

func (this *Presenter) Present() (interface{}, error) {
	result, err := this.query.Execute()
	if err != nil {
		return nil, err
	} else {
		return query.Result{"result": this.deductResult(result)}, nil
	}
}

func (this *Presenter) deductResult(result interface{}) interface{} {
	switch result.(type) {
	case []query.Result:
		return result
	}
	r := result.(query.Result)
	return r["result"]
}

func (this *Presenter) PresentInterval(t *TimeFrame, i *Interval) (interface{}, error) {
	result, err := this.collectIntervalResult(t, i)
	if this.query.Arel.GroupByGiven() {
		this.ensureGroupField(result)
	}
	return query.Result{"result": result}, err
}

func (this *Presenter) collectIntervalResult(t *TimeFrame, i *Interval) ([]query.Result, error) {
	var results []query.Result

	start := t.Start
	for start.Before(t.End) {
		end := i.NextTime(start)
		this.query.Arel.Where(
			query.NewCondition("chocolat.created_at", "gt", start),
			query.NewCondition("chocolat.created_at", "lt", end))

		if result, err := this.query.Execute(); err != nil {
			return nil, err
		} else {
			results = append(results, query.Result{
				"timeframe": TimeFrame{Start: start, End: end},
				"value":     this.deductResult(result),
			})
		}

		start = end
	}

	return results, nil
}

func (this *Presenter) ensureGroupField(results []query.Result) {
	// Support only one group_by field
	groupName := this.query.Arel.ArelNodes.GroupBy.Group[0]
	values := mapset.NewSet()
	for _, result := range results {
		values = values.Union(this.collectGroupValues(result["value"].([]query.Result), groupName))
	}

	for _, result := range results {
		this.ensureGroupValues(&result, groupName, values)
	}
}

func (this *Presenter) collectGroupValues(result []query.Result, name string) mapset.Set {
	values := mapset.NewSet()
	for _, group := range result {
		values.Add(group[name])
	}
	return values
}

func (this *Presenter) ensureGroupValues(result *query.Result, groupName string, values mapset.Set) {
	results := (*result)["value"]
	missing := values.Difference(this.collectGroupValues(results.([]query.Result), groupName)).ToSlice()

	for _, v := range missing {
		missingResult := make(query.Result)
		missingResult[groupName] = v
		missingResult["result"] = nil

		results = append(results.([]query.Result), missingResult)
	}
	(*result)["value"] = results
}
