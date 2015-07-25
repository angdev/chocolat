package api

import (
	"github.com/angdev/chocolat/service"
	"github.com/angdev/chocolat/support/repo"
	"github.com/ant0ine/go-json-rest/rest"
	"net/http"
)

var QueriesRoutes = Routes(
	rest.Get("/projects/:project_id/queries/count", queryCountGet),
	rest.Post("/projects/:project_id/queries/count", queryCountPost),
)

func queryCountGet(w rest.ResponseWriter, req *rest.Request) {
	dbName := req.PathParam("project_id")
	collName := req.FormValue("event_collection")

	result, err := service.Count(dbName, &service.CountParams{
		CollectionName: collName,
	})

	if err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
	} else {
		w.WriteJson(&repo.Doc{"result": result})
	}
}

func queryCountPost(w rest.ResponseWriter, req *rest.Request) {
	dbName := req.PathParam("project_id")
	payload := repo.Doc{}

	if err := req.DecodeJsonPayload(&payload); err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	collName := payload["event_collection"].(string)
	result, err := service.Count(dbName, &service.CountParams{
		CollectionName: collName,
	})

	if err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
	} else {
		w.WriteJson(result)
	}
}
