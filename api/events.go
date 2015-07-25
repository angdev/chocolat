package api

import (
	"github.com/angdev/chocolat/service"
	"github.com/angdev/chocolat/support/repo"
	"github.com/ant0ine/go-json-rest/rest"
	"net/http"
)

var EventsRoutes = Routes(
	rest.Get("/projects/:project_id/events/:event_name", createEvent),
)

func createEvent(w rest.ResponseWriter, req *rest.Request) {
	dbName := req.PathParam("project_id")
	collName := req.PathParam("event_name")
	encodedData := req.FormValue("data")

	_, err := service.CreateEvent(dbName, &service.CreateEventParams{
		CollectionName: collName,
		EncodedData:    encodedData,
	})

	if err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
	} else {
		w.WriteJson(&repo.Doc{"created": true})
	}
}
