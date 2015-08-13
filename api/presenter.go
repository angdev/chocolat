package api

import (
	"github.com/deckarep/golang-set"
	"labix.org/v2/mgo"
)

func NewPresenter(a Aggregator, p *QueryParams) *Presenter {
	return &Presenter{aggregator: a, params: p}
}

type Presenter struct {
	aggregator Aggregator
	params     *QueryParams
}

func (this *Presenter) Present() (interface{}, error) {
	if this.params.Interval.IsGiven() && this.params.TimeFrame.IsGiven() {
		return this.presentInterval()
	} else {
		return this.present()
	}
}

func (this *Presenter) groupByGiven() bool {
	return len(this.params.GroupBy) != 0
}

func (this *Presenter) aggregate() (interface{}, error) {
	if this.groupByGiven() {
		var result queryGroupResultArray
		if err := this.aggregator(this.params, &result); err != nil {
			return nil, err
		} else {
			return result, nil
		}
	} else {
		var result queryResult
		if err := this.aggregator(this.params, &result); err != nil {
			return nil, err
		} else {
			return result, nil
		}
	}
}

func (this *Presenter) present() (interface{}, error) {
	result, err := this.aggregate()
	if err != nil {
		return nil, err
	}

	switch result.(type) {
	case queryResult:
		return result, nil
	case queryGroupResultArray:
		return RawResult{"result": result}, nil
	}

	return nil, nil
}

func (this *Presenter) presentInterval() (interface{}, error) {
	var (
		result interface{}
		err    error
	)

	if result, err = this.collectIntervalResult(); err != nil {
		return nil, err
	}

	if this.groupByGiven() {
		this.ensureGroupField(result.([]queryIntervalResult))
	}

	return RawResult{"result": result}, nil
}

func (this *Presenter) collectIntervalResult() (interface{}, error) {
	var results []queryIntervalResult

	t := this.params.TimeFrame
	i := this.params.Interval

	start := t.Start
	for start.Before(t.End) {
		end := i.NextTime(start)
		this.params.TimeFrame.Start = start
		this.params.TimeFrame.End = end

		if result, err := this.aggregate(); err != nil && err != mgo.ErrNotFound {
			return nil, err
		} else {
			results = append(results, queryIntervalResult{
				Result:    result,
				TimeFrame: this.params.TimeFrame,
			})
		}

		start = end
	}

	return results, nil
}

func (this *Presenter) ensureGroupField(results []queryIntervalResult) {
	// Support only one group_by field
	groupName := this.params.GroupBy[0]
	values := mapset.NewSet()
	for _, result := range results {
		values = values.Union(this.collectGroupValues(result.Result.(queryGroupResultArray), groupName))
	}

	for i, _ := range results {
		results[i].Result = this.ensureGroupValues(results[i].Result.(queryGroupResultArray), groupName, values)
	}
}

func (this *Presenter) collectGroupValues(results queryGroupResultArray, name string) mapset.Set {
	values := mapset.NewSet()
	for _, result := range results {
		values.Add(result.Groups[name])
	}
	return values
}

func (this *Presenter) ensureGroupValues(results queryGroupResultArray, groupName string, values mapset.Set) queryGroupResultArray {
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
