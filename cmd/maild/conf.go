package main

import (
	"github.com/urfave/cli"
	"gopkg.in/go-playground/validator.v9"
	"gopkg.in/mcuadros/go-defaults.v1"
)

type Conf struct {
	Provider         string `validate:"required" default:"sendgrid"`
	ProviderEndpoint string `validate:"required"`
	Protocol         string `validate:"required" default:"amqp"`
	HTTPPort         string `default:"50101"`
	AMQPUrl          string `default:"localhost:5672"`
	LogLevel         string `validate:"required" default:"INFO"`
	TemplatesDir     string `validate:"required" default:"./templates"`
}

func New(app *cli.App) *Conf {
	var cfg Conf
	defaults.SetDefaults(&cfg)

	if app == nil {
		return &cfg
	}

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "provider",
			Value:       cfg.Provider,
			Usage:       "SMTP provider name, allowed values: sendgrid",
			Destination: &cfg.Provider,
		},
		cli.StringFlag{
			Name:        "provider-endpoint",
			Value:       cfg.ProviderEndpoint,
			Usage:       "Provider connection information",
			Destination: &cfg.ProviderEndpoint,
		},
		cli.StringFlag{
			Name:        "http-port",
			Value:       cfg.HTTPPort,
			Usage:       "HTTP service port number, ignored if protocol is AMQP",
			Destination: &cfg.HTTPPort,
		},
		cli.StringFlag{
			Name:        "amqp-url",
			Value:       cfg.AMQPUrl,
			Usage:       "AMQP broker URL, ignored if protocol is HTTP",
			Destination: &cfg.AMQPUrl,
		},
		cli.StringFlag{
			Name:        "protocol",
			Value:       cfg.Protocol,
			Usage:       "Service protocol name, allowed values: [HTTP, AMQP]",
			Destination: &cfg.Protocol,
		},
		cli.StringFlag{
			Name:        "log-levle",
			Value:       cfg.LogLevel,
			Usage:       "Log level, allowed values: [INFO, WARNING, ERROR, DEBUG]",
			Destination: &cfg.LogLevel,
		},
		cli.StringFlag{
			Name:        "templates-dir",
			Value:       cfg.TemplatesDir,
			Usage:       "Templates directory path",
			Destination: &cfg.TemplatesDir,
		},
	}

	return &cfg
}

func (c *Conf) Validate() error {
	return validator.New().Struct(c)
}
