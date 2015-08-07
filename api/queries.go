package api

import (
	"errors"
	"github.com/ant0ine/go-json-rest/rest"
	"net/http"
)

var QueriesRoutes = Routes(
	rest.Get("/projects/:project_id/queries/count", RequireReadKey(handleQueryCount)),
	rest.Post("/projects/:project_id/queries/count", RequireReadKey(handleQueryCount)),
	rest.Get("/projects/:project_id/queries/count_unique", RequireReadKey(handleQueryUniqueCount)),
	rest.Post("/projects/:project_id/queries/count_unique", RequireReadKey(handleQueryUniqueCount)),
	rest.Get("/projects/:project_id/queries/min", RequireReadKey(handleQueryMin)),
	rest.Post("/projects/:project_id/queries/min", RequireReadKey(handleQueryMin)),
	rest.Get("/projects/:project_id/queries/max", RequireReadKey(handleQueryMax)),
	rest.Post("/projects/:project_id/queries/max", RequireReadKey(handleQueryMax)),
	rest.Get("/projects/:project_id/queries/sum", RequireReadKey(handleQuerySum)),
	rest.Post("/projects/:project_id/queries/sum", RequireReadKey(handleQuerySum)),
	rest.Get("/projects/:project_id/queries/average", RequireReadKey(handleQueryAverage)),
	rest.Post("/projects/:project_id/queries/average", RequireReadKey(handleQueryAverage)),
	rest.Get("/projects/:project_id/queries/percentile", RequireReadKey(handleQueryPercentile)),
	rest.Post("/projects/:project_id/queries/percentile", RequireReadKey(handleQueryPercentile)),
)

func ensureEventCollection(req *rest.Request, params *QueryParams) error {
	if params.CollectionName != "" {
		return nil
	}

	if v := req.FormValue("event_collection"); v != "" {
		params.CollectionName = v
		return nil
	} else {
		return errors.New("event_collection is missing")
	}
}

func handleQueryCount(w rest.ResponseWriter, req *rest.Request) {
	project := CurrentProject(req)

	var params QueryParams
	if err := req.DecodeJsonPayload(&params); err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	ensureEventCollection(req, &params)

	result, err := count(project, &params)

	if err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
	} else {
		w.WriteJson(result)
	}
}

func handleQueryUniqueCount(w rest.ResponseWriter, req *rest.Request) {
	project := CurrentProject(req)

	var params struct {
		*QueryParams
		TargetProperty string `json:"target_property"`
	}

	if err := req.DecodeJsonPayload(&params); err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	ensureEventCollection(req, params.QueryParams)

	result, err := countUnique(project, params.TargetProperty, params.QueryParams)

	if err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
	} else {
		w.WriteJson(result)
	}
}

func handleQueryMin(w rest.ResponseWriter, req *rest.Request) {
	project := CurrentProject(req)

	var params struct {
		*QueryParams
		TargetProperty string `json:"target_property"`
	}

	if err := req.DecodeJsonPayload(&params); err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	ensureEventCollection(req, params.QueryParams)

	result, err := min(project, params.TargetProperty, params.QueryParams)

	if err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
	} else {
		w.WriteJson(result)
	}
}

func handleQueryMax(w rest.ResponseWriter, req *rest.Request) {
	project := CurrentProject(req)

	var params struct {
		*QueryParams
		TargetProperty string `json:"target_property"`
	}

	if err := req.DecodeJsonPayload(&params); err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	ensureEventCollection(req, params.QueryParams)

	result, err := max(project, params.TargetProperty, params.QueryParams)

	if err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
	} else {
		w.WriteJson(result)
	}
}

func handleQuerySum(w rest.ResponseWriter, req *rest.Request) {
	project := CurrentProject(req)

	var params struct {
		*QueryParams
		TargetProperty string `json:"target_property"`
	}

	if err := req.DecodeJsonPayload(&params); err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	ensureEventCollection(req, params.QueryParams)

	result, err := sum(project, params.TargetProperty, params.QueryParams)

	if err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
	} else {
		w.WriteJson(result)
	}
}

func handleQueryAverage(w rest.ResponseWriter, req *rest.Request) {
	project := CurrentProject(req)

	var params struct {
		*QueryParams
		TargetProperty string `json:"target_property"`
	}

	if err := req.DecodeJsonPayload(&params); err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	ensureEventCollection(req, params.QueryParams)

	result, err := average(project, params.TargetProperty, params.QueryParams)

	if err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
	} else {
		w.WriteJson(result)
	}
}

func handleQueryPercentile(w rest.ResponseWriter, req *rest.Request) {
	project := CurrentProject(req)

	var params struct {
		*QueryParams
		TargetProperty string  `json:"target_property"`
		Percent        float64 `json:"percent"`
	}

	if err := req.DecodeJsonPayload(&params); err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	ensureEventCollection(req, params.QueryParams)

	result, err := percentile(project, params.TargetProperty, params.Percent, params.QueryParams)

	if err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
	} else {
		w.WriteJson(result)
	}
}
