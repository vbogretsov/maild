// Package maild provides sending emails.
package maild

import (
	"bytes"
	"errors"
	"sync"
	"text/template"

	"gopkg.in/yaml.v2"
)

// Message represents an email message.
type Message struct {
	Subject  string   `yaml:"Subject,omitempty"`
	Body     string   `yaml:"Body,omitempty"`
	BodyType string   `yaml:"BodyType,omitempty"`
	From     string   `yaml:"From,omitempty"`
	To       []string `yaml:"To,omitempty"`
	Cc       []string `yaml:"Cc,omitempty"`
}

// MailData contains parameters required to build email from template and send it.
type MailData struct {
	Lang       string
	TemplateID string
	Args       interface{}
}

// A Sender implements sending email that will be generated from template
// with the id provided using the args provided.
type Sender interface {
	Send(md MailData) error
}

// A TemplateLoader implements loading tempalte.
type TemplateLoader interface {
	Load(lang string, id string) (*template.Template, error)
}

// A Provider implements sending email message.
type Provider interface {
	Send(msg Message) error
}

type maildServer struct {
	templateLoader TemplateLoader
	provider       Provider
	mutex          sync.RWMutex
	templatesCache map[string]*template.Template
}

func (server *maildServer) Send(md MailData) error {
	templateKey := md.Lang + md.TemplateID
	var template *template.Template

	server.mutex.RLock()
	template, ok := server.templatesCache[templateKey]
	server.mutex.RUnlock()

	if !ok {
		tmpl, err := server.templateLoader.Load(md.Lang, md.TemplateID)
		if err != nil {
			return err
		}

		server.mutex.Lock()
		server.templatesCache[templateKey] = tmpl
		server.mutex.Unlock()

		template = tmpl
	}

	buf := bytes.Buffer{}
	err := template.Execute(&buf, md.Args)
	if err != nil {
		return err
	}

	msg := Message{}
	err = yaml.Unmarshal(buf.Bytes(), &msg)
	if err != nil {
		return err
	}

	return server.provider.Send(msg)
}

// NewSender creates a new implementation of Sender interface.
func NewSender(templateLoader TemplateLoader, provider Provider) Sender {
	return &maildServer{
		templateLoader: templateLoader,
		provider:       provider,
		mutex:          sync.RWMutex{},
		templatesCache: map[string]*template.Template{},
	}
}
