package server

import (
	"encoding/json"

	"github.com/sendgrid/sendgrid-go"
	"gopkg.in/go-playground/validator.v9"

	"github.com/vbogretsov/maild/model"
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

type sendGridProvider struct {
	url string
	key string
}

func (sp *sendGridProvider) SendMail(message model.Message) error {
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

	request := sendgrid.GetRequest(sp.key, "/v3/mail/send", sp.url)
	request.Method = "POST"
	request.Body = data

	_, err := sendgrid.API(request)
	return err
}

func NewSendGridProvider(url, key string) Provider {
	return &sendgridProvider{
		url: url,
		key: key,
	}
}
