package pubsub

// TODO(vbogretsov): move to separate repository where RabbitMQ main patters
// will be implemented and covered with integration tests.

import (
	"github.com/op/go-logging"

	"github.com/streadway/amqp"
	"github.com/vmihailenco/msgpack"
)

var (
	log = logging.MustGetLogger("pubsub")
)

// Consumer represents interface of the delayed queue consumer.
type Consumer interface {
	Create() interface{}
	Consume(interface{}) error
}

type Producer struct {
	channel *amqp.Channel
	topic   string
}

func (p *Producer) Produce(v interface{}) error {
	body, err := msgpack.Marshal(v)

	if err != nil {
		return err
	}

	pub := amqp.Publishing{
		Body: body,
	}

	return p.channel.Publish(p.topic, p.topic, false, false, pub)
}

func (p *Producer) Close() error {
	return p.channel.Close()
}

func NewProducer(conn *amqp.Connection, topic string) (*Producer, error) {
	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	producer := Producer{
		channel: ch,
		topic:   topic,
	}
	return &producer, nil
}

type Queue struct {
	channel *amqp.Channel
}

func (q *Queue) Close() error {
	return q.channel.Close()
}

func NewQueue(conn *amqp.Connection, consumer Consumer, topic string) (*Queue, error) {
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
		topic,    // name
		"direct", // type
		true,     // durable
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)
	if err != nil {
		return nil, err
	}

	q, err := ch.QueueDeclare(
		topic, // name
		true,  // durable
		false, // delete when usused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return nil, err
	}

	err = ch.QueueBind(
		q.Name, // queue name
		topic,  // routing key
		topic,  // exchange
		false,  // no-wait
		nil,    // arguments
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
			item := consumer.Create()

			if err := msgpack.Unmarshal(msg.Body, item); err != nil {
				log.Errorf("msgpack.Unmarhshal failed: %v", err)
				continue
			}

			if err := consumer.Consume(item); err != nil {
				log.Errorf("consumer.Consume failed: %v", err)
				msg.Nack(false, true)
			} else {
				msg.Ack(false)
			}
		}
	}()

	return &Queue{channel: ch}, nil
}
