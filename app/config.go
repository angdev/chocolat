package app

import (
	"bytes"
	"gopkg.in/yaml.v2"
	"os"
	"path/filepath"
	"text/template"
)

const (
	RepoConfigPath = "config/repoconf.yml"
	DbConfigPath   = "config/dbconf.yml"
)

type repoConf struct {
	Open string
}

type dbConf struct {
	Driver string
	Open   string
}

func parseConfigYaml(path string, out interface{}) error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	configPath := filepath.Join(wd, path)
	tmp, err := template.ParseFiles(configPath)
	if err != nil {
		return err
	}

	var bytes bytes.Buffer
	if err = tmp.Execute(&bytes, envMap()); err != nil {
		return err
	}

	if err = yaml.Unmarshal(bytes.Bytes(), out); err != nil {
		return err
	}

	return nil
}
