package main

import (
	"fmt"
	"net/rpc"
	"os"
	"text/template"

	"github.com/op/go-logging"
	"github.com/streadway/amqp"
	"github.com/urfave/cli"
	"github.com/vbogretsov/amqprpc"
	"gopkg.in/go-playground/validator.v9"
	"gopkg.in/mcuadros/go-defaults.v1"

	"github.com/vbogretsov/maild/model"
	"github.com/vbogretsov/maild/pubsub"
	"github.com/vbogretsov/maild/sendgrid"
	"github.com/vbogretsov/maild/server"
)

const (
	name    = "maild"
	usage   = "notification service for micro service architecture"
	version = "0.0.0"
	logfmt  = `%{color}#%{id:03x} [%{pid}] %{time:15:04:05.000} %{shortfunc} ▶ %{level:.4s} %{message}%{color:reset}`

	brokerURLUsage    = `URL of the broker which holds the queue of requests`
	serviceNameUsage  = `name of the service is used for routing requests`
	providerURLUsage  = `URL of SMTP service provider`
	providerKeyUsage  = `SMTP service provider security key`
	providerNameUsage = `SMTP provider name, allowed valus: [sendgrid, log]`
	templateDirUsage  = `email templates location`
	logLevelUsage     = `log level, allowed values: [DEBUG, INFO, NOTICE, WARNING, ERROR, CRITICAL]`
)

var (
	log = logging.MustGetLogger(name)
)

type conf struct {
	BrokerURL   string `validate:"required" default:"amqp://localhost:5672"`
	ServiceName string `validate:"required" default:"maild"`
	Provider    struct {
		URL  string `validate:"required"`
		Key  string `validate:"required"`
		Name string `validate:"required" default:"log"`
	} `validate:"-"`
	// TODO(vbogretsov): validate dir exists
	TemplateDir string `validate:"required"`
	LogLevel    string `validate:"required" default:"INFO"`
}

func (c *conf) Validate() error {
	v := validator.New()
	if err := v.Struct(c); err != nil {
		return err
	}

	if c.Provider.Name != "log" {
		if err := v.Struct(c.Provider); err != nil {
			return err
		}
	}

	return nil
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
			Value:       cfg.Provider.URL,
			Destination: &cfg.Provider.URL,
		},
		cli.StringFlag{
			// TODO(vbogretsov): may be it should be file name for security reason
			Name:        "provider-key",
			Usage:       providerKeyUsage,
			Value:       cfg.Provider.Key,
			Destination: &cfg.Provider.Key,
		},
		cli.StringFlag{
			Name:        "provider-name",
			Usage:       providerNameUsage,
			Value:       cfg.Provider.Name,
			Destination: &cfg.Provider.Name,
		},
		cli.StringFlag{
			Name:        "template-dir",
			Usage:       templateDirUsage,
			Value:       cfg.TemplateDir,
			Destination: &cfg.TemplateDir,
		},
		cli.StringFlag{
			Name:        "log-level",
			Usage:       logLevelUsage,
			Value:       cfg.LogLevel,
			Destination: &cfg.LogLevel,
		},
	}
}

type consumer struct {
	sendmail server.SendMail
}

func (c *consumer) Create() interface{} {
	return &model.Message{}
}

func (c *consumer) Consume(v interface{}) error {
	m := v.(*model.Message)
	return c.sendmail(*m)
}

func newConsumer(cfg *conf) *consumer {
	if cfg.Provider.Name == "sendgrid" {
		provider := sendgrid.NewProvider(cfg.Provider.URL, cfg.Provider.Key)
		return &consumer{provider.SendMail}
	}

	return &consumer{func(m model.Message) error {
		log.Infof("sent\n%v", m)
		return nil
	}}
}

func run(cfg *conf) error {
	log.Info("connecting broker ....")
	amqpconn, err := amqp.Dial(cfg.BrokerURL)
	if err != nil {
		return err
	}
	defer amqpconn.Close()
	log.Info("connecting broker done")

	topic := fmt.Sprintf("%s-bg", cfg.ServiceName)

	producer, err := pubsub.NewProducer(amqpconn, topic)
	if err != nil {
		return err
	}
	defer producer.Close()

	sendmail := func(m model.Message) error {
		return producer.Produce(&m)
	}

	globpattern := fmt.Sprintf("%s/*.msg", cfg.TemplateDir)

	log.Info("loading templates ....")
	templates, err := template.ParseGlob(globpattern)
	if err != nil {
		return err
	}
	templateNames := []string{}
	for _, t := range templates.Templates() {
		templateNames = append(templateNames, t.Name())
	}
	log.Infof("loaded tempaltes:\n%v", templateNames)

	maild := server.NewMaild(sendmail, templates)
	rpc.Register(maild)

	queue, err := pubsub.NewQueue(amqpconn, newConsumer(cfg), topic)
	if err != nil {
		return err
	}
	defer queue.Close()

	serverCodec, err := amqprpc.NewServerCodec(amqpconn, cfg.ServiceName)
	if err != nil {
		log.Fatal(err)
	}

	log.Info("maild stared")

	rpc.ServeCodec(serverCodec)
	return nil
}

func main() {
	logging.SetBackend(logging.NewLogBackend(os.Stderr, "", 0))

	app := newApp()
	cfg := newConf()
	bind(app, cfg)

	app.Action = func(c *cli.Context) error {
		if err := cfg.Validate(); err != nil {
			log.Fatal(err)
			return err
		}

		logging.SetFormatter(logging.MustStringFormatter(logfmt))

		if err := run(cfg); err != nil {
			log.Error(err)
			return err
		}

		return nil
	}

	app.Run(os.Args)
}
