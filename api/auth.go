package api

import (
	"github.com/angdev/chocolat/model"
	"github.com/ant0ine/go-json-rest/rest"
)

func requireApiKey(w rest.ResponseWriter, req *rest.Request, scopes ...model.ApiScope) error {
	project := currentProject(req)
	if project == nil {
		return NotFoundError
	}

	authKey, ok := authKeyValue(req)
	if !ok {
		return AuthKeyError
	}

	if !boundScope(project, authKey, scopes...) {
		return InvalidApiKeyError
	}

	return nil
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

func RequireReadKey(w rest.ResponseWriter, req *rest.Request) error {
	return requireApiKey(w, req, model.ApiReadKey, model.ApiMasterKey)
}

func RequireWriteKey(w rest.ResponseWriter, req *rest.Request) error {
	return requireApiKey(w, req, model.ApiWriteKey, model.ApiMasterKey)
}

func RequireMasterKey(w rest.ResponseWriter, req *rest.Request) error {
	return requireApiKey(w, req, model.ApiMasterKey)
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
