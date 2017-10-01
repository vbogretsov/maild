package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli"
	"gopkg.in/go-playground/validator.v9"
	"gopkg.in/mcuadros/go-defaults.v1"
)

var (
	valid = validator.New()
)

type commonConf struct {
	AMQPUrl        string `default:"amqp://localhost:5672"`
	AMQPRoutingKey string `validate:"required" default:"maild"`
}

type sendConf struct {
	Lang       string `validate:"required"`
	TemplateID string `validate:"required"`
	Args       string `default:"{}"`
}

type uploadConf struct {
	Lang       string `validate:"required"`
	TemplateID string `validate:"required"`
	Template   string `validate:"required"`
}

type conf struct {
	Common commonConf
	Send   sendConf
	Upload uploadConf
}

func (c *commonConf) Validate() {
	if err := valid.Struct(c); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(2)
	}
}

func (c *sendConf) Validate() {
	if err := valid.Struct(c); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(2)
	}
}

func (c *uploadConf) Validate() {
	if err := valid.Struct(c); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(2)
	}
}

func newConf(app *cli.App) *conf {
	var cfg conf
	defaults.SetDefaults(&cfg)

	if app == nil {
		return &cfg
	}

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "amqp-url",
			Value:       cfg.Common.AMQPUrl,
			Usage:       "AMQP broker URL",
			Destination: &cfg.Common.AMQPUrl,
		},
		cli.StringFlag{
			Name:        "amqp-routing",
			Value:       cfg.Common.AMQPRoutingKey,
			Usage:       "AMQP routing key",
			Destination: &cfg.Common.AMQPRoutingKey,
		},
	}

	app.Commands = []cli.Command{
		{
			Name:  "send",
			Usage: "send email",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:        "lang",
					Value:       cfg.Send.Lang,
					Usage:       "template language",
					Destination: &cfg.Send.Lang,
				},
				cli.StringFlag{
					Name:        "template-id",
					Value:       cfg.Send.TemplateID,
					Usage:       "template id",
					Destination: &cfg.Send.TemplateID,
				},
				cli.StringFlag{
					Name:        "args",
					Value:       cfg.Send.Args,
					Usage:       "email template arguments",
					Destination: &cfg.Send.Args,
				},
			},
			Action: func(ctx *cli.Context) {
				cfg.Common.Validate()

				cfg.Send.Validate()

				fmt.Println("send", cfg)
			},
		},
		{
			Name:  "upload",
			Usage: "upload email template",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:        "lang",
					Value:       cfg.Upload.Lang,
					Usage:       "template language",
					Destination: &cfg.Upload.Lang,
				},
				cli.StringFlag{
					Name:        "template-id",
					Value:       cfg.Upload.TemplateID,
					Usage:       "template id",
					Destination: &cfg.Upload.TemplateID,
				},
				cli.StringFlag{
					Name:        "template",
					Value:       cfg.Upload.Template,
					Usage:       "template file path",
					Destination: &cfg.Upload.Template,
				},
			},
			Action: func(ctx *cli.Context) {
				cfg.Common.Validate()
				cfg.Upload.Validate()

				fmt.Println("upload", cfg)
			},
		},
	}

	return &cfg
}
