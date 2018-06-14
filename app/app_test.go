package app

import (
	"fmt"
	"io"
	"strings"
	"testing"

	"github.com/vbogretsov/gotest"

	"github.com/vbogretsov/maild/model"
)

const example = `
From:
  Email: levelup@levelup.com
  Name: LevelUp
Subject: Subject
BodyType: "text/plain"
Body: |
  Hello {{.Username}}!
  This is test body!
`

type Loader struct {
	templates map[string]string
}

func (self *Loader) Load(lang, name string) (io.Reader, error) {
	id := fmt.Sprintf("%s-%s", lang, name)
	text, ok := self.templates[id]
	if !ok {
		return nil, fmt.Errorf("template %s not found", id)
	}
	return strings.NewReader(text), nil
}

type Sender struct {
	inboxes map[model.Address][]model.Message
}

func (self *Sender) send(msg model.Message, recs []model.Address) {
	for _, rec := range recs {
		_, ok := self.inboxes[rec]
		if !ok {
			self.inboxes[rec] = []model.Message{}
		}
		self.inboxes[rec] = append(self.inboxes[rec], msg)
	}
}

func (self *Sender) Send(msg model.Message) error {
	self.send(msg, msg.To)
	self.send(msg, msg.Cc)
	self.send(msg, msg.Bcc)
	return nil
}

func configure(t *testing.T) *gotest.Session {
	s := gotest.NewSession(t)

	s.Register(gotest.Fixture{
		Provider: func() *Loader {
			return &Loader{
				templates: map[string]string{
					"en-1": example,
				},
			}
		},
	})

	s.Register(gotest.Fixture{
		Provider: func() *Sender {
			return &Sender{
				inboxes: map[model.Address][]model.Message{},
			}
		},
	})

	s.Register(gotest.Fixture{
		Provider: func(loader *Loader, sender *Sender) *App {
			return New(loader, sender)
		},
	})

	return s
}

func TestApp(t *testing.T) {
	s := configure(t)
	defer s.Close()

	s.Run(func(sender *Sender, app *App, assert *gotest.Assert) {
		recipient := "to1@mail.com"
		body := fmt.Sprintf("Hello user@mail.com!\nThis is test body!\n")

		req := model.Request{
			TemplateLang: "en",
			TemplateName: "1",
			TemplateArgs: map[string]interface{}{
				"Username": "user@mail.com",
			},
			To: []model.Address{
				{
					Email: recipient,
				},
			},
		}

		err := app.SendMail(req)

		assert.Nil(err)
		inbox, ok := sender.inboxes[model.Address{Email: recipient}]
		assert.True(ok)
		assert.True(len(inbox) > 0)
		msg := inbox[0]
		assert.Equal(msg.Subject, "Subject")
		assert.Equal(msg.BodyType, "text/plain")
		assert.Equal(msg.From.Email, "levelup@levelup.com")
		assert.Equal(msg.From.Name, "LevelUp")
		assert.Equal(msg.Body, body)

	}, "TestSendMailSuccess")

	s.Run(func(sender *Sender, app *App, assert *gotest.Assert) {
		req := model.Request{
			TemplateLang: "en",
			TemplateName: "1",
			TemplateArgs: map[string]interface{}{
				"Username": "user@mail.com",
			},
		}

		err := app.SendMail(req)

		assert.NotNil(err)

	}, "TestSendMailFailedIfMissingRecipients")

	s.Run(func(sender *Sender, app *App, assert *gotest.Assert) {
		recipient := "to1@mail.com"
		req := model.Request{
			TemplateName: "1",
			TemplateArgs: map[string]interface{}{
				"Username": "user@mail.com",
			},
			To: []model.Address{
				{
					Email: recipient,
				},
			},
		}

		err := app.SendMail(req)

		assert.NotNil(err)

	}, "TestSendMailFailedIfMissingLang")

	s.Run(func(sender *Sender, app *App, assert *gotest.Assert) {
		recipient := "to1@mail.com"
		req := model.Request{
			TemplateLang: "en",
			TemplateArgs: map[string]interface{}{
				"Username": "user@mail.com",
			},
			To: []model.Address{
				{
					Email: recipient,
				},
			},
		}

		err := app.SendMail(req)

		assert.NotNil(err)

	}, "TestSendMailFailedIfMissingName")

	s.Run(func(sender *Sender, app *App, assert *gotest.Assert) {
		recipient := "to1@mail.com"
		req := model.Request{
			TemplateLang: "en",
			TemplateName: "2",
			TemplateArgs: map[string]interface{}{
				"Username": "user@mail.com",
			},
			To: []model.Address{
				{
					Email: recipient,
				},
			},
		}

		err := app.SendMail(req)

		assert.NotNil(err)

	}, "TestSendMailFailedIfTempalteNotFound")
}
