package api

import (
	"github.com/ant0ine/go-json-rest/rest"
	"net/http"
)

func HandleQueryCount(w rest.ResponseWriter, req *rest.Request) {
	if err := RequireReadKey(w, req); err != nil {
		rest.Error(w, err.Error(), err.(StatusError).Code)
		return
	}

	var params struct {
		*QueryParams
	}

	requires := []string{"event_collection"}
	if err := NewParams(req).Require(requires...).Parse(&params); err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	project := currentProject(req)
	result, err := count(project, params.QueryParams)

	if err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
	} else {
		w.WriteJson(result)
	}
}

func HandleQueryUniqueCount(w rest.ResponseWriter, req *rest.Request) {
	if err := RequireReadKey(w, req); err != nil {
		rest.Error(w, err.Error(), err.(StatusError).Code)
		return
	}

	var params struct {
		*QueryParams
		TargetProperty string `json:"target_property"`
	}

	requires := []string{"event_collection", "target_property"}
	if err := NewParams(req).Require(requires...).Parse(&params); err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	project := currentProject(req)
	result, err := countUnique(project, params.TargetProperty, params.QueryParams)

	if err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
	} else {
		w.WriteJson(result)
	}
}

func HandleQueryMin(w rest.ResponseWriter, req *rest.Request) {
	if err := RequireReadKey(w, req); err != nil {
		rest.Error(w, err.Error(), err.(StatusError).Code)
		return
	}

	var params struct {
		*QueryParams
		TargetProperty string `json:"target_property"`
	}

	requires := []string{"event_collection", "target_property"}
	if err := NewParams(req).Require(requires...).Parse(&params); err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	project := currentProject(req)
	result, err := min(project, params.TargetProperty, params.QueryParams)

	if err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
	} else {
		w.WriteJson(result)
	}
}

func HandleQueryMax(w rest.ResponseWriter, req *rest.Request) {
	if err := RequireReadKey(w, req); err != nil {
		rest.Error(w, err.Error(), err.(StatusError).Code)
		return
	}

	var params struct {
		*QueryParams
		TargetProperty string `json:"target_property"`
	}

	requires := []string{"event_collection", "target_property"}
	if err := NewParams(req).Require(requires...).Parse(&params); err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	project := currentProject(req)
	result, err := max(project, params.TargetProperty, params.QueryParams)

	if err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
	} else {
		w.WriteJson(result)
	}
}

func HandleQuerySum(w rest.ResponseWriter, req *rest.Request) {
	if err := RequireReadKey(w, req); err != nil {
		rest.Error(w, err.Error(), err.(StatusError).Code)
		return
	}

	var params struct {
		*QueryParams
		TargetProperty string `json:"target_property"`
	}

	requires := []string{"event_collection", "target_property"}
	if err := NewParams(req).Require(requires...).Parse(&params); err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	project := currentProject(req)
	result, err := sum(project, params.TargetProperty, params.QueryParams)

	if err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
	} else {
		w.WriteJson(result)
	}
}

func HandleQueryAverage(w rest.ResponseWriter, req *rest.Request) {
	if err := RequireReadKey(w, req); err != nil {
		rest.Error(w, err.Error(), err.(StatusError).Code)
		return
	}

	var params struct {
		*QueryParams
		TargetProperty string `json:"target_property"`
	}

	requires := []string{"event_collection", "target_property"}
	if err := NewParams(req).Require(requires...).Parse(&params); err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	project := currentProject(req)
	result, err := average(project, params.TargetProperty, params.QueryParams)

	if err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
	} else {
		w.WriteJson(result)
	}
}

func HandleQueryPercentile(w rest.ResponseWriter, req *rest.Request) {
	if err := RequireReadKey(w, req); err != nil {
		rest.Error(w, err.Error(), err.(StatusError).Code)
		return
	}

	var params struct {
		*QueryParams
		TargetProperty string  `json:"target_property"`
		Percent        float64 `json:"percent"`
	}

	requires := []string{"event_collection", "target_property", "percentile"}
	if err := NewParams(req).Require(requires...).Parse(&params); err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	project := currentProject(req)
	result, err := percentile(project, params.TargetProperty, params.Percent, params.QueryParams)

	if err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
	} else {
		w.WriteJson(result)
	}
}

func HandleQueryMedian(w rest.ResponseWriter, req *rest.Request) {
	if err := RequireReadKey(w, req); err != nil {
		rest.Error(w, err.Error(), err.(StatusError).Code)
		return
	}

	var params struct {
		*QueryParams
		TargetProperty string `json:"target_property"`
	}

	requires := []string{"event_collection", "target_property"}
	if err := NewParams(req).Require(requires...).Parse(&params); err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	project := currentProject(req)
	result, err := median(project, params.TargetProperty, params.QueryParams)

	if err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
	} else {
		w.WriteJson(result)
	}
}

func HandleQuerySelectUnique(w rest.ResponseWriter, req *rest.Request) {
	if err := RequireReadKey(w, req); err != nil {
		rest.Error(w, err.Error(), err.(StatusError).Code)
		return
	}

	var params struct {
		*QueryParams
		TargetProperty string `json:"target_property"`
	}

	requires := []string{"event_collection", "target_property", "timeframe"}
	if err := NewParams(req).Require(requires...).Parse(&params); err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	project := currentProject(req)
	result, err := selectUnique(project, params.TargetProperty, params.QueryParams)

	if err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
	} else {
		w.WriteJson(result)
	}
}
