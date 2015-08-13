package app

import (
	log "github.com/Sirupsen/logrus"
	"labix.org/v2/mgo"
)

const (
	RepoConfigPath = "config/repoconf.yml"
)

type repoConf struct {
	Open string
}

func initRepo(app *App) *mgo.Session {
	var configs map[string]repoConf
	if err := parseConfigYaml(RepoConfigPath, &configs); err != nil {
		log.Fatal(err.Error())
	}

	conf := configs[app.Env.Env]
	log.Printf("Repository initializing (%s)\n", conf)

	if session, err := mgo.Dial(conf.Open); err != nil {
		log.Fatal(err)
	} else {
		return session
	}

	return nil
}
