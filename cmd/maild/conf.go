package main

import (
	"github.com/urfave/cli"
	"gopkg.in/mcuadros/go-defaults.v1"
)

type conf struct {
	BrokerURL   string `validate:"required" default:"amqp://localhost:5672"`
	ServiceName string `validate:"required" default:"maild"`
	ProviderURL string `validate:"required"`
	ProviderKey string `validate:"required"`
	DatabaseDSN string `validate:"required"`
	LogLevel    string `validate:"required" default:"INFO"`
}

func newConf() *conf {
	cfg := conf{}
	defaults.SetDefaults(&cfg)
	return &cfg
}

func bind(app *cli.App, cfg *conf) {
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "broker-url",
			Usage:       brokerURLUsage,
			Value:       cfg.BrokerURL,
			Destination: &cfg.BrokerURL,
		},
		cli.StringFlag{
			Name:        "service-name",
			Usage:       serviceNameUsage,
			Value:       cfg.ServiceName,
			Destination: &cfg.ServiceName,
		},
		cli.StringFlag{
			Name:        "provider-url",
			Usage:       providerURLUsage,
			Value:       cfg.ProviderURL,
			Destination: &cfg.ProviderURL,
		},
		cli.StringFlag{
			// TODO(vbogretsov): may be it should be file name for security reason
			Name:        "provider-key",
			Usage:       providerKeyUsage,
			Value:       cfg.ProviderKey,
			Destination: &cfg.ProviderKey,
		},
		cli.StringFlag{
			Name:        "dbdsn",
			Usage:       dbDSNUsage,
			Value:       cfg.DatabaseDSN,
			Destination: &cfg.DatabaseDSN,
		},
		cli.StringFlag{
			Name:        "log-level",
			Usage:       logLevelUsage,
			Value:       cfg.LogLevel,
			Destination: &cfg.LogLevel,
		},
	}
}
