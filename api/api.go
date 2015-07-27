package api

import (
	"github.com/angdev/chocolat/model"
	"github.com/ant0ine/go-json-rest/rest"
)

func Routes(routes ...*rest.Route) []*rest.Route {
	return routes
}

func CurrentProject(req *rest.Request) *model.Project {
	uuid := req.PathParam("project_id")
	return model.ProjectByUUID(uuid)
}
