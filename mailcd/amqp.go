package mailcd

import (
	"net/rpc"

	"github.com/streadway/amqp"
	"github.com/vbogretsov/maild/mail"
)

type amqpClient struct {
	client *rpc.Client
}

func (self *Client) Send(req mail.Request) error {
	reply := struct{}{}
	return self.client.Call("Maild.Send", req, &reply)
}

func New(client *rpc.Client) (mail.Sender, error) {
	return &amqpClient{client: client}, nil
}
