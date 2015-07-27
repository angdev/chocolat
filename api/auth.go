package api

import (
	"github.com/angdev/chocolat/model"
	"github.com/ant0ine/go-json-rest/rest"
	"net/http"
)

func requireApiKey(h rest.HandlerFunc, scopes ...model.ApiScope) rest.HandlerFunc {
	return func(w rest.ResponseWriter, req *rest.Request) {
		project := CurrentProject(req)
		if project == nil {
			rest.NotFound(w, req)
			return
		}

		authKey, ok := authKeyValue(req)
		if !ok {
			rest.Error(w, "Auth Key is missing", http.StatusBadRequest)
			return
		}

		if boundScope(project, authKey, scopes...) {
			h(w, req)
		} else {
			rest.Error(w, "Requires a valid api key", http.StatusUnauthorized)
		}
	}
}

func boundScope(project *model.Project, authKey string, scopes ...model.ApiScope) bool {
	ok := false

	for _, scope := range scopes {
		switch scope {
		case model.ApiReadKey:
			ok = (project.ReadKey().Value == authKey)
		case model.ApiWriteKey:
			ok = (project.WriteKey().Value == authKey)
		case model.ApiMasterKey:
			ok = (project.MasterKey().Value == authKey)
		}

		if ok {
			return ok
		}
	}

	return ok
}

func RequireReadKey(h rest.HandlerFunc) rest.HandlerFunc {
	return requireApiKey(h, model.ApiReadKey, model.ApiMasterKey)
}

func RequireWriteKey(h rest.HandlerFunc) rest.HandlerFunc {
	return requireApiKey(h, model.ApiWriteKey, model.ApiMasterKey)
}

func RequireMasterKey(h rest.HandlerFunc) rest.HandlerFunc {
	return requireApiKey(h, model.ApiMasterKey)
}

func authKeyValue(req *rest.Request) (string, bool) {
	key := req.FormValue("api_key")
	if key != "" {
		return key, true
	}

	key = req.Header.Get("Authorization")
	if key == "" {
		return "", false
	} else {
		return key, true
	}
}
