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
)

type testProvider struct {
	Recv model.Message
}

func (tp *testProvider) SendMail(m *model.Message) error {
	tp.Recv = *m
	return nil
}

type senderTest struct {
	Request model.Request
	Message model.Message
	Error   error
}

func testSender(sender *server.Maild, provider *testProvider, st senderTest) func(*testing.T) {
	return func(t *testing.T) {
		err := sender.Send(&st.Request, &struct{}{})
		if err != nil && st.Error == nil {
			t.Errorf("unexpected error %v", err)
		}
		if err == nil && st.Error != nil {
			t.Errorf("expected error got nil")
		}
		if !reflect.DeepEqual(st.Message, provider.Recv) {
			t.Errorf(
				"error: message received do not match expected message: %v",
				pretty.Diff(st.Message, provider.Recv),
			)
		}
	}
}

func TestSender(t *testing.T) {
	templateDB := map[model.TemplateID]string{
		{Lang: "en", Name: "test"}: entest,
	}

	templateLoader := func(key model.TemplateID) (*template.Template, error) {
		if value, ok := templateDB[key]; ok {
			return template.New(key.Lang + key.Name).Parse(value)
		}

		return nil, errors.New("record not found")
	}

	provider := testProvider{}
	sender := server.NewMaild(&provider, templateLoader, 10)

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
	}

	for _, test := range tests {
		t.Run("TestSender", testSender(sender, &provider, test))
	}
}
