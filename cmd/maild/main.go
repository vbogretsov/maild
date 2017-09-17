package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli"

	"github.com/vbogretsov/maild/maild/conf"
)

const (
	Name    = "maild"
	Usage   = "notification service for micro service architecture"
	Version = "0.0.0"
)

func run(cfg *conf.Conf) {
	fmt.Println("maild started with configuration:", cfg)
}

func createApp() *cli.App {
	app := cli.NewApp()
	app.Name = Name
	app.Usage = Usage
	app.Version = Version
	return app
}

func main() {
	app := createApp()
	cfg := conf.New(app)

	app.Action = func(c *cli.Context) {
		run(cfg)
	}

	app.Run(os.Args)
}
