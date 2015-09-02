package app

import (
	"fmt"
	log "github.com/Sirupsen/logrus"

	"github.com/angdev/chocolat/config"
	"github.com/angdev/chocolat/model"
	"github.com/angdev/chocolat/repo"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"

	"labix.org/v2/mgo"

	"github.com/ant0ine/go-json-rest/rest"
	"net/http"

	"github.com/jinzhu/gorm"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

var Chocolat = App{}

type App struct {
	Env Environment
}

func (this *App) Init() {
	this.initEnv()
	this.initModel()
	this.initRepo()
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

	router, err := rest.MakeRouter(config.Routes...)
	if err != nil {
		log.Fatal(err)
	}

	apiServer.SetApp(router)

	log.Info("Start serving, port=", this.Env.Port)
	log.Fatal(http.ListenAndServe(this.Port(), apiServer.MakeHandler()))
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

func (this *App) initModel() {
	var configs map[string]dbConf
	if err := parseConfigYaml(DbConfigPath, &configs); err != nil {
		log.Fatal(err.Error())
	}

	conf := configs[this.Env.Env]
	log.Debugf("Database initializing (%s)\n", conf)

	if opened, err := gorm.Open(conf.Driver, conf.Open); err != nil {
		log.Fatal(err.Error())
	} else {
		model.Init(&opened)
	}
}

func (this *App) initRepo() {
	var configs map[string]repoConf
	if err := parseConfigYaml(RepoConfigPath, &configs); err != nil {
		log.Fatal(err.Error())
	}

	conf := configs[this.Env.Env]
	log.Debugf("Repository initializing (%s)\n", conf)

	if session, err := mgo.Dial(conf.Open); err != nil {
		log.Fatal(err)
	} else {
		repo.Init(session)
	}
}
