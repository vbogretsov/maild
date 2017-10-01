package mail

import (
	"bytes"
	"errors"
	"text/template"

	"github.com/bluele/gcache"
	"gopkg.in/go-playground/validator.v9"
	"gopkg.in/yaml.v2"
)

var valid = validator.New()

type TemplateID struct {
	Lang string `validate:"required"`
	ID   string `validate:"required"`
}

// Request represents the parameters required to build an email from template
// and send it.
type Request struct {
	TemplateID TemplateID `validate:"required"`
	Args       interface{}
}

// Message represents an email message.
type Message struct {
	From     string   `yaml:"From,omitempty"`
	To       []string `yaml:"To,omitempty"`
	Cc       []string `yaml:"Cc,omitempty"`
	Subject  string   `yaml:"Subject,omitempty"`
	Body     string   `yaml:"Body,omitempty"`
	BodyType string   `yaml:"BodyType,omitempty"`
}

// A Sender represents interface for sending email according to request.
type Sender interface {
	Send(Request) error
}

// A Provider represents interface for sending email.
type Provider interface {
	SendMail(Message) error
}

type sender struct {
	templatesCache gcache.Cache
	provider       Provider
}

// Send builds email from template and send it via provider.
func (snd *sender) Send(req Request) error {
	err := valid.Struct(req)
	if err != nil {
		return errors.New("invalid request: " + err.Error())
	}

	item, err := snd.templatesCache.Get(req.TemplateID)
	if err != nil {
		return err
	}

	template, ok := item.(*template.Template)
	if !ok {
		return errors.New("unable to cast template")
	}

	buf := bytes.Buffer{}

	err = template.Execute(&buf, req.Args)
	if err != nil {
		return err
	}

	msg := Message{}

	err = yaml.Unmarshal(buf.Bytes(), &msg)
	if err != nil {
		return err
	}

	return snd.provider.SendMail(msg)
}

// Create new sender.
func NewSender(provider Provider, templatesCache gcache.Cache) Sender {
	return &sender{
		provider:       provider,
		templatesCache: templatesCache,
	}
}
