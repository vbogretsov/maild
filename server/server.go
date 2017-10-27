package server

import (
	"bytes"
	"text/template"

	"gopkg.in/go-playground/validator.v9"
	"gopkg.in/yaml.v1"

	"github.com/vbogretsov/maild/model"
)

var (
	valid = validator.New()
)

type Provider interface {
	SendMail(*model.Message) error
}

type TemplateCache interface {
	Get(lang, name string) (*template.Template, error)
}

type Maild struct {
	provider      Provider
	templateCache TemplateCache
}

func (m *Maild) Send(r *model.Request) error {
	if err := valid.Struct(r); err != nil {
		return err
	}

	tmpl, err := m.templateCache.Get(r.TemplateLang, r.TemplateName)
	if err != nil {
		return err
	}

	buff := bytes.Buffer{}

	if err := tmpl.Execute(&buff, r.TemplateArgs); err != nil {
		return err
	}

	msg := model.Message{}

	if err := yaml.Unmarshal(buff.Bytes(), &msg); err != nil {
		return err
	}

	msg.To = append(msg.To, r.To...)
	msg.Cc = append(msg.Cc, r.Cc...)

	return m.provider.SendMail(&msg)
}

func New(provider Provider, templateCache TemplateCache) *Maild {
	return &Maild{
		provider:      provider,
		templateCache: templateCache,
	}
}
