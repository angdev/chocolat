package main

import (
	"github.com/angdev/chocolat/app"
	"os"
)

func init() {
	app.Chocolat.Init()
}

func main() {
	cli := app.NewCli()
	cli.Run(os.Args)
}
