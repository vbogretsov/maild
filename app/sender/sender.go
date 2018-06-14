package sender

import (
	"github.com/vbogretsov/maild/model"
)

type Sender interface {
	Send(model.Message) error
}

func New(name, url, key string) (Sender, error) {
	return nil, nil
}
