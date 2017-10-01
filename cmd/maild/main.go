package main

import (
	"os"

	"github.com/urfave/cli"
)

const (
	Name    = "maild"
	Usage   = "notification service for micro service architecture"
	Version = "0.0.0"
)

func newApp() *cli.App {
	app := cli.NewApp()
	app.Name = name
	app.Usage = usage
	app.Version = version
	return app
}

func main() {
	app := newApp()
	cfg := newConf(app)

	app.Action = func(ctx *cli.Context) {

	}

	app.Run(os.Args)
}
