package main

import (
	"github.com/k0kubun/pp"
	"github.com/angdev/chocolat/api"
	"github.com/angdev/chocolat/model"
	"github.com/ant0ine/go-json-rest/rest"
	"github.com/joho/godotenv"
	"net/http"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		pp.Print(err.Error())
	}

	model.InitDB()

	db := model.DB()
	db.LogMode(true)

	apiServer := rest.NewApi()
	apiServer.Use(rest.DefaultDevStack...)

	// Cors
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

	// Jsonp
	apiServer.Use(&rest.JsonpMiddleware{
		CallbackNameKey: "jsonp",
	})

	routes := mergeRouteSet(api.EventsRoutes, api.QueriesRoutes)

	router, err := rest.MakeRouter(routes...)
	if err != nil {
		pp.Fatal(err)
	}

	apiServer.SetApp(router)
	pp.Fatal(http.ListenAndServe(":5000", apiServer.MakeHandler()))
}

func mergeRouteSet(routeSets ...[]*rest.Route) []*rest.Route {
	var routes []*rest.Route
	for _, routeSet := range routeSets {
		routes = append(routes, routeSet...)
	}
	return routes
}
