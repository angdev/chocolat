package config

import (
	"github.com/angdev/chocolat/api"
	"github.com/ant0ine/go-json-rest/rest"
)

var Routes []*rest.Route

func init() {
	Routes = []*rest.Route{
		rest.Get("/projects/:project_id/events/:event_name", api.HandleCreateEvent),
		rest.Post("/projects/:project_id/events/:event_name", api.HandleCreateEvent),
		rest.Post("/projects/:project_id/events", api.HandleCreateMultiEvents),

		rest.Get("/projects/:project_id/queries/count", api.HandleQueryCount),
		rest.Post("/projects/:project_id/queries/count", api.HandleQueryCount),
		rest.Get("/projects/:project_id/queries/count_unique", api.HandleQueryUniqueCount),
		rest.Post("/projects/:project_id/queries/count_unique", api.HandleQueryUniqueCount),
		rest.Get("/projects/:project_id/queries/min", api.HandleQueryMin),
		rest.Post("/projects/:project_id/queries/min", api.HandleQueryMin),
		rest.Get("/projects/:project_id/queries/max", api.HandleQueryMax),
		rest.Post("/projects/:project_id/queries/max", api.HandleQueryMax),
		rest.Get("/projects/:project_id/queries/sum", api.HandleQuerySum),
		rest.Post("/projects/:project_id/queries/sum", api.HandleQuerySum),
		rest.Get("/projects/:project_id/queries/average", api.HandleQueryAverage),
		rest.Post("/projects/:project_id/queries/average", api.HandleQueryAverage),
		rest.Get("/projects/:project_id/queries/percentile", api.HandleQueryPercentile),
		rest.Post("/projects/:project_id/queries/percentile", api.HandleQueryPercentile),
		rest.Get("/projects/:project_id/queries/select_unique", api.HandleQuerySelectUnique),
		rest.Post("/projects/:project_id/queries/select_unique", api.HandleQuerySelectUnique),
	}
}
