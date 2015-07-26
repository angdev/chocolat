package model

import (
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/k0kubun/pp"
	"github.com/kardianos/osext"
	"gopkg.in/yaml.v2"
	"io/ioutil"
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
		pp.Fatal("Database is not initialized!")
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
		pp.Fatal(err.Error())
	}

	pp.Printf("Database initializing (%s)\n", conf)

	if opened, err := gorm.Open(conf.Driver, conf.Open); err != nil {
		pp.Fatal(err.Error())
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
		pp.Fatal("Cannot get executable directory path!")
	}

	configPath := filepath.Join(wd, ConfigPath)
	bytes, err := ioutil.ReadFile(configPath)
	if err != nil {
		pp.Fatal("Cannot read database config file! - " + configPath)
	}

	var conf map[string]*dbConf
	err = yaml.Unmarshal(bytes, &conf)
	if err != nil {
		pp.Fatal("Cannot unmarshal config file!")
	}

	if value, ok := conf[env]; ok {
		return value, nil
	} else {
		return nil, errors.New("Unknown environment")
	}
}
