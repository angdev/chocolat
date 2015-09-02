package routes

import (
	log "github.com/Sirupsen/logrus"
	"github.com/ant0ine/go-json-rest/rest"
	"strings"
)

type Collector struct {
	Routes []*rest.Route
}

func NewCollector() *Collector {
	return &Collector{Routes: []*rest.Route{}}
}

func (this *Collector) collect(routes ...*rest.Route) {
	this.Routes = append(this.Routes, routes...)
}

type Builder struct {
	RouteStack []string
	collector  *Collector
}

func NewBuilder() *Builder {
	return &Builder{RouteStack: []string{}, collector: NewCollector()}
}

func (this *Builder) Routes() []*rest.Route {
	return this.collector.Routes
}

func (this *Builder) Namespace(ns string, block func(*Builder)) {
	builder := this.clone()
	builder.RouteStack = append(builder.RouteStack, ns)
	block(builder)
}

func (this *Builder) Get(path string, handler rest.HandlerFunc) {
	builder := this.clone().appendPath(path)
	this.collector.collect(rest.Get(builder.joinedPath(), handler))
}

func (this *Builder) Post(path string, handler rest.HandlerFunc) {
	builder := this.clone().appendPath(path)
	this.collector.collect(rest.Post(builder.joinedPath(), handler))
}

func (this *Builder) Method(methods []string, path string, handler rest.HandlerFunc) {
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

func (this *Builder) joinedPath() string {
	return strings.Join(this.RouteStack, "")
}

func (this *Builder) clone() *Builder {
	return &Builder{RouteStack: this.RouteStack, collector: this.collector}
}

func (this *Builder) appendPath(path string) *Builder {
	this.RouteStack = append(this.RouteStack, path)
	return this
}
