package main

import (
	"os"

	"github.com/urfave/cli"
)

const (
	name    = "mailctl"
	usage   = "notification service control"
	version = "0.0.0"
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
	newConf(app)
	app.Run(os.Args)
}
