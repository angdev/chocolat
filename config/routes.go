package config

import (
	"github.com/angdev/chocolat/api"
	"github.com/angdev/chocolat/lib/routes"
	"github.com/ant0ine/go-json-rest/rest"
)

var (
	Routes []*rest.Route
	r      *routes.Builder = routes.NewBuilder()
)

func init() {
	r.Namespace("/3.0", func(r *routes.Builder) {
		r.Namespace("/projects/:project_id", func(r *routes.Builder) {
			r.Get("/events/:event_name", api.HandleInspectAndCreateEvent)
			r.Post("/events/:event_name", api.HandleCreateEvent)
			r.Get("/events", api.HandleInspectAllCollections)
			r.Post("/events", api.HandleCreateMultiEvents)
		})

		r.Namespace("/projects/:project_id/queries", func(r *routes.Builder) {
			r.Get("/count", api.HandleQueryCount)
			r.Post("/count", api.HandleQueryCount)
			r.Get("/count_unique", api.HandleQueryUniqueCount)
			r.Post("/count_unique", api.HandleQueryUniqueCount)
			r.Get("/minimum", api.HandleQueryMin)
			r.Post("/minimum", api.HandleQueryMin)
			r.Get("/maximum", api.HandleQueryMax)
			r.Post("/maximum", api.HandleQueryMax)
			r.Get("/sum", api.HandleQuerySum)
			r.Post("/sum", api.HandleQuerySum)
			r.Get("/average", api.HandleQueryAverage)
			r.Post("/average", api.HandleQueryAverage)
			r.Get("/percentile", api.HandleQueryPercentile)
			r.Post("/percentile", api.HandleQueryPercentile)
			r.Get("/median", api.HandleQueryMedian)
			r.Post("/median", api.HandleQueryMedian)
			r.Get("/select_unique", api.HandleQuerySelectUnique)
			r.Post("/select_unique", api.HandleQuerySelectUnique)
		})
	})

	Routes = r.Routes()
}
