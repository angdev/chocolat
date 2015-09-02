package config

import (
	log "github.com/Sirupsen/logrus"
	"github.com/angdev/chocolat/api"
	"github.com/ant0ine/go-json-rest/rest"
	"strings"
)

var (
	Routes []*rest.Route
	r      *RouteBuilder = &RouteBuilder{}
)

func init() {
	r.Namespace("/3.0", func(r *RouteBuilder) {
		r.Namespace("/projects/:project_id", func(r *RouteBuilder) {
			r.Get("/events/:event_name", api.HandleCreateEvent)
			r.Post("/events/:event_name", api.HandleCreateEvent)
			r.Get("/events", api.HandleCreateMultiEvents)
		})

		r.Namespace("/projects/:project_id/queries", func(r *RouteBuilder) {
			r.Get("/count", api.HandleQueryCount)
			r.Post("/count", api.HandleQueryCount)
			r.Get("/count_unique", api.HandleQueryUniqueCount)
			r.Post("/count_unique", api.HandleQueryUniqueCount)
			r.Get("/min", api.HandleQueryMin)
			r.Post("/min", api.HandleQueryMin)
			r.Get("/max", api.HandleQueryMax)
			r.Post("/max", api.HandleQueryMax)
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
}

type RouteBuilder struct {
	RouteStack []string
}

func (this *RouteBuilder) joinedPath() string {
	return strings.Join(this.RouteStack, "")
}

func (this *RouteBuilder) clone() *RouteBuilder {
	return &RouteBuilder{RouteStack: this.RouteStack}
}

func (this *RouteBuilder) appendPath(path string) *RouteBuilder {
	this.RouteStack = append(this.RouteStack, path)
	return this
}

func (this *RouteBuilder) Namespace(ns string, block func(*RouteBuilder)) {
	block(&RouteBuilder{RouteStack: append(this.RouteStack, ns)})
}

func (this *RouteBuilder) Get(path string, handler rest.HandlerFunc) {
	builder := this.clone().appendPath(path)
	addRoutes(rest.Get(builder.joinedPath(), handler))
}

func (this *RouteBuilder) Post(path string, handler rest.HandlerFunc) {
	builder := this.clone().appendPath(path)
	addRoutes(rest.Post(builder.joinedPath(), handler))
}

func (this *RouteBuilder) Method(methods []string, path string, handler rest.HandlerFunc) {
	for _, method := range methods {
		method = strings.ToUpper(method)
		switch method {
		case "GET":
			this.Get(path, handler)
		case "POST":
			this.Post(path, handler)
		default:
			log.Fatal("Not supported method")
		}
	}
}

func addRoutes(routes ...*rest.Route) {
	Routes = append(Routes, routes...)
}
