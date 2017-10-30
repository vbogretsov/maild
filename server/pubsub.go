package server

import (
	"github.com/op/go-logging"
	"github.com/streadway/amqp"
	"github.com/vmihailenco/msgpack"

	"github.com/vbogretsov/maild/model"
)

var (
	log = logging.MustGetLogger("maild")
)

type AMQPProducer struct {
	channel  *amqp.Channel
	exchange string
	key      string
}

func (p *AMQPProducer) SendMail(message *model.Message) error {
	body, err := msgpack.Marshal(message)

	if err != nil {
		return err
	}

	pub := amqp.Publishing{
		Body: body,
	}

	return p.channel.Publish(p.exchange, p.key, false, false, pub)
}

func (p *AMQPProducer) Close() error {
	return p.channel.Close()
}

func NewAMQPProducer(conn *amqp.Connection, serviceName string) (*AMQPProducer, error) {
	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	producer := AMQPProducer{
		channel:  ch,
		exchange: serviceName,
		key:      serviceName,
	}
	return &producer, nil
}

type AMQPConsumer struct {
	channel *amqp.Channel
}

func (c *AMQPConsumer) Close() error {
	return c.channel.Close()
}

func NewAMQPConsumer(conn *amqp.Connection, provier model.Provider, serviceName string) (*AMQPConsumer, error) {
	var err error

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			ch.Close()
		}
	}()

	err = ch.ExchangeDeclare(
		serviceName, // name
		"direct",    // type
		true,        // durable
		false,       // auto-deleted
		false,       // internal
		false,       // no-wait
		nil,         // arguments
	)
	if err != nil {
		return nil, err
	}

	q, err := ch.QueueDeclare(
		serviceName, // name
		true,        // durable
		false,       // delete when usused
		false,       // exclusive
		false,       // no-wait
		nil,         // arguments
	)
	if err != nil {
		return nil, err
	}

	err = ch.QueueBind(
		q.Name,      // queue name
		serviceName, // routing key
		serviceName, // exchange
		false,       // no-wait
		nil,         // arguments
	)
	if err != nil {
		return nil, err
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
		return nil, err
	}

	go func() {
		for msg := range msgs {
			message := model.Message{}

			if err := msgpack.Unmarshal(msg.Body, &message); err != nil {
				log.Errorf("msgpack.Unmarhshal failed: %v", err)
				continue
			}

			if err := provier.SendMail(&message); err != nil {
				log.Errorf("provier.SendMail failed: %v", err)
				msg.Nack(false, true)
			} else {
				msg.Ack(false)
			}
		}
	}()

	return &AMQPConsumer{channel: ch}, nil
}
