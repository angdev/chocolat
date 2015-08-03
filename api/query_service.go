package api

import (
	"github.com/angdev/chocolat/lib/query"
	"github.com/angdev/chocolat/model"
	"github.com/angdev/chocolat/support/repo"
)

func count(p *model.Project, params *QueryParams) (interface{}, error) {
	r := repo.NewRepository(p.RepoName())
	defer r.Close()

	q := query.New(r.C(params.CollectionName), params.ToQuery().Count())

	presenter := NewPresenter(q)

	if params.Interval.IsGiven() {
		return presenter.PresentInterval(&params.TimeFrame, &params.Interval)
	} else {
		return presenter.Present()
	}
}

func countUnique(p *model.Project, target string, params *QueryParams) (interface{}, error) {
	return nil, nil
}

func min(target string) (interface{}, error) {
	// $group Op = $min
	return nil, nil
}

func max(target string) (interface{}, error) {
	// $group Op = $max
	return nil, nil
}

func sum(target string) (interface{}, error) {
	// $group Op = $sum
	return nil, nil
}

func average(target string) (interface{}, error) {
	// $group Op = $avg
	return nil, nil
}

func percentile(target string, percent int) (interface{}, error) {
	// 0. Get Count results
	// 1. For each result, Timeframe, GroupBy, Filter -> to one filter
	// 2. Match {1} -> Skip Percentile(Count, #percent) -> Limit 1
	return nil, nil
}

func median(target string) (interface{}, error) {
	// return Percentile(target, 50)
	return nil, nil
}

func selectUnique(target string) (interface{}, error) {
	// Unique w/ no counting
	return nil, nil
}
