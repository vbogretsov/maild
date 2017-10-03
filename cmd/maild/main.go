package main

import (
	"os"
	"time"

	"github.com/op/go-logging"
	"github.com/urfave/cli"
)

const (
	name    = "maild"
	usage   = "notification service for micro service architecture"
	version = "0.0.0"
	logfmt  = `%{color}#%{id:03x} [%{pid}]%{time:15:04:05.000} %{shortfunc} ▶ %{level:.4s} %{message}%{color:reset}`
)

var (
	log = logging.MustGetLogger(name)
)

func init() {
	logging.SetBackend(logging.NewLogBackend(os.Stderr, "", 0))
}

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
		// TODO(vbogretsov): set log level

		for {
			if err := run(cfg); err != nil {
				log.Errorf("AMQP connection failed %v", err)
			}
			time.Sleep(time.Second * 1)
		}
	}

	app.Run(os.Args)
}
