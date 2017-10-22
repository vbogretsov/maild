package main

import (
	"fmt"

	"github.com/op/go-logging"
	"github.com/streadway/amqp"
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

const email = `
From:
  Email: levelup@levelup.com
  Name: LevelUP
To:
  - Address:
    Email: to1@mail.com
  - Address:
    Email: to2@mail.com
Subject: Subject
BodyType: "text/plain"
Body: |
  Hello!
  This is test body!
`

func newApp() *cli.App {
	app := cli.NewApp()
	app.Name = name
	app.Usage = usage
	app.Version = version
	return app
}

func main() {
	logging.SetBackend(logging.NewLogBackend(os.Stderr, "", 0))

	app := newApp()
	cfg := newConf()

	app.Action = func(c *cli.Context) {
		if err := cfg.Validate(); err != nil {
			log.Fatalf("error: %v", err)
		}

		logging.SetFormatter(logging.MustStringFormatter(logfmt))

		for {
			if err := run(cfg); err != nil {
				log.Errorf(err)
			}
			time.Sleep(time.Second * 1)
		}
	}

	app.Run(os.Args)
}
