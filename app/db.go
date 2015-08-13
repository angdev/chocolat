package app

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/jinzhu/gorm"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

const (
	DbConfigPath = "config/dbconf.yml"
)

type dbConf struct {
	Driver string
	Open   string
}

func (this *dbConf) String() string {
	return fmt.Sprintf("Driver: %s, Open: %s", this.Driver, this.Open)
}

func initDB(app *App) *gorm.DB {
	var configs map[string]dbConf
	if err := parseConfigYaml(DbConfigPath, &configs); err != nil {
		log.Fatal(err.Error())
	}

	conf := configs[app.Env.Env]
	log.Printf("Database initializing (%s)\n", conf)

	if opened, err := gorm.Open(conf.Driver, conf.Open); err != nil {
		log.Fatal(err.Error())
	} else {
		return &opened
	}
	return nil
}
