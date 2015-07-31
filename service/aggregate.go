package service

import (
	"github.com/angdev/chocolat/lib/query"
	"github.com/angdev/chocolat/model"
	"github.com/angdev/chocolat/support/repo"
)

type AggregateResult struct {
	Result interface{} `json:"result"`
}

type IntervalResult struct {
	Value     interface{} `json:"value"`
	TimeFrame *TimeFrame  `json:"timeframe"`
}

type Aggregator struct {
	project *model.Project
	params  *QueryParams
}

func NewAggregator(p *model.Project, params *QueryParams) *Aggregator {
	return &Aggregator{project: p, params: params}
}

// func intervalAggregate(p *Pipeline, params *QueryParams) (*AggregateResult, error) {
// 	if !params.TimeFrame.IsGiven() {
// 		return nil, errors.New("timeframe is not given")
// 	}

// 	var results []IntervalResult

// 	start := params.TimeFrame.Start
// 	frameEnd := params.TimeFrame.End

// 	for start.Before(frameEnd) {
// 		end := params.Interval.NextTime(start)
// 		t := TimeFrame{Start: start, End: end}
// 		if result, err := p.Copy().Prepend(t.Pipe()).Result(); err != nil {
// 			return nil, err
// 		} else {
// 			results = append(results, IntervalResult{Value: result, TimeFrame: &t})
// 			start = end
// 		}
// 	}

// 	return &AggregateResult{Result: results}, nil
// }

func (this *Aggregator) Count() (interface{}, error) {
	r := repo.NewRepository(this.project.RepoName())
	defer r.Close()

	// arel := NewArel().Where(this.params.Filters...).GroupBy(this.params.GroupBy...).Count()

	var conds []*query.Condition
	for _, f := range this.params.Filters {
		conds = append(conds, query.NewCondition(f.PropertyName, f.Operator, f.PropertyValue))
	}

	t := this.params.TimeFrame
	if t.IsGiven() {
		conds = append(conds,
			query.NewCondition("chocolat.created_at", "gt", t.Start),
			query.NewCondition("chocolat.created_at", "lt", t.End))
	}

	arel := query.NewArel().Where(conds...).GroupBy(this.params.GroupBy...).Count()
	q := query.New(r.C(this.params.CollectionName), arel)

	p := NewPresenter(q)

	i := this.params.Interval
	if i.IsGiven() {
		return p.PresentInterval(&t, &i)
	} else {
		return p.Present()
	}

	// c := r.C(this.params.CollectionName)
	// q := NewQuery(c, arel)

	// i := this.params.Interval

	// var result interface{}
	// var err error

	// if i.IsGiven() {
	// 	result, err = q.ExecuteWithInterval(t, i)
	// } else {
	// 	result, err = q.Execute()
	// }

	// return &AggregateResult{Result: result}, err
}

func (this *Aggregator) CountUnique(target string) (interface{}, error) {
	// r := repo.NewRepository(this.project.RepoName())
	// defer r.Close()

	// arel := NewArel().Where(this.params.Filters...).GroupBy(this.params.GroupBy...).CountUnique(target)
	// pp.Println(arel.Pipeline())

	// c := r.C(this.params.CollectionName)
	// pipeline := countUniquePipeline(c, this.params, target)

	// if this.params.Interval.IsGiven() {
	// 	result, err := intervalAggregate(pipeline, this.params)
	// 	return result, err
	// }

	// return aggregate(pipeline, this.params)
	return nil, nil
}

func (this *Aggregator) Min(target string) (interface{}, error) {
	// $group Op = $min
	return nil, nil
}

func (this *Aggregator) Max(target string) (interface{}, error) {
	// $group Op = $max
	return nil, nil
}

func (this *Aggregator) Sum(target string) (interface{}, error) {
	// $group Op = $sum
	return nil, nil
}

func (this *Aggregator) Average(target string) (interface{}, error) {
	// $group Op = $avg
	return nil, nil
}

func (this *Aggregator) Percentile(target string, percent int) (interface{}, error) {
	// 0. Get Count results
	// 1. For each result, Timeframe, GroupBy, Filter -> to one filter
	// 2. Match {1} -> Skip Percentile(Count, #percent) -> Limit 1
	return nil, nil
}

func (this *Aggregator) Median(target string) (interface{}, error) {
	// return Percentile(target, 50)
	return nil, nil
}

func (this *Aggregator) SelectUnique(target string) (interface{}, error) {
	// Unique w/ no counting
	return nil, nil
}
