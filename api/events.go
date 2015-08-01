package api

import (
	"github.com/angdev/chocolat/support"
	"github.com/angdev/chocolat/support/repo"
	"github.com/ant0ine/go-json-rest/rest"
	"net/http"
)

// Including routes related with events.
var EventsRoutes = Routes(
	rest.Get("/projects/:project_id/events/:event_name", RequireWriteKey(handleCreateEvent)),
	rest.Post("/projects/:project_id/events/:event_name", RequireWriteKey(handleCreateEvent)),
	rest.Post("/projects/:project_id/events", RequireWriteKey(handleCreateMultiEvents)),
)

// Require a write key.
// Create a event.
func handleCreateEvent(w rest.ResponseWriter, req *rest.Request) {
	project := CurrentProject(req)
	event := req.PathParam("event_name")

	var data repo.Doc
	var err error
	if data, err = eventData(req); err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	events := map[string][]repo.Doc{}
	events[event] = []repo.Doc{data}

	result, err := createEvent(project, &CreateEventParams{
		CollectionName: event,
		Events:         events,
	})

	if err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
	} else {
		w.WriteJson(result)
	}
}

func eventData(req *rest.Request) (repo.Doc, error) {
	encoded := req.FormValue("data")
	if encoded != "" {
		if doc, err := support.DecodeData(encoded); err != nil {
			return nil, err
		} else {
			return doc, nil
		}
	}

	var payload repo.Doc
	if err := req.DecodeJsonPayload(&payload); err != nil {
		return nil, err
	} else {
		return payload, nil
	}
}

// Require a write key.
// Create multiple events with a single request.
func handleCreateMultiEvents(w rest.ResponseWriter, req *rest.Request) {
	project := CurrentProject(req)
	event := req.PathParam("event_name")
	var events map[string][]repo.Doc

	if err := req.DecodeJsonPayload(&events); err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result, err := createMultipleEvents(project, &CreateEventParams{
		CollectionName: event,
		Events:         events,
	})

	if err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
	} else {
		w.WriteJson(result)
	}
}
