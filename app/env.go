package app

import (
	"os"
	"strings"
)

type Environment struct {
	Env  string
	Port string
}

func (this *App) defaultEnv() Environment {
	return Environment{
		Env:  "development",
		Port: "5000",
	}
}

func envMap() map[string]string {
	environ := os.Environ()
	env := make(map[string]string)

	for _, v := range environ {
		pair := strings.SplitN(v, "=", 2)
		env[pair[0]] = pair[1]
	}

	return env
}
