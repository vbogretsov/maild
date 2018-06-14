package sendgrid

import (
	"encoding/json"
	"github.com/pkg/errors"
	"net/http"

	api "github.com/sendgrid/sendgrid-go"
	"gopkg.in/go-playground/validator.v9"

	"github.com/vbogretsov/maild/model"
)

const v3URL = "/v3/mail/send"

var valid = validator.New()

type content struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

type personalization struct {
	Subject string          `json:"subject"`
	To      []model.Address `json:"to"`
	Cc      []model.Address `json:"cc,omitempty"`
}

type message struct {
	From            model.Address     `json:"from"`
	Personalization sgPersonalization `json:"personalization"`
	Content         sgContent         `json:"content"`
}

type Sender struct {
	url string
	key string
}

func (self Sender) Send(msg model.Message) error {
	if err := valid.Struct(message); err != nil {
		return err
	}

	data := message{
		From: msg.From,
		Personalization: personalization{
			Subject: msg.Subject,
			To:      msg.To,
			Cc:      msg.Cc,
			Bcc:     msg.Bcc,
		},
		Content: content{
			Type:  msg.BodyType,
			Value: msg.Body,
		},
	}

	data, err := json.Marshal(&data)
	if err != nil {
		return err
	}

	request := api.GetRequest(self.key, v3URL, self.url)
	request.Method = http.MethodPost
	request.Body = data

	resp, err := api.API(request)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		// TODO: log respinse body.
		return errors.Errorf("SendGrid API error %d", resp.StatusCode)
	}

	return nil
}
