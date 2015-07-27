package api

import (
	"github.com/angdev/chocolat/service"
	"github.com/angdev/chocolat/support/repo"
	"github.com/ant0ine/go-json-rest/rest"
	"net/http"
)

var QueriesRoutes = Routes(
	rest.Get("/projects/:project_id/queries/count", RequireReadKey(queryCount)),
	rest.Post("/projects/:project_id/queries/count", RequireReadKey(queryCount)),
)

func queryCount(w rest.ResponseWriter, req *rest.Request) {
	project := CurrentProject(req)

	var payload repo.Doc
	if err := req.DecodeJsonPayload(&payload); err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	collName := eventCollection(req, payload)
	params, err := service.NewCountParams(collName, payload)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result, err := service.Count(project, params)

	if err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
	} else {
		w.WriteJson(result)
	}
}

func eventCollection(req *rest.Request, payload repo.Doc) string {
	if v, ok := payload["event_collection"]; ok {
		return v.(string)
	}

	v := req.FormValue("event_collection")
	return v
}
