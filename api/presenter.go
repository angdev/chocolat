package api

import (
	"encoding/json"
	"github.com/angdev/chocolat/lib/query"
	"github.com/deckarep/golang-set"
)

type queryResult struct {
	Result interface{} `json:"result"`
}

type queryGroupResult struct {
	Result interface{}
	Groups RawResult
}

func (this *queryGroupResult) MarshalJSON() ([]byte, error) {
	result := make(RawResult)
	result["result"] = this.Result

	for k, v := range this.Groups {
		result[k] = v
	}

	return json.Marshal(result)
}

type queryIntervalResult struct {
	Result    interface{} `json:"value"`
	TimeFrame TimeFrame   `json:"timeframe"`
}

func (this *queryIntervalResult) MarshalJSON() ([]byte, error) {
	result := make(RawResult)

	switch this.Result.(type) {
	case queryResult:
		result["value"] = this.Result.(queryResult).Result
	default:
		result["value"] = this.Result
	}

	result["timeframe"] = this.TimeFrame
	return json.Marshal(result)
}

func NewPresenter(q *query.Query) *Presenter {
	return &Presenter{query: q}
}

type Presenter struct {
	query *query.Query
}

func (this *Presenter) Present() (interface{}, error) {
	if this.query.Arel.GroupByGiven() {
		var result []queryGroupResult
		if err := this.query.Execute(&result); err != nil {
			return nil, err
		} else {
			return RawResult{"result": result}, nil
		}
	} else {
		var result queryResult

		if err := this.query.Execute(&result); err != nil {
			return nil, err
		} else {
			return result, nil
		}
	}
}

func (this *Presenter) PresentInterval(t *TimeFrame, i *Interval) (interface{}, error) {
	if this.query.Arel.GroupByGiven() {
		if result, err := this.collectIntervalGroupResult(t, i); err != nil {
			return nil, err
		} else {
			this.ensureGroupField(result)
			return RawResult{"result": result}, nil
		}
	} else {
		if result, err := this.collectIntervalResult(t, i); err != nil {
			return nil, err
		} else {
			return RawResult{"result": result}, nil
		}
	}
}

func (this *Presenter) collectIntervalResult(t *TimeFrame, i *Interval) ([]queryIntervalResult, error) {
	var results []queryIntervalResult
	var result queryResult

	start := t.Start
	for start.Before(t.End) {
		end := i.NextTime(start)
		this.query.Arel.Where(
			query.NewCondition("chocolat.created_at", "gt", start),
			query.NewCondition("chocolat.created_at", "lt", end))

		if err := this.query.Execute(&result); err != nil && err.Error() != "not found" {
			return nil, err
		} else {
			results = append(results, queryIntervalResult{
				Result: result,
				TimeFrame: TimeFrame{
					Start: start,
					End:   end,
				},
			})
		}

		start = end
	}

	return results, nil
}

func (this *Presenter) collectIntervalGroupResult(t *TimeFrame, i *Interval) ([]queryIntervalResult, error) {
	var results []queryIntervalResult
	var result []queryGroupResult

	start := t.Start
	for start.Before(t.End) {
		end := i.NextTime(start)
		this.query.Arel.Where(
			query.NewCondition("chocolat.created_at", "gt", start),
			query.NewCondition("chocolat.created_at", "lt", end))

		if err := this.query.Execute(&result); err != nil {
			return nil, err
		} else {
			results = append(results, queryIntervalResult{
				Result: result,
				TimeFrame: TimeFrame{
					Start: start,
					End:   end,
				},
			})
		}

		start = end
	}

	return results, nil
}

func (this *Presenter) ensureGroupField(results []queryIntervalResult) {
	// Support only one group_by field
	groupName := this.query.Arel.ArelNodes.GroupBy.Group[0]
	values := mapset.NewSet()
	for _, result := range results {
		values = values.Union(this.collectGroupValues(result.Result.([]queryGroupResult), groupName))
	}

	for i, _ := range results {
		results[i].Result = this.ensureGroupValues(results[i].Result.([]queryGroupResult), groupName, values)
	}
}

func (this *Presenter) collectGroupValues(results []queryGroupResult, name string) mapset.Set {
	values := mapset.NewSet()
	for _, result := range results {
		values.Add(result.Groups[name])
	}
	return values
}

func (this *Presenter) ensureGroupValues(results []queryGroupResult, groupName string, values mapset.Set) []queryGroupResult {
	groups := this.collectGroupValues(results, groupName)
	missingGroups := values.Difference(groups).ToSlice()
	ensured := results

	for _, v := range missingGroups {
		group := make(RawResult)
		group[groupName] = v

		ensured = append(ensured, queryGroupResult{
			Result: nil,
			Groups: group,
		})
	}

	return ensured
}
