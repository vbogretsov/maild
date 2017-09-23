package client

import (
	"net/rpc"

	"github.com/hashicorp/go-multierror"
	"github.com/streadway/amqp"
	"github.com/vbogretsov/amqprpc"

	"github.com/vbogretsov/maild"
)

type AMQPClient struct {
	codec  rpc.ClientCodec
	client *rpc.Client
}

func (c *AMQPClient) Send(md maild.MailData) error {
	res := 0
	return c.client.Call("Maild.Send", md, &res)
}

func (c *AMQPClient) Close() error {
	res := &multierror.Error{}

	if err := c.codec.Close(); err != nil {
		res = multierror.Append(res, err)
	}

	if err := c.client.Close(); err != nil {
		res = multierror.Append(res, err)
	}

	return res.ErrorOrNil()
}

func NewAMQPClient(conn *amqp.Connection, routing string) (*AMQPClient, error) {
	codec, err := amqprpc.NewClientCodec(conn, routing)
	if err != nil {
		return nil, err
	}

	client := rpc.NewClientWithCodec(codec)
	return &AMQPClient{codec, client}, nil
}
