package main

import (
	"fmt"
	"net/rpc"
	"os"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/op/go-logging"
	"github.com/streadway/amqp"
	"github.com/urfave/cli"
	"github.com/vbogretsov/amqprpc"
	"gopkg.in/go-playground/validator.v9"

	"github.com/vbogretsov/maild/sendgrid"
	"github.com/vbogretsov/maild/server"
)

var (
	log = logging.MustGetLogger(name)
)

func newApp() *cli.App {
	app := cli.NewApp()
	app.Name = name
	app.Usage = usage
	app.Version = version
	return app
}

func run(cfg *conf) error {
	amqpconn, err := amqp.Dial(cfg.BrokerURL)
	if err != nil {
		return err
	}
	defer amqpconn.Close()

	dbconn, err := gorm.Open("postgres", cfg.DatabaseDSN)
	if err != nil {
		return err
	}
	defer dbconn.Close()

	routingKey := fmt.Sprintf("%s-bg", cfg.ServiceName)

	producer, err := server.NewAMQPProducer(amqpconn, routingKey)
	if err != nil {
		return err
	}
	defer producer.Close()

	sender := server.NewMaild(producer, server.NewDbLoader(dbconn), 100)
	rpc.Register(sender)
	provider := sendgrid.NewProvider(cfg.ProviderURL, cfg.ProviderKey)

	consumer, err := server.NewAMQPConsumer(amqpconn, provider, routingKey)
	if err != nil {
		return err
	}
	defer consumer.Close()

	serverCodec, err := amqprpc.NewServerCodec(amqpconn, cfg.ServiceName)
	if err != nil {
		log.Fatal(err)
	}

	rpc.ServeCodec(serverCodec)
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
