package main

import (
	"github.com/urfave/cli"
	"gopkg.in/go-playground/validator.v9"
	"gopkg.in/mcuadros/go-defaults.v1"
)

type conf struct {
	Provider         string `validate:"required" default:"sendgrid"`
	ProviderEndpoint string `validate:"required"`
	AMQPUrl          string `default:"amqp://localhost:5672"`
	LogLevel         string `validate:"required" default:"INFO"`
	TemplatesDir     string `validate:"required" default:"./templates"`
	Routing          string `validate:"required" default: "maild"`
}

func newConf(app *cli.App) *conf {
	var cfg conf
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
			Name:        "amqp-url",
			Value:       cfg.AMQPUrl,
			Usage:       "AMQP broker URL, ignored if protocol is HTTP",
			Destination: &cfg.AMQPUrl,
		},
		cli.StringFlag{
			Name:        "routing",
			Value:       cfg.Routing,
			Usage:       "AMQP routing key",
			Destination: &cfg.Routing,
		},
		cli.StringFlag{
			Name:        "log-level",
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

func (c *conf) Validate() error {
	return validator.New().Struct(c)
}
