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
	AMQPRoutingKey   string `validate:"required" default:"maild"`
	LogLevel         string `validate:"required" default:"INFO"`
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
			Name:        "amqp-routing",
			Value:       cfg.AMQPRoutingKey,
			Usage:       "AMQP routing key",
			Destination: &cfg.AMQPRoutingKey,
		},
		cli.StringFlag{
			Name:        "log-level",
			Value:       cfg.LogLevel,
			Usage:       "Log level, allowed values: [INFO, WARNING, ERROR, DEBUG]",
			Destination: &cfg.LogLevel,
		},
	}

	return &cfg
}

func (c *conf) Validate() error {
	return validator.New().Struct(c)
}
