package app_test

import (
	"testing"

	"github.com/stretchr/testify/require"

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
			require.Equal(t, fx.result, err)
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
		require.Nil(t, err)
		require.Len(t, sd.inbox, 1)

		exp := model.Message{
			Subject:  "Subject",
			From:     model.Address{Email: "user@mail.com", Name: "Sender"},
			BodyType: "text/plain",
			Body:     expectedBody,
			To:       to,
			Cc:       []model.Address{},
			Bcc:      []model.Address{},
		}
		act := sd.inbox[0]
		require.Equal(t, exp, act)
	})
}
