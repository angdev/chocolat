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

func listProject() {
	initialize()

	db := model.DB()
	var projects []model.Project

	if db.Find(&projects).RecordNotFound() {
		fmt.Println("No project found")
		return
	}

	for _, p := range projects {
		fmt.Println(p.UUID)
	}
}

func inspectProject(uuid string) {
	initialize()

	project := model.ProjectByUUID(uuid)

	if project == nil {
		fmt.Println("No project found")
		return
	}

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
		cli.BoolFlag{
			Name:  "create, c",
			Usage: "Create a new project",
		},
		cli.BoolFlag{
			Name:  "run, r",
			Usage: "Run the server",
		},
		cli.BoolFlag{
			Name:  "list, l",
			Usage: "Listing projects",
		},
		cli.StringFlag{
			Name:  "project, p",
			Usage: "Inspect a project",
		},
	}

	app.Action = func(c *cli.Context) {
		if c.Bool("create") {
			createProject()
		} else if c.Bool("list") {
			listProject()
		} else if c.String("project") != "" {
			inspectProject(c.String("project"))
		} else if c.Bool("run") {
			run()
		} else {
			cli.ShowAppHelp(c)
		}
	}

	app.Run(os.Args)
}
