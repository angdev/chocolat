package app

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/angdev/chocolat/api"
	"github.com/angdev/chocolat/model"
	"github.com/ant0ine/go-json-rest/rest"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"net/http"
)

var Chocolat = App{}

type Environment struct {
	Env  string
	Port string
}

type App struct {
	Env Environment
}

func (this *App) Init() {
	this.initEnv()
	this.initModel()
}

func (this *App) Run() {
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
		log.Fatal(err)
	}

	apiServer.SetApp(router)

	log.Info("Start serving, port=", this.Env.Port)

	log.Fatal(http.ListenAndServe(this.Port(), apiServer.MakeHandler()))
}

func mergeRouteSet(routeSets ...[]*rest.Route) []*rest.Route {
	var routes []*rest.Route
	for _, routeSet := range routeSets {
		routes = append(routes, routeSet...)
	}
	return routes
}

func (this *App) Port() string {
	return fmt.Sprintf(":%s", this.Env.Port)
}

func (this *App) initEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Print(err.Error())
	}

	this.Env = this.defaultEnv()
	if err := envconfig.Process("chocolat", &this.Env); err != nil {
		log.Fatal("Environment parse error")
	}
}

func (this *App) defaultEnv() Environment {
	return Environment{
		Env:  "development",
		Port: "5000",
	}
}

func (this *App) initModel() {
	db := initDB(this)
	model.Init(db)
}
