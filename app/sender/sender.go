package sender

import (
	"fmt"

	"github.com/vbogretsov/maild/app"
	"github.com/vbogretsov/maild/app/sender/sendgrid"
)

type factory func(url, key string) app.Sender

var senders = map[string]factory{
	"sendgrid": sendgrid.New,
}

// New creates a new sender.
func New(name, url, key string) (app.Sender, error) {
	fn, ok := senders[name]
	if !ok {
		return nil, fmt.Errorf("unsupported sender %s", name)
	}
	return fn(url, key), nil
}
