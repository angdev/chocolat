package model

import (
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/kardianos/osext"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

const (
	ConfigPath = "db/dbconf.yml"
)

var db *gorm.DB

func DB() *gorm.DB {
	if db == nil {
		log.Fatal("Database is not initialized!")
	}
	return db
}

type dbConf struct {
	Driver string
	Open   string
}

func (this *dbConf) String() string {
	return fmt.Sprintf("Driver: %s, Open: %s", this.Driver, this.Open)
}

func InitDB() {
	env := envOrDefault()
	conf, err := dbConfEnv(env)
	if err != nil {
		log.Fatal(err.Error())
	}

	log.Printf("Database initializing (%s)", conf)

	if opened, err := gorm.Open(conf.Driver, conf.Open); err != nil {
		log.Fatal(err.Error())
	} else {
		db = &opened
	}
}

func envOrDefault() string {
	env := os.Getenv("CHOCOLAT_ENV")
	if env == "" {
		env = "development"
	}
	return env
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
