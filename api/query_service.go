package api

import (
	"github.com/angdev/chocolat/lib/query"
	"github.com/angdev/chocolat/model"
	"github.com/angdev/chocolat/repo"
)

func count(p *model.Project, params *QueryParams) (interface{}, error) {
	var aggregator = func(params *QueryParams, out interface{}) error {
		r := repo.NewRepository(p.RepoName())
		defer r.Close()

		q := query.New(r.C(params.CollectionName), params.ToQuery().Count())
		return q.Execute(out)
	}

	presenter := NewPresenter(aggregator, params)
	return presenter.Present()
}

func countUnique(p *model.Project, target string, params *QueryParams) (interface{}, error) {
	var aggregator = func(params *QueryParams, out interface{}) error {
		r := repo.NewRepository(p.RepoName())
		defer r.Close()

		q := query.New(r.C(params.CollectionName), params.ToQuery().CountUnique(target))
		return q.Execute(out)
	}

	presenter := NewPresenter(aggregator, params)
	return presenter.Present()
}

func min(p *model.Project, target string, params *QueryParams) (interface{}, error) {
	var aggregator = func(params *QueryParams, out interface{}) error {
		r := repo.NewRepository(p.RepoName())
		defer r.Close()

		q := query.New(r.C(params.CollectionName), params.ToQuery().Min(target))
		return q.Execute(out)
	}

	presenter := NewPresenter(aggregator, params)
	return presenter.Present()
}

func max(p *model.Project, target string, params *QueryParams) (interface{}, error) {
	var aggregator = func(params *QueryParams, out interface{}) error {
		r := repo.NewRepository(p.RepoName())
		defer r.Close()

		q := query.New(r.C(params.CollectionName), params.ToQuery().Max(target))
		return q.Execute(out)
	}

	presenter := NewPresenter(aggregator, params)
	return presenter.Present()
}

func sum(p *model.Project, target string, params *QueryParams) (interface{}, error) {
	var aggregator = func(params *QueryParams, out interface{}) error {
		r := repo.NewRepository(p.RepoName())
		defer r.Close()

		q := query.New(r.C(params.CollectionName), params.ToQuery().Sum(target))
		return q.Execute(out)
	}

	presenter := NewPresenter(aggregator, params)
	return presenter.Present()
}

func average(p *model.Project, target string, params *QueryParams) (interface{}, error) {
	var aggregator = func(params *QueryParams, out interface{}) error {
		r := repo.NewRepository(p.RepoName())
		defer r.Close()

		q := query.New(r.C(params.CollectionName), params.ToQuery().Average(target))
		return q.Execute(out)
	}

	presenter := NewPresenter(aggregator, params)
	return presenter.Present()
}

func percentile(p *model.Project, target string, percent float64, params *QueryParams) (interface{}, error) {
	var aggregator = func(params *QueryParams, out interface{}) error {
		r := repo.NewRepository(p.RepoName())
		defer r.Close()

		q := query.New(r.C(params.CollectionName), params.ToQuery().
			OrderBy(&query.Order{Field: target, Order: query.ASC}).Collect(target))

		if err := q.Execute(out); err != nil {
			return err
		}

		switch out.(type) {
		case *queryGroupResultArray:
			results := *out.(*queryGroupResultArray)
			for i, _ := range results {
				result := results[i].Result.([]interface{})
				offset := int(float64(len(result)) * percent / 100)
				results[i].Result = result[offset]
			}
		case *queryResult:
			result := out.(*queryResult)
			results := result.Result.([]interface{})
			offset := int(float64(len(results)) * percent / 100)
			result.Result = results[offset]
		}

		return nil
	}

	presenter := NewPresenter(aggregator, params)
	return presenter.Present()
}

func median(p *model.Project, target string, params *QueryParams) (interface{}, error) {
	return percentile(p, target, 50, params)
}

func selectUnique(p *model.Project, target string, params *QueryParams) (interface{}, error) {
	var aggregator = func(params *QueryParams, out interface{}) error {
		r := repo.NewRepository(p.RepoName())
		defer r.Close()

		q := query.New(r.C(params.CollectionName), params.ToQuery().SelectUnique(target))
		return q.Execute(out)
	}

	presenter := NewPresenter(aggregator, params)
	return presenter.Present()
}
