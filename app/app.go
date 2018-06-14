package app

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"text/template"

	"github.com/vbogretsov/go-validation"
	"gopkg.in/go-playground/validator.v9"
	"gopkg.in/yaml.v1"

	"github.com/vbogretsov/maild/app/loader"
	"github.com/vbogretsov/maild/app/sender"
	"github.com/vbogretsov/maild/model"
)

var Validate = validator.New()

type App struct {
	loader loader.Loader
	sender sender.Sender
}

func New(loader loader.Loader, sender sender.Sender) *App {
	return &App{
		loader: loader,
		sender: sender,
	}
}

func (a *App) SendMail(req model.Request) error {
	if err := Validate.Struct(req); err != nil {
		return validation.Error(err)
	}

	if len(req.To) == 0 && len(req.Cc) == 0 && len(req.Bcc) == 0 {
		return validation.Error(errors.New("missing recipients"))
	}

	body, err := a.loader.Load(req.TemplateLang, req.TemplateName)
	if err != nil {
		return err
	}

	key := fmt.Sprintf("%s-%s", req.TemplateLang, req.TemplateName)

	text, err := ioutil.ReadAll(body)
	if err != nil {
		return err
	}

	tml, err := template.New(key).Parse(string(text))
	if err != nil {
		return err
	}

	buf := new(bytes.Buffer)
	if err := tml.Execute(buf, req.TemplateArgs); err != nil {
		return err
	}

	msg := model.Message{
		To:  []model.Address{},
		Cc:  []model.Address{},
		Bcc: []model.Address{},
	}
	if err := yaml.Unmarshal(buf.Bytes(), &msg); err != nil {
		return err
	}

	for _, rec := range req.To {
		msg.To = append(msg.To, rec)
	}

	for _, rec := range req.Cc {
		msg.Cc = append(msg.Cc, rec)
	}

	for _, rec := range req.Bcc {
		msg.Bcc = append(msg.Bcc, rec)
	}

	return a.sender.Send(msg)
}
