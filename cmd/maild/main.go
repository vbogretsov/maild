package main

import (
	"fmt"
	"net/rpc"
	"os"

	"github.com/streadway/amqp"
	"github.com/urfave/cli"
	"github.com/vbogretsov/amqprpc"
	"gopkg.in/go-playground/validator.v9"
	"gopkg.in/mcuadros/go-defaults.v1"

	"github.com/vbogretsov/maild"
	"github.com/vbogretsov/maild/provider"
	"github.com/vbogretsov/maild/tmplstore"
)

const (
	name    = "maild"
	usage   = "notification service for micro service architecture"
	version = "0.0.0"
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

// Maild is the maild.Sender RPC implementation.
type Maild struct {
	sender maild.Sender
}

// Send implements maild.Sender interface.
func (server *Maild) Send(md maild.MailData, out *int) error {
	return server.sender.Send(md)
}

func run(cfg *conf) error {
	conn, err := amqp.Dial(cfg.AMQPUrl)
	if err != nil {
		return err
	}
	defer conn.Close()

	templateStore, err := tmplstore.NewDiskStore(cfg.TemplatesDir)
	if err != nil {
		return err
	}

	provider := provider.NewConsolePorovider()

	sender := maild.NewSender(templateStore, provider)
	server := Maild{sender}

	serverCodec, err := amqprpc.NewServerCodec(conn, cfg.Routing)
	if err != nil {
		return err
	}
	defer serverCodec.Close()

	rpc.Register(&server)
	rpc.ServeCodec(serverCodec)

	return nil
}

func createApp() *cli.App {
	app := cli.NewApp()
	app.Name = name
	app.Usage = usage
	app.Version = version
	return app
}

func main() {
	app := createApp()
	cfg := newConf(app)

	app.Action = func(c *cli.Context) {
		if err := cfg.Validate(); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(2)
		}
		if err := run(cfg); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
	}

	app.Run(os.Args)
}
