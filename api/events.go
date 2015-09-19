package api

import (
	"github.com/ant0ine/go-json-rest/rest"
	"net/http"
)

func HandleInspectAllCollections(w rest.ResponseWriter, req *rest.Request) {
	if err := RequireReadKey(w, req); err != nil {
		rest.Error(w, err.Error(), err.(StatusError).Code)
		return
	}

	project := currentProject(req)
	result, err := inspectAllCollections(project)

	if err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
	} else {
		w.WriteJson(result)
	}
}

func HandleInspectAndCreateEvent(w rest.ResponseWriter, req *rest.Request) {
	if err := RequireReadKey(w, req); err != nil {
		// Try to handle with create event handler.
		// Same API url, but actions depend upon api key scope.
		HandleCreateEvent(w, req)
		return
	}

	project := currentProject(req)
	name := req.PathParam("event_name")
	if result, err := inspectCollection(project, name); err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
	} else {
		w.WriteJson(result)
	}
}

// Require a write key.
// Create a event.
func HandleCreateEvent(w rest.ResponseWriter, req *rest.Request) {
	if err := RequireWriteKey(w, req); err != nil {
		rest.Error(w, err.Error(), err.(StatusError).Code)
		return
	}

	project := currentProject(req)
	event := req.PathParam("event_name")

	var data CreateSingleEventParams
	var err error
	if err = eventData(req, &data); err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	events := make(map[string][]interface{})
	events[event] = []interface{}{data}

	result, err := createEvent(project, event, data)

	if err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
	} else {
		w.WriteJson(result)
	}
}

// Require a write key.
// Create multiple events with a single request.
func HandleCreateMultiEvents(w rest.ResponseWriter, req *rest.Request) {
	if err := RequireWriteKey(w, req); err != nil {
		rest.Error(w, err.Error(), err.(StatusError).Code)
		return
	}

	project := currentProject(req)
	var events CreateMultipleEventParams

	if err := req.DecodeJsonPayload(&events); err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result, err := createMultipleEvents(project, events)

	if err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
	} else {
		w.WriteJson(result)
	}
}

func eventData(req *rest.Request, out interface{}) error {
	encoded := req.FormValue("data")
	if encoded != "" {
		if err := decodeData(encoded, &out); err != nil {
			return err
		} else {
			return nil
		}
	}

	if err := req.DecodeJsonPayload(&out); err != nil {
		return err
	} else {
		return nil
	}
}
