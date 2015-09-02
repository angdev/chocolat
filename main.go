package main

import (
	"github.com/angdev/chocolat/app"
	"os"
)

func main() {
	cli := app.NewCli()
	cli.Run(os.Args)
}
