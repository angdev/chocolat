package service

import (
	"github.com/angdev/chocolat/lib/query"
	"github.com/angdev/chocolat/model"
	"github.com/angdev/chocolat/support/repo"
)

type Aggregator struct {
	project *model.Project
	params  *QueryParams
}

func NewAggregator(p *model.Project, params *QueryParams) *Aggregator {
	return &Aggregator{project: p, params: params}
}

func (this *Aggregator) Count() (interface{}, error) {
	r := repo.NewRepository(this.project.RepoName())
	defer r.Close()

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
}

func (this *Aggregator) CountUnique(target string) (interface{}, error) {
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
