package main

import (
	"fmt"
	"github.com/angdev/chocolat/api"
	"github.com/angdev/chocolat/model"
	"github.com/ant0ine/go-json-rest/rest"
	"github.com/codegangsta/cli"
	"github.com/joho/godotenv"
	"github.com/k0kubun/pp"
	"github.com/satori/go.uuid"
	"net/http"
	"os"
)

func initialize() {
	err := godotenv.Load()
	if err != nil {
		pp.Print(err.Error())
	}
	model.InitDB()
}

func run() {
	initialize()

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

func createProject() {
	initialize()

	db := model.DB()
	project := model.Project{
		UUID: uuid.NewV4().String(),
	}
	db.Create(&project)

	fmt.Printf("Created a new project.\n")
	fmt.Printf("Project UUID - %s\n", project.UUID)
	fmt.Printf("Project Master Key - %s\n", project.MasterKey().Value)
	fmt.Printf("Project Read Key - %s\n", project.ReadKey().Value)
	fmt.Printf("Project Write Key - %s\n", project.WriteKey().Value)
}

func main() {
	app := cli.NewApp()
	app.Name = "chocolat"
	app.Usage = "Yet Another Data Aggregation Server"
	app.Version = "0.0.1"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "create, c",
			Value: "project",
			Usage: "Create a new project",
		},
	}

	app.Action = func(c *cli.Context) {
		if c.String("create") == "project" {
			createProject()
		} else {
			run()
		}
	}

	app.Run(os.Args)
}
