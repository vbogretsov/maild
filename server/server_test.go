package server

import (
	"errors"
	"reflect"
	"testing"
	"text/template"

	"github.com/kr/pretty"

	"github.com/vbogretsov/maild/model"
	"github.com/vbogretsov/maild/server"
)

const (
	entest = `
From:
  Email: levelup@levelup.com
  Name: LevelUP
Subject: Subject
BodyType: "text/plain"
Body: |
  Hello {{.Username}}!
  This is test body!
`
	eninval = `
From:
x
Subject: Subject
BodyType: "text/plain"
Body: |
  Hello {{.Username}}!
  This is test body!
`
)

type senderTest struct {
	Request model.Request
	Message model.Message
	Error   error
}

func testSender(sender *server.Maild, m *model.Message, st senderTest) func(*testing.T) {
	return func(t *testing.T) {
		err := sender.Send(st.Request, &struct{}{})
		if err != nil && st.Error == nil {
			t.Errorf("unexpected error %v", err)
		}
		if err == nil && st.Error != nil {
			t.Errorf("expected error got nil")
		}
		if err != nil && st.Error != nil {
			return
		}
		if !reflect.DeepEqual(st.Message, *m) {
			t.Errorf(
				"error: message received do not match expected message: %v",
				pretty.Diff(st.Message, *m),
			)
		}
	}
}

func TestSender(t *testing.T) {
	var expectedMsg model.Message

	templates := template.New("maild")

	templates.AddParseTree(
		"en-test.msg",
		template.Must(template.New("en-test.msg").Parse(entest)).Copy())

	templates.AddParseTree(
		"en-inval.msg",
		template.Must(template.New("en-inval.msg").Parse(eninval)).Copy())

	sendmail := func(m model.Message) error {
		expectedMsg = m
		return nil
	}

	sender := server.NewMaild(sendmail, templates)

	tests := []senderTest{
		{
			Request: model.Request{
				TemplateLang: "en",
				TemplateName: "test",
				TemplateArgs: map[string]interface{}{
					"Username": "Donald",
				},
				To: []model.Address{
					{Email: "to1@mail.com"},
					{Email: "to2@mail.com"},
				},
				Cc: []model.Address{
					{Email: "cc1@mail.com"},
					{Email: "cc2@mail.com"},
				},
			},
			Message: model.Message{
				From: model.Address{
					Name:  "LevelUP",
					Email: "levelup@levelup.com",
				},
				To: []model.Address{
					{Email: "to1@mail.com"},
					{Email: "to2@mail.com"},
				},
				Cc: []model.Address{
					{Email: "cc1@mail.com"},
					{Email: "cc2@mail.com"},
				},
				Subject:  "Subject",
				BodyType: "text/plain",
				Body:     "Hello Donald!\nThis is test body!\n",
			},
			Error: nil,
		},
		{
			Request: model.Request{
				TemplateLang: "en",
				TemplateName: "test",
				TemplateArgs: map[string]interface{}{
					"Username": "Donald",
				},
				To: []model.Address{
					{Email: "to1@mail.com"},
					{Email: "to2@mail.com"},
				},
			},
			Message: model.Message{
				From: model.Address{
					Name:  "LevelUP",
					Email: "levelup@levelup.com",
				},
				To: []model.Address{
					{Email: "to1@mail.com"},
					{Email: "to2@mail.com"},
				},
				Subject:  "Subject",
				BodyType: "text/plain",
				Body:     "Hello Donald!\nThis is test body!\n",
			},
			Error: nil,
		},
		{
			Request: model.Request{
				TemplateLang: "en",
				TemplateName: "test2",
				TemplateArgs: map[string]interface{}{
					"Username": "Donald",
				},
				To: []model.Address{
					{Email: "to1@mail.com"},
					{Email: "to2@mail.com"},
				},
				Cc: []model.Address{
					{Email: "cc1@mail.com"},
					{Email: "cc2@mail.com"},
				},
			},
			Message: model.Message{},
			Error:   errors.New("template not found"),
		},
		{
			Request: model.Request{
				TemplateLang: "en",
				TemplateName: "test",
				TemplateArgs: map[string]interface{}{
					"Username": "Donald",
				},
				Cc: []model.Address{
					{Email: "cc1@mail.com"},
					{Email: "cc2@mail.com"},
				},
			},
			Message: model.Message{},
			Error:   errors.New("invalid request"),
		},
		{
			Request: model.Request{
				TemplateLang: "en",
				TemplateName: "test",
				To: []model.Address{
					{Email: "to1@mail.com"},
					{Email: "to2@mail.com"},
				},
			},
			Message: model.Message{},
			Error:   errors.New("invalid request"),
		},
		{
			Request: model.Request{
				TemplateLang: "en",
				TemplateArgs: map[string]interface{}{
					"Username": "Donald",
				},
				To: []model.Address{
					{Email: "to1@mail.com"},
					{Email: "to2@mail.com"},
				},
			},
			Message: model.Message{},
			Error:   errors.New("template not found"),
		},
		{
			Request: model.Request{
				TemplateName: "test2",
				TemplateArgs: map[string]interface{}{
					"Username": "Donald",
				},
				To: []model.Address{
					{Email: "to1@mail.com"},
					{Email: "to2@mail.com"},
				},
			},
			Message: model.Message{},
			Error:   errors.New("template not found"),
		},
		{
			Request: model.Request{
				TemplateLang: "en",
				TemplateName: "inval",
				TemplateArgs: map[string]interface{}{
					"Username": "Donald",
				},
				To: []model.Address{
					{Email: "to1@mail.com"},
					{Email: "to2@mail.com"},
				},
			},
			Message: model.Message{},
			Error:   errors.New("unable to build message"),
		},
	}

	for _, test := range tests {
		t.Run("TestSender", testSender(sender, &expectedMsg, test))
	}
}
