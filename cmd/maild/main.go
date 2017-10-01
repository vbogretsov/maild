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

	app.Action = func(c *cli.Context) {
		if err := cfg.Validate(); err != nil {
			log.Fatalf("error: %v", err)
		}

		logging.SetFormatter(logging.MustStringFormatter(logfmt))

		for {
			if err := run(cfg); err != nil {
				log.Errorf("AMQP connection failed %v", err)
			}
			time.Sleep(time.Second * 1)
		}
	}

	app.Run(os.Args)
}
