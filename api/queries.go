package api

import (
	"errors"
	"github.com/angdev/chocolat/service"
	"github.com/ant0ine/go-json-rest/rest"
	"github.com/k0kubun/pp"
	"net/http"
)

var QueriesRoutes = Routes(
	rest.Get("/projects/:project_id/queries/count", RequireReadKey(queryCount)),
	rest.Post("/projects/:project_id/queries/count", RequireReadKey(queryCount)),
)

func queryCount(w rest.ResponseWriter, req *rest.Request) {
	project := CurrentProject(req)

	var params service.CountParams
	if err := req.DecodeJsonPayload(&params); err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	ensureEventCollection(req, &params.QueryParams)

	pp.Println(params)

	result, err := service.Count(project, &params)

	if err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
	} else {
		w.WriteJson(result)
	}
}

func ensureEventCollection(req *rest.Request, params *service.QueryParams) error {
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
