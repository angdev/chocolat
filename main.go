package main

import (
	"github.com/angdev/chocolat/api"
	"github.com/ant0ine/go-json-rest/rest"
	"log"
	"net/http"
)

func main() {
	apiServer := rest.NewApi()
	apiServer.Use(rest.DefaultDevStack...)
	apiServer.Use(&rest.CorsMiddleware{
		RejectNonCorsRequests: false,
		OriginValidator: func(origin string, request *rest.Request) bool {
			return true
		},
		AllowedMethods: []string{"GET", "POST", "PUT"},
		AllowedHeaders: []string{
			"Authorization", "Accept", "Content-Type", "X-Custom-Header", "Origin",
		},
		AccessControlAllowCredentials: true,
		AccessControlMaxAge:           3600,
	})

	routes := mergeRouteSet(api.EventsRoutes, api.QueriesRoutes)

	router, err := rest.MakeRouter(routes...)
	if err != nil {
		log.Fatal(err)
	}

	apiServer.SetApp(router)
	log.Fatal(http.ListenAndServe(":5000", apiServer.MakeHandler()))
}

func mergeRouteSet(routeSets ...[]*rest.Route) []*rest.Route {
	var routes []*rest.Route
	for _, routeSet := range routeSets {
		routes = append(routes, routeSet...)
	}
	return routes
}
