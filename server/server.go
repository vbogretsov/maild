package server

import (
	"bytes"
	"fmt"
	"text/template"

	"github.com/pkg/errors"
	"gopkg.in/go-playground/validator.v9"
	"gopkg.in/yaml.v1"

	"github.com/vbogretsov/maild/model"
)

var (
	valid = validator.New()
)

type SendMail func(model.Message) error

// Maild represents a maild server.
type Maild struct {
	template *template.Template
	sendmail SendMail
}

// Send builds an email from template and sends it via sendmail function.
func (m *Maild) Send(r model.Request, out *struct{}) error {
	if err := valid.Struct(r); err != nil {
		return errors.Wrap(err, "invalid request")
	}

	key := fmt.Sprintf("%s-%s.msg", r.TemplateLang, r.TemplateName)
	tml := m.template.Lookup(key)

	if tml == nil {
		return fmt.Errorf("template not found: %v", key)
	}

	buf := new(bytes.Buffer)
	if err := tml.Execute(buf, r.TemplateArgs); err != nil {
		return errors.Wrap(err, "template execution failed")
	}

	msg := model.Message{}

	if err := yaml.Unmarshal(buf.Bytes(), &msg); err != nil {
		return errors.Wrap(err, "unable to build mail message")
	}

	msg.To = append(msg.To, r.To...)
	msg.Cc = append(msg.Cc, r.Cc...)

	return m.sendmail(msg)
}

// NewMaild creates a new Maild server.
func NewMaild(sendmail SendMail, template *template.Template) *Maild {
	return &Maild{template: template, sendmail: sendmail}
}
