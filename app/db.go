package app

import (
	"errors"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/jinzhu/gorm"
	"github.com/kardianos/osext"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"path/filepath"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

const (
	ConfigPath = "db/dbconf.yml"
)

type dbConf struct {
	Driver string
	Open   string
}

func (this *dbConf) String() string {
	return fmt.Sprintf("Driver: %s, Open: %s", this.Driver, this.Open)
}

func initDB(app *App) *gorm.DB {
	conf, err := dbConfEnv(app.Env.Env)
	if err != nil {
		log.Fatal(err.Error())
	}

	log.Printf("Database initializing (%s)\n", conf)

	if opened, err := gorm.Open(conf.Driver, conf.Open); err != nil {
		log.Fatal(err.Error())
	} else {
		return &opened
	}
	return nil
}

func dbConfEnv(env string) (*dbConf, error) {
	wd, err := osext.ExecutableFolder()
	if err != nil {
		log.Fatal("Cannot get executable directory path!")
	}

	configPath := filepath.Join(wd, ConfigPath)
	bytes, err := ioutil.ReadFile(configPath)
	if err != nil {
		log.Fatal("Cannot read database config file! - " + configPath)
	}

	var conf map[string]*dbConf
	err = yaml.Unmarshal(bytes, &conf)
	if err != nil {
		log.Fatal("Cannot unmarshal config file!")
	}

	if value, ok := conf[env]; ok {
		return value, nil
	} else {
		return nil, errors.New("Unknown environment")
	}
}
