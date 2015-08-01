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
)

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

	var params QueryParams
	if err := req.DecodeJsonPayload(&params); err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	ensureEventCollection(req, &params)

	var target string
	result, err := countUnique(project, target, &params)

	if err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
	} else {
		w.WriteJson(result)
	}
}

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
