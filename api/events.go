package api

import (
	"github.com/ant0ine/go-json-rest/rest"
	"net/http"
)

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
