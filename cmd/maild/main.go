package main

import (
	"net/rpc"
	"os"
	"time"

	"github.com/op/go-logging"
	"github.com/streadway/amqp"
	"github.com/urfave/cli"
	"github.com/vbogretsov/amqprpc"

	"github.com/vbogretsov/maild"
	"github.com/vbogretsov/maild/provider"
	"github.com/vbogretsov/maild/tmplstore"
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

// Maild is the maild.Sender RPC implementation.
type Maild struct {
	sender maild.Sender
}

// Send implements maild.Sender interface.
func (server *Maild) Send(md maild.MailData, out *int) error {
	return server.sender.Send(md)
}

func run(cfg *conf) error {
	log.Info("connecting AMQP")

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
	logging.SetBackend(logging.NewLogBackend(os.Stderr, "", 0))

	app := createApp()
	cfg := newConf(app)

	app.Action = func(c *cli.Context) {
		if err := cfg.Validate(); err != nil {
			log.Fatalf("error: %v", err)
		}

		logging.SetFormatter(logging.MustStringFormatter(logfmt))

		for {
			if err := run(cfg); err != nil {
				log.Errorf("AMQP connection failed %v", err)
			}
			time.Sleep(time.Second * 1)
		}
	}

	app.Run(os.Args)
}
