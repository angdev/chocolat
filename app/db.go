package app

import (
	"bytes"
	"errors"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/jinzhu/gorm"
	"github.com/kardianos/osext"
	"gopkg.in/yaml.v2"
	"os"
	"path/filepath"
	"strings"
	"text/template"

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
		log.Fatal(err)
	}

	configPath := filepath.Join(wd, ConfigPath)
	tmp, err := template.ParseFiles(configPath)
	if err != nil {
		log.Fatal(err)
	}

	var bytes bytes.Buffer
	if err = tmp.Execute(&bytes, envMap()); err != nil {
		log.Fatal(err)
	}

	var conf map[string]*dbConf
	if err = yaml.Unmarshal(bytes.Bytes(), &conf); err != nil {
		log.Fatal(err)
	}

	if value, ok := conf[env]; ok {
		return value, nil
	} else {
		return nil, errors.New("Unknown environment")
	}
}

func envMap() map[string]string {
	environ := os.Environ()
	env := make(map[string]string)

	for _, v := range environ {
		pair := strings.Split(v, "=")
		env[pair[0]] = pair[1]
	}

	return env
}
