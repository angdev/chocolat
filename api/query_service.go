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
	r := repo.NewRepository(p.RepoName())
	defer r.Close()

	q := query.New(r.C(params.CollectionName), params.ToQuery().CountUnique(target))

	presenter := NewPresenter(q)

	if params.Interval.IsGiven() {
		return presenter.PresentInterval(&params.TimeFrame, &params.Interval)
	} else {
		return presenter.Present()
	}
}

func min(p *model.Project, target string, params *QueryParams) (interface{}, error) {
	r := repo.NewRepository(p.RepoName())
	defer r.Close()

	q := query.New(r.C(params.CollectionName), params.ToQuery().Min(target))

	presenter := NewPresenter(q)

	if params.Interval.IsGiven() {
		return presenter.PresentInterval(&params.TimeFrame, &params.Interval)
	} else {
		return presenter.Present()
	}
}

func max(p *model.Project, target string, params *QueryParams) (interface{}, error) {
	r := repo.NewRepository(p.RepoName())
	defer r.Close()

	q := query.New(r.C(params.CollectionName), params.ToQuery().Max(target))

	presenter := NewPresenter(q)

	if params.Interval.IsGiven() {
		return presenter.PresentInterval(&params.TimeFrame, &params.Interval)
	} else {
		return presenter.Present()
	}
}

func sum(p *model.Project, target string, params *QueryParams) (interface{}, error) {
	r := repo.NewRepository(p.RepoName())
	defer r.Close()

	q := query.New(r.C(params.CollectionName), params.ToQuery().Sum(target))

	presenter := NewPresenter(q)

	if params.Interval.IsGiven() {
		return presenter.PresentInterval(&params.TimeFrame, &params.Interval)
	} else {
		return presenter.Present()
	}
}

func average(p *model.Project, target string, params *QueryParams) (interface{}, error) {
	r := repo.NewRepository(p.RepoName())
	defer r.Close()

	q := query.New(r.C(params.CollectionName), params.ToQuery().Average(target))

	presenter := NewPresenter(q)

	if params.Interval.IsGiven() {
		return presenter.PresentInterval(&params.TimeFrame, &params.Interval)
	} else {
		return presenter.Present()
	}
}

func percentile(p *model.Project, target string, percent int, params *QueryParams) (interface{}, error) {
	// 0. Get Count results
	// 1. For each result, Timeframe, GroupBy, Filter -> to one filter
	// 2. Match {1} -> Skip Percentile(Count, #percent) -> Limit 1
	return nil, nil
}

func median(p *model.Project, target string, params *QueryParams) (interface{}, error) {
	// return Percentile(target, 50)
	return nil, nil
}

func selectUnique(p *model.Project, target string, params *QueryParams) (interface{}, error) {
	// Unique w/ no counting
	return nil, nil
}
