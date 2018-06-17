package app_test

import (
	"reflect"
	"testing"

	"github.com/kr/pretty"

	"github.com/vbogretsov/maild/app"
	"github.com/vbogretsov/maild/model"
)

func TestSendMail(t *testing.T) {
	lr := NewLoader()
	sd := NewSender()
	ap := app.New(lr, sd)

	for _, fx := range fixtures {
		t.Run(fx.name, func(t *testing.T) {
			err := ap.SendMail(fx.request)
			if !reflect.DeepEqual(fx.result, err) {
				t.Error(pretty.Diff(fx.result, err))
				t.Logf("error: %v", err)
			}
		})
	}

	t.Run("MailSent", func(t *testing.T) {
		to := []model.Address{
			{
				Email: "user@mail.com",
				Name:  "",
			},
		}
		req := model.Request{
			TemplateLang: "en",
			TemplateName: "valid",
			TemplateArgs: map[string]interface{}{
				"Username": "SuperUser",
			},
			To: to,
		}
		err := ap.SendMail(req)
		if err != nil {
			t.Errorf("expected nil, but got error '%v'", err)
			t.FailNow()
		}
		exp := model.Message{
			Subject:  "Subject",
			From:     model.Address{Email: "user@mail.com", Name: "Sender"},
			BodyType: "text/plain",
			Body:     expectedBody,
			To:       to,
			Cc:       []model.Address{},
			Bcc:      []model.Address{},
		}
		if len(sd.inbox) == 0 {
			t.Error("message was not sent")
			t.FailNow()
		}
		if len(sd.inbox) > 1 {
			t.Error("to many messages sent")
			t.FailNow()
		}
		act := sd.inbox[0]
		if !reflect.DeepEqual(exp, act) {
			t.Error(pretty.Diff(exp, act))
		}
	})
}
