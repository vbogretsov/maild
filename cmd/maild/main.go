package main

import (
	"os"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/op/go-logging"
	"github.com/streadway/amqp"
	"github.com/urfave/cli"
	"github.com/vmihailenco/msgpack"
	"gopkg.in/go-playground/validator.v9"
	"gopkg.in/mcuadros/go-defaults.v1"

	"github.com/vbogretsov/maild/model"
	"github.com/vbogretsov/maild/server"
)

var (
	log = logging.MustGetLogger(name)
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

func newApp() *cli.App {
	app := cli.NewApp()
	app.Name = name
	app.Usage = usage
	app.Version = version
	return app
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

func run(cfg *conf) error {
	amqpconn, err := amqp.Dial(cfg.BrokerURL)
	if err != nil {
		return err
	}
	defer amqpconn.Close()

	ch, err := amqpconn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	if err := ch.ExchangeDeclare(
		cfg.ServiceName, // name
		"direct",        // type
		true,            // durable
		false,           // auto-deleted
		false,           // internal
		false,           // no-wait
		nil,             // arguments
	); err != nil {
		return err
	}

	q, err := ch.QueueDeclare(
		"",    // name
		true,  // durable
		false, // delete when usused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return err
	}

	if err := ch.QueueBind(
		q.Name,          // queue name
		cfg.ServiceName, // routing key
		cfg.ServiceName, // exchange
		false,           // no-wait
		nil,             // arguments
	); err != nil {
		return err
	}

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // arguments
	)
	if err != nil {
		return err
	}

	// TODO(vbogretsov) get database type from dsn.
	dbconn, err := gorm.Open("postgres", cfg.DatabaseDSN)
	if err != nil {
		return err
	}
	defer dbconn.Close()

	provider := server.NewSendGridProvider(cfg.ProviderURL, cfg.ProviderKey)
	tmlcache := server.NewDbCache(100, dbconn)
	mailserv := server.New(provider, tmlcache)

	log.Info("ready")

	for msg := range msgs {
		request := model.Request{}
		if err := msgpack.Unmarshal(msg.Body, &request); err != nil {
			log.Errorf("unable to unpack request %v", err)
			continue
		}

		if err := mailserv.Send(&request); err != nil {
			log.Errorf("send failed %v", err)
			msg.Nack(false, true)
		} else {
			msg.Ack(false)
		}
	}

	return nil
}

func main() {
	logging.SetBackend(logging.NewLogBackend(os.Stderr, "", 0))

	app := newApp()
	cfg := newConf()
	bind(app, cfg)

	app.Action = func(c *cli.Context) error {
		if err := validator.New().Struct(cfg); err != nil {
			log.Fatal(err)
			return err
		}

		logging.SetFormatter(logging.MustStringFormatter(logfmt))

		for {
			if err := run(cfg); err != nil {
				log.Error(err)
			}
			time.Sleep(time.Second * 1)
		}

		return nil
	}

	app.Run(os.Args)
}
