package sendgrid

import (
	"encoding/json"

	api "github.com/sendgrid/sendgrid-go"
	"gopkg.in/go-playground/validator.v9"

	"github.com/vbogretsov/maild/model"
)

var (
	valid = validator.New()
)

type content struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

type personalizations struct {
	Subject string          `json:"subject"`
	To      []model.Address `json:"to"`
	Cc      []model.Address `json:"cc,omitempty"`
}

type sendGridData struct {
	From             model.Address    `json:"from"`
	Personalizations personalizations `json:"personalization"`
	Content          content          `json:"content"`
}

type Provider struct {
	url string
	key string
}

func (p *Provider) SendMail(message *model.Message) error {
	if err := valid.Struct(message); err != nil {
		return err
	}

	sgdata := sendGridData{
		From: message.From,
		Personalizations: personalizations{
			Subject: message.Subject,
			To:      message.To,
			Cc:      message.Cc,
		},
		Content: content{
			Type:  message.BodyType,
			Value: message.Body,
		},
	}

	data, err := json.Marshal(&sgdata)
	if err != nil {
		return err
	}

	request := api.GetRequest(p.key, "/v3/mail/send", p.url)
	request.Method = "POST"
	request.Body = data

	_, err = api.API(request)
	return err
}

func NewProvider(url, key string) *Provider {
	return &Provider{url: url, key: key}
}
