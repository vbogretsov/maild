package main

import (
	"errors"
	"net/rpc"

	"github.com/bluele/gcache"
	"github.com/streadway/amqp"
	"github.com/vbogretsov/amqprpc"
	"github.com/vmihailenco/msgpack"

	"github.com/vbogretsov/maild/mail"
)

type Template struct {
	ID      mail.TemplateID
	Concent []byte
}

type Maild struct {
	sender    mail.Sender
	channel   *amqp.Channel
	queueName string
}

func NewServer(conn *amqp.Connection, sender mail.Sender) (*Maild, error) {
	channel, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	queue, err := channel.QueueDeclare(
		"", true, false, false, false, nil)

	if err != nil {
		return nil, err
	}

	requests, err := channel.Consume(
		queue.Name, "", false, false, false, false, nil)

	if err != nil {
		return nil, err
	}

	go func() {
		for req := range requests {
			request := mail.Request{}

			err := msgpack.Unmarshal(req.Body, &requests)
			if err != nil {
				log.Fatal("msgpack.Unmarshal call failed")
			}

			err = sender.Send(request)

			if err != nil {
				log.Debug("send request failed")
				req.Nack(false, true)
			} else {
				log.Debug("send request succeed")
				req.Ack(false)
			}
		}
	}()

	server := Maild{
		sender:    sender,
		channel:   channel,
		queueName: queue.Name,
	}

	return &server, nil
}

func (self *Maild) Send(req *mail.Request, out *struct{}) error {
	log.Debugf("accepted request %v", req)

	body, err := msgpack.Marshal(req)
	if err != nil {
		return err
	}

	publishing := amqp.Publishing{
		CorrelationId: self.queueName,
		Body:          body,
	}

	return self.channel.Publish("", "", true, true, publishing)
}

func (self *Maild) Upload(req *Template, out *struct{}) error {
	return errors.New("not implemented")
}

type logProvider struct {
}

func (self *logProvider) SendMail(msg mail.Message) error {
	log.Infof("sent %v", msg)
	return nil
}

func run(cfg *conf) error {
	log.Debugf("connecting AMQP borker %s", cfg.AMQPUrl)

	conn, err := amqp.Dial(cfg.AMQPUrl)
	if err != nil {
		return err
	}
	defer conn.Close()

	serverCodec, err := amqprpc.NewServerCodec(conn, cfg.AMQPRoutingKey)
	if err != nil {
		return err
	}
	defer serverCodec.Close()

	sender := mail.NewSender(&logProvider{}, gcache.New(100).LRU().Build())

	maild, err := NewServer(conn, sender)
	if err != nil {
		return err
	}

	rpc.Register(maild)
	rpc.ServeCodec(serverCodec)

	return nil
}
